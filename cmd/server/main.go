package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/novintriantonius/cosmos-validator-service/internal/handlers"
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
	// Initialize configuration
	cfg := NewConfig()
	
	// Initialize validator store
	validatorStore := store.NewInMemoryValidatorStore()
	
	// Setup router
	router := handlers.SetupRouter(validatorStore)
	
	// Start server
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
} 