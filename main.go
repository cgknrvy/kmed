package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"kmed/api/cmd"
	"kmed/api/internal/config"
	"kmed/api/internal/logger"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/browser"
)

const APPNAME = "kmed"

// Injected by GoReleaser's -ldflags during linking
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

//go:embed all:frontend/dist
var frontend embed.FS

func main() {
	configStore := config.NewConfigStore(APPNAME)
	cfg, err := configStore.Load()
	if err != nil {
		slog.Error("failed to load config", "err", err)
		panic(err)
	}

	_, close := logger.Init(cfg.LogFilePath, APPNAME)
	defer close() // close rotator

	// Strip the "frontend/dist" prefix so index.html is at "/index.html"
	dist, err := fs.Sub(frontend, "frontend/dist")
	if err != nil {
		slog.Error("failed to get frontend file system", "err", err)
		panic(err)
	}

	const addr = "127.0.0.1:54322"

	// Try to take the port first. If it's already in use, assume
	// another instance is running and just open the browser to it.
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		if isAddrInUse(err) {
			slog.Info("another instance is already running, opening browser...")
			if err := browser.OpenURL("http://" + addr); err != nil {
				slog.Error("failed to open browser for running instance", "error", err)
			}
			return
		}
		slog.Error("failed to listen", "address", addr, "error", err)
		panic(err)
	}

	handler := http.NewServeMux()

	// API
	api := cmd.NewApi(cfg)
	handler.Handle("/api/v1/", http.StripPrefix("/api/v1", api.Router))

	// Handle the quiting the app so that it can shutdown gracefully
	quitCtx, cancel := context.WithCancel(context.Background())
	handler.Handle("/api/quit", quitHandler(cancel))

	// Tanstack Router react frontend
	handler.Handle("/", spaHandler(dist))

	server := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(&server, api, quitCtx)

	slog.Info("starting server", "port", server.Addr)

	// Open the url in the browser
	go func() {
		if err := browser.OpenURL("http://" + addr); err != nil {
			slog.Warn("failed to open browser", "error", err)
		}
	}()

	if err := server.Serve(ln); err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s\n", err))
	}

	slog.Info("Graceful shutdown complete.")
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

func quitHandler(cancel context.CancelFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Quiting from UI")
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

	slog.Info("shutting down gracefully, press Ctrl+C again to force")

	// Close the Api
	if err := api.Close(); err != nil {
		slog.Warn("api failed to close successfully", "error", err)
	}

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Warn("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exiting")
}

// isAddrInUse checks if the given addres is being used by another process
func isAddrInUse(err error) bool {
	// net.OpError wrapping a syscall.Errno for EADDRINUSE
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		var sysErr *os.SyscallError
		if errors.As(opErr.Err, &sysErr) {
			return errors.Is(sysErr.Err, syscall.EADDRINUSE)
		}
	}
	return errors.Is(err, syscall.EADDRINUSE)
}
