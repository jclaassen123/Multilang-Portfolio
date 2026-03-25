package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed static/*
var staticFiles embed.FS

// main bootstraps the Gin server, serves the embedded frontend, and handles
// graceful shutdown requests from either the terminal or the web UI.
func main() {
	gin.SetMode(gin.ReleaseMode)

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("failed to load embedded static files: %v", err)
	}

	shutdownRequested := make(chan struct{}, 1)
	router := newRouter(staticFS, shutdownRequested)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go startServer(server)
	waitForShutdown(shutdownRequested)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}

// newRouter configures the HTTP routes for static assets and the exit endpoint.
func newRouter(staticFS fs.FS, shutdownRequested chan struct{}) *gin.Engine {
	router := gin.Default()
	router.StaticFS("/", http.FS(staticFS))
	router.POST("/exit", handleExitRequest(shutdownRequested))
	return router
}

// handleExitRequest returns a Gin handler that acknowledges the request and
// signals the main goroutine to stop the server.
func handleExitRequest(shutdownRequested chan struct{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Shutting down 2048..."})

		// Non-blocking send prevents duplicate clicks from hanging the request.
		select {
		case shutdownRequested <- struct{}{}:
		default:
		}
	}
}

// startServer begins listening for HTTP requests and stops the program if the
// server fails unexpectedly.
func startServer(server *http.Server) {
	log.Println("2048 is running at http://localhost:8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}

// waitForShutdown blocks until either an OS signal arrives or the web UI asks
// the application to exit.
func waitForShutdown(shutdownRequested <-chan struct{}) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigChan:
		log.Println("shutdown signal received")
	case <-shutdownRequested:
		log.Println("exit requested from web interface")
	}
}
