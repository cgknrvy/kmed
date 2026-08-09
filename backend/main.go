package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"kmed/api/cmd"
	"kmed/api/internal/httpx"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	api := cmd.NewApi()

	root := http.NewServeMux()
	root.Handle("/v1/", http.StripPrefix("/v1", api.Router))

	server := http.Server{
		Addr:    ":8080",
		Handler: httpx.CORS(root),
	}

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(&server, api)

	fmt.Printf("starting server on %s\n", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s\n", err))
	}

	log.Println("Graceful shutdown complete.")
}

func gracefulShutdown(server *http.Server, api *cmd.Api) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// Close the Api
	if err := api.Close(); err != nil {
		log.Printf("api failed to close successfully: %v\n", err)
	}

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v\n", err)
	}

	log.Println("Server exiting")
}
