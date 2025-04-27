package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/novintriantonius/cosmos-validator-service/internal/services"
	"github.com/novintriantonius/cosmos-validator-service/internal/store"
)

// SetupRouter configures all the routes for the application
func SetupRouter(validatorStore store.ValidatorStore, cosmosService *services.CosmosService) *mux.Router {
	router := mux.NewRouter()
	
	// Create handler instances
	validatorHandler := NewValidatorHandler(validatorStore)
	
	// Validator routes
	router.HandleFunc("/validators", validatorHandler.GetAll).Methods("GET")
	router.HandleFunc("/validators/{address}", validatorHandler.GetByAddress).Methods("GET")
	router.HandleFunc("/validators", validatorHandler.Create).Methods("POST")
	router.HandleFunc("/validators/{address}", validatorHandler.Update).Methods("PUT")
	router.HandleFunc("/validators/{address}", validatorHandler.Delete).Methods("DELETE")
	
	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Service is healthy"))
	}).Methods("GET")
	
	return router
} 