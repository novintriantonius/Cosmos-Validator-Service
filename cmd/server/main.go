package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/novintriantonius/cosmos-validator-service/internal/database"
	"github.com/novintriantonius/cosmos-validator-service/internal/routes"
	"github.com/novintriantonius/cosmos-validator-service/internal/scheduler"
	"github.com/novintriantonius/cosmos-validator-service/internal/services"
	"github.com/novintriantonius/cosmos-validator-service/internal/store"
)

// Config holds application configuration
type Config struct {
	ServerPort int
}

// NewConfig creates a new config with values from environment or defaults
func NewConfig() *Config {
	port := 8080
	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}
	return &Config{
		ServerPort: port,
	}
}

func main() {
	// Initialize database connection
	dbConfig := database.NewConfig()
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize stores with PostgreSQL
	validatorStore := store.NewValidatorStore(db)
	delegationStore := store.NewDelegationStore(db)
	
	// Initialize cosmos service
	cosmosService := services.NewCosmosService()
	
	// Set up router with all dependencies
	router := routes.SetupRouter(validatorStore, delegationStore, cosmosService)
	
	// Initialize and setup scheduler with all tasks
	sched := scheduler.SetupScheduler(validatorStore, delegationStore, cosmosService)
	
	// Start the scheduler
	sched.Start()
	defer sched.Stop()
	
	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", NewConfig().ServerPort),
		Handler: router,
	}
	
	// Start the server in a goroutine
	go func() {
		log.Printf("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	
	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down server...")
	
	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Doesn't block if no connections, but will otherwise wait until the timeout
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Server exited gracefully")
} 