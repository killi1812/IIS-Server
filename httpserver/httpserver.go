package httpserver

import (
	"context"
	"errors"
	"iis_server/apidata"
	"iis_server/httpserver/httpio"
	"iis_server/httpserver/restapi/secure"
	"iis_server/httpserver/restapi/upload"
	"iis_server/httpserver/xmlrpc"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const port = ":5555"

func Start(ctx context.Context, wg *sync.WaitGroup, schedulerCancel context.CancelFunc) {
	defer wg.Done()

	router := mux.NewRouter()
	setupHandlers(router, schedulerCancel)

	srv := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	go func() {
		err := srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			zap.S().Debugf("HTTP server is closed")
		} else if err != nil {
			zap.S().Errorf("Cannot start HTTP server, err = %v", err)
		}
	}()
	zap.S().Infof("Started HTTP listen, address = http://localhost%v", srv.Addr)

	// wait for context cancellation
	<-ctx.Done()

	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer timeoutCancel()
	err := srv.Shutdown(timeoutCtx)
	if err != nil {
		zap.S().Errorf("Cannot shut down HTTP server, err = %v", err)
	}
	zap.S().Info("HTTP server was shut down")
}

func headersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		headers.Add("Access-Control-Allow-Origin", "http://localhost:3000")
		headers.Add("Vary", "Origin")
		headers.Add("Vary", "Access-Control-Request-Method")
		headers.Add("Vary", "Access-Control-Request-Headers")
		headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token, Access-Control-Allow-Origin, Authorization")
		headers.Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
		headers.Add("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func setupHandlers(router *mux.Router, schedulerCancel context.CancelFunc) {
	router.Use(headersMiddleware)

	// Basic ping
	helloFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})
	router.HandleFunc("/", helloFunc).Methods("GET", "OPTIONS")
	router.Handle("/secure", secure.Protect(helloFunc)).Methods("GET", "OPTIONS")

	// REST api - Validator
	uploadAndValidate := http.HandlerFunc(upload.HandleUploadFile)
	router.HandleFunc("/upload/xsd", uploadAndValidate).Methods("POST", "OPTIONS")
	router.HandleFunc("/upload/rng", uploadAndValidate).Methods("POST", "OPTIONS")

	router.HandleFunc("/search/{username}", handleSearch).Methods("GET", "OPTIONS")

	// Secure
	secure.RegisterEndpoints(router)

	// XML-RPC
	xmlrpc.RegisterEndpoint(router)
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		zap.S().Error("Username parameter missing!")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	zap.S().Debug("Search for user with username = %s", username)
	rez, err := apidata.Search(username)
	if err != nil {
		zap.S().Errorf("Failed to find data err = %+v", err)
		httpio.WriteStandardHTTPResponse(w, http.StatusInternalServerError, "", err)
		return
	}

	httpio.WriteStandardHTTPResponse(w, http.StatusOK, rez, nil)
}
