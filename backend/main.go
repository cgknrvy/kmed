package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"kmed/api/cmd"
	"log"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed all:web/dist
var frontend embed.FS

func main() {
	// Strip the "web/dist" prefix so index.html is at "/index.html"
	dist, err := fs.Sub(frontend, "web/dist")
	if err != nil {
		log.Fatal(err)
	}

	handler := http.NewServeMux()

	// API
	api := cmd.NewApi()
	handler.Handle("/v1/", http.StripPrefix("/v1", api.Router))

	// Tanstack Router react frontend
	handler.Handle("/", spaHandler(dist))

	server := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(&server, api)

	fmt.Printf("starting server on %s\n", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s\n", err))
	}

	log.Println("Graceful shutdown complete.")
}

// spaHandler handles requests to the single page react application.
func spaHandler(fsys fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(fsys))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")

		if path != "" {
			// If the requested file exists serve it
			if f, err := fsys.Open(path); err == nil {
				defer f.Close()

				if stat, err := f.Stat(); err == nil && !stat.IsDir() {
					fileServer.ServeHTTP(w, r)
					return
				}
			}
		}

		// Otherwise let the SPA router handle it
		http.ServeFileFS(w, r, fsys, "index.html")
	})
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
