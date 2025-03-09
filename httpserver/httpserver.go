package httpserver

import (
	"context"
	"errors"
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

func setupHandlers(router *mux.Router, schedulerCancel context.CancelFunc) {
	// Basic ping
	helloFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})
	router.HandleFunc("/", helloFunc).Methods("GET")

	// REST api - Validator
	uploadAndValidate := http.HandlerFunc(upload.HandleUploadFile)
	router.HandleFunc("/upload/xsd", uploadAndValidate).Methods("POST")
	router.HandleFunc("/upload/rng", uploadAndValidate).Methods("POST")

	// Secure
	router.HandleFunc("/login", secure.Login).Methods("POST")
	router.Handle("/secure", secure.Protect(helloFunc)).Methods("GET")

	// SOAP

	// XML-RPC
	xmlrpc.RegisterEndpoint(router)
}
