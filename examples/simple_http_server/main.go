package main

import (
	"context"
	"fmt"
	"github.com/NicolasBissig/gotel"
	"github.com/NicolasBissig/gotel/gotelhttp"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() (err error) {
	// Set up OpenTelemetry.
	sdk, err := gotel.Setup()
	if err != nil {
		return fmt.Errorf("failed to set up OpenTelemetry: %w", err)
	}

	// Start HTTP server.
	srv := &http.Server{
		Addr:         ":8080",
		BaseContext:  func(_ net.Listener) context.Context { return sdk.RootContext },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      newHTTPHandler(),
	}
	srvErr := make(chan error, 1)

	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	return <-srvErr
}

func newHTTPHandler() http.Handler {
	mux := gotelhttp.NewServeMux()

	mux.HandleFunc("GET /rolldice", rolldice)

	return mux
}
