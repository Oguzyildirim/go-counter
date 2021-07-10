package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Oguzyildirim/go-counter/internal/rest"
	"github.com/Oguzyildirim/go-counter/internal/service"
	"github.com/Oguzyildirim/go-counter/internal/storage"
)

func main() {
	var env, address string
	flag.StringVar(&env, "env", "db", "path")
	flag.StringVar(&address, "address", ":9234", "HTTP Server Address")
	flag.Parse()

	if _, err := os.Stat(env); os.IsNotExist(err) {
		_, err := os.Create(env)
		if err != nil {
			log.Fatalf("Couldn't run: %s", err)
		}
	}

	errC, err := run(address, env)
	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("Error while running: %s", err)
	}
}

func run(address, db string) (<-chan error, error) {

	errC := make(chan error, 1)

	srv := newServer(address, db)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		log.Printf("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {
			stop()
			cancel()
			close(errC)
		}()

		srv.SetKeepAlivesEnabled(false)

		if err := srv.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}

		log.Printf("Shutdown completed")
	}()

	go func() {
		log.Printf("Listening and serving %s", address)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errC <- err
		}
	}()

	return errC, nil

}

func newServer(address string, db string) *http.Server {
	r := http.NewServeMux()

	repo := storage.NewCounter(db)
	svc := service.NewCounter(repo)

	rest.NewCounterHandler(svc).Register(r)

	return &http.Server{
		Handler:           r,
		Addr:              address,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}
}
