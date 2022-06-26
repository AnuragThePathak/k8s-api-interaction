package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

type ServerConfig struct {
	Port int
}

type Server interface {
	ListenAndServe()
}

func NewServer(
	endpoints []Endpoints,
	config ServerConfig,
	logger *zap.Logger,
) Server {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	for _, endpoint := range endpoints {
		endpoint.Register(router)
	}

	return &server{
		config: config,
		logger: logger,
		handler: cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST"},
			AllowedHeaders: []string{"*"},
		}).Handler(router),
	}
}

type server struct {
	handler http.Handler
	logger  *zap.Logger
	config  ServerConfig
}

// serve the http requests
func (s *server) ListenAndServe() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: s.handler,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sig
		s.logger.Info("Shutting down server...")

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				s.logger.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			s.logger.Fatal(err.Error())
		}
		serverStopCtx()
	}()

	// Run the server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Fatal(err.Error())
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
