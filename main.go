package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	exitOK = iota
	exitError
)

func main() {
	os.Exit(run())
}

func run() int {

	// Prepare mux and listener
	mux := http.NewServeMux()
	server := &http.Server{
		Handler: mux,
	}
	mux.Handle("/", testHandler())

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen HTTP port: %v", err)
		return exitError
	}

	// Run server
	eg, _ := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		return server.Serve(listener)
	})

	// Wait for signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	sig := <-sigCh
	log.Printf("received %v signal. this server will shutdown gracefully.", sig)

	// shutdown server gracefully
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctxTimeout); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	if err := eg.Wait(); err != nil {
		log.Printf("goroutine err message: %v", err)
	}

	return exitOK
}

func testHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.WriteHeader(200)
		if _, err := w.Write([]byte("Hello world!\n")); err != nil {
			log.Printf("Write response data error: %v", err)
		}
	})
}
