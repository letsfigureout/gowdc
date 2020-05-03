package main

import (
	"context"
	"fmt"
	"gowdc/cmd/gowdc/internal/handlers"
	"gowdc/internal/web"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	if err := run(); err != nil {
		log.Printf("error with main: %v", err)
	}
}

func run() error {

	// Logger
	log := log.New(os.Stdout, "GoWDC: ", log.LstdFlags|log.Lmicroseconds)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Error channel
	apiErrors := make(chan error, 1)

	router := handlers.NewRouter(log)
	srv := web.NewBackendServer(router, "0.0.0.0:8000", 60, 60, 60)

	go func() {
		log.Printf("Starting API on %v", srv.Addr)
		apiErrors <- srv.Server.ListenAndServe()
	}()

	select {
	case err := <-apiErrors:
		return fmt.Errorf("api error: %v", err)

	case r := <-shutdown:
		log.Printf("sigterm detected: %v, attempting graceful shutdown", r)
		ctx, cancel := context.WithTimeout(context.Background(), srv.ShutdownTimeout)
		defer cancel()

		if err := srv.Server.Shutdown(ctx); err != nil {
			log.Printf("graceful shutdown did not complete within the timeout period %v:%v", srv.ShutdownTimeout, err)
		}
	}

	return nil
}
