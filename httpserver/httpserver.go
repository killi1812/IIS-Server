package httpserver

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

const port = ":5555"

func Start(ctx context.Context, wg *sync.WaitGroup, schedulerCancel context.CancelFunc) {
	defer wg.Done()

	mux := http.NewServeMux()
	setupHandlers(mux, schedulerCancel)

	// setup server
	srv := &http.Server{
		Addr:         port,
		Handler:      mux,
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
	zap.S().Debugf("Started HTTP listen, address = http://localhost%v", srv.Addr)

	// wait for context cancellation
	<-ctx.Done()

	// shutdown server
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer timeoutCancel()
	err := srv.Shutdown(timeoutCtx)
	if err != nil {
		zap.S().Errorf("Cannot shut down HTTP server, err = %v", err)
	}
	zap.S().Debugf("HTTP server was shut down")
}

func setupHandlers(mux *http.ServeMux, schedulerCancel context.CancelFunc) {
	/*
		// server management
		getServerStatus := http.HandlerFunc(mgmt.GetServerStatus)
		getServerVersion := http.HandlerFunc(mgmt.GetServerVersion)
		getServerStatistics := http.HandlerFunc(mgmt.GetServerStatistics)
		// for shutdown, get handler function
		shutdownServer := mgmt.ShutdownServer(schedulerCancel)

		mux.HandleFunc("/getServerStatus", getServerStatus)
		mux.HandleFunc("/getServerVersion", getServerVersion)
		mux.HandleFunc("/getServerStatistics", getServerStatistics)
		mux.HandleFunc("/shutdownServer", shutdownServer)

		// upload
		uploadDocument := http.HandlerFunc(upload.UploadDocument)
		uploadFile := http.HandlerFunc(upload.UploadFile)

		mux.HandleFunc("/uploadDocument", uploadDocument)
		mux.HandleFunc("/uploadFile", uploadFile)

		// data
		getData := http.HandlerFunc(data.GetData)
		mux.HandleFunc("/getData", getData)
	*/
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})
}
