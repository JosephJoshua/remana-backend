package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/JosephJoshua/repair-management-backend/internal/genapi"
	"github.com/JosephJoshua/repair-management-backend/internal/shared"
)

const (
	ReadHeaderTimeout = 5 * time.Second
	IdleTimeout       = 60 * time.Second
	ReadTimeout       = 30 * time.Second
	WriteTimeout      = 30 * time.Second
	ShutdownTimeout   = 10 * time.Second
)

func run(ctx context.Context, stdout, stderr io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	oasServer, err := genapi.NewServer(shared.Server{}, nil)
	if err != nil {
		fmt.Fprintf(stderr, "run() > error creating server: %s\n", err)
	}

	httpServer := &http.Server{
		Addr:              net.JoinHostPort("localhost", "8080"),
		Handler:           oasServer,
		ReadHeaderTimeout: ReadHeaderTimeout,
		IdleTimeout:       IdleTimeout,
		ReadTimeout:       ReadTimeout,
		WriteTimeout:      WriteTimeout,
	}

	go func() {
		log.Printf("run() > listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(stderr, "run() > error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()

		log.Println("run() > shutting down server")

		shutdownCtx := context.Background()
		shutdownCtx, cancel = context.WithTimeout(shutdownCtx, ShutdownTimeout)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(stderr, "run() > error shutting down http server: %s\n", err)
		}
	}()

	wg.Wait()
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
