package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"context"

	"github.com/streamersonglist/dragonfly-flex/internal/api"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() (err error) {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start HTTP server.
	srv := &http.Server{
		Addr:         ":5500",
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      api.NewHandler(),
	}

	srvErr := make(chan error, 1)

	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
		return
	case <-ctx.Done():
		stop()
	}

	err = srv.Shutdown(context.Background())
	return err
}
