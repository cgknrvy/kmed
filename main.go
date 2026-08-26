package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"kmed/api/cmd"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed all:frontend/dist
var frontend embed.FS

func main() {
	// Strip the "frontend/dist" prefix so index.html is at "/index.html"
	dist, err := fs.Sub(frontend, "frontend/dist")
	if err != nil {
		log.Fatal(err)
	}

	handler := http.NewServeMux()

	// API
	api := cmd.NewApi()
	handler.Handle("/api/v1/", http.StripPrefix("/api/v1", api.Router))

	// Handle the quiting the app so that it can shutdown gracefully
	quitCtx, cancel := context.WithCancel(context.Background())
	handler.Handle("/api/quit", quitHandler(cancel))

	// Tanstack Router react frontend
	handler.Handle("/", spaHandler(dist))

	server := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(&server, api, quitCtx)

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

func quitHandler(cancel context.CancelFunc ) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		log.Println("Quiting from UI")
		// Cancel the quitCtx so that it's Done and the gracefulShutdown can run
		cancel()
	})
}

func gracefulShutdown(server *http.Server, api *cmd.Api, ctx context.Context) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
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
