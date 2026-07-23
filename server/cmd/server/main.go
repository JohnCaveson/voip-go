package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/voip-app/internal/signaling"
)

func main() {
	port := flag.Int("port", 9321, "signaling server port")
	flag.Parse()

	hub := signaling.NewHub()

	mux := http.NewServeMux()
	mux.Handle("/signaling", hub)

	addr := fmt.Sprintf(":%d", *port)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		log.Printf("VoIP signaling server listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}
}
