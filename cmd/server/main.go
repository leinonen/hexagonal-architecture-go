package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apiClient "github.com/leinonen/hexagonal-architecture-go/internal/adapters/api"
	httpHandler "github.com/leinonen/hexagonal-architecture-go/internal/adapters/http"
	"github.com/leinonen/hexagonal-architecture-go/internal/adapters/repository/memory"
	"github.com/leinonen/hexagonal-architecture-go/internal/application"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	weatherAPIKey := os.Getenv("WEATHER_API_KEY")
	if weatherAPIKey == "" {
		weatherAPIKey = "demo-key"
		log.Println("Warning: WEATHER_API_KEY not set, using demo key")
	}

	userRepo := memory.NewUserRepository()
	weatherClient := apiClient.NewWeatherClient(weatherAPIKey)

	userService := application.NewUserService(userRepo)
	weatherService := application.NewWeatherService(weatherClient, userRepo)

	handler := httpHandler.NewHandler(userService, weatherService)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			log.Printf("Started %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
			log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
		})
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
