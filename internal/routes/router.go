package routes

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/novintriantonius/cosmos-validator-service/internal/services"
	"github.com/novintriantonius/cosmos-validator-service/internal/store"
)

// SetupRouter configures all the routes for the application
func SetupRouter(validatorStore store.ValidatorStore, delegationStore store.DelegationStore, cosmosService *services.CosmosService) *mux.Router {
	router := mux.NewRouter()
	
	// Create handler instances
	validatorHandler := NewValidatorHandler(validatorStore)
	delegationHandler := NewDelegationHandler(delegationStore)
	
	// API routes
	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	
	// Validator routes
	apiRouter.HandleFunc("/validators", validatorHandler.GetAll).Methods("GET")
	apiRouter.HandleFunc("/validators/{address}", validatorHandler.GetByAddress).Methods("GET")
	apiRouter.HandleFunc("/validators", validatorHandler.Create).Methods("POST")
	apiRouter.HandleFunc("/validators/{address}", validatorHandler.Update).Methods("PUT")
	apiRouter.HandleFunc("/validators/{address}", validatorHandler.Delete).Methods("DELETE")
	
	// Delegation routes
	apiRouter.HandleFunc("/validators/{validator_address}/delegations/hourly", delegationHandler.GetHourlyDelegations).Methods("GET")
	apiRouter.HandleFunc("/validators/{validator_address}/delegations/daily", delegationHandler.GetDailyDelegations).Methods("GET")
	apiRouter.HandleFunc("/validators/{validator_address}/delegator/{delegator_address}/history", delegationHandler.GetDelegatorHistory).Methods("GET")

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Service is healthy"))
	}).Methods("GET")
	
	return router
} 