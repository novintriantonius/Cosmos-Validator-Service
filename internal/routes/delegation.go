package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/novintriantonius/cosmos-validator-service/internal/models"
	"github.com/novintriantonius/cosmos-validator-service/internal/store"
)

// DelegationHandler handles delegation-related HTTP requests
type DelegationHandler struct {
	store store.DelegationStore
}

// NewDelegationHandler creates a new delegation handler
func NewDelegationHandler(store store.DelegationStore) *DelegationHandler {
	return &DelegationHandler{store: store}
}

// GetHourlyDelegations handles GET /api/v1/validators/{validator_address}/delegations/hourly
// Returns hourly snapshot of delegations for a validator
func (h *DelegationHandler) GetHourlyDelegations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	validatorAddress := vars["validator_address"]

	// Get all delegations for this validator
	delegations, err := h.store.GetDelegations(validatorAddress)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to retrieve delegations",
			"errors":  []string{err.Error()},
		})
		return
	}

	// Group delegations by hour
	hourlyDelegations := groupDelegationsByHour(delegations)

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"code":    http.StatusOK,
		"message": "Hourly delegations retrieved successfully",
		"data": map[string]interface{}{
			"validator_address": validatorAddress,
			"hourly_delegations": hourlyDelegations,
			"count":              len(hourlyDelegations),
		},
	})
}

// GetDailyDelegations handles GET /api/v1/validators/{validator_address}/delegations/daily
// Returns daily snapshot of delegations for a validator
func (h *DelegationHandler) GetDailyDelegations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	validatorAddress := vars["validator_address"]

	// Get all delegations for this validator
	delegations, err := h.store.GetDelegations(validatorAddress)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to retrieve delegations",
			"errors":  []string{err.Error()},
		})
		return
	}

	// Group delegations by day
	dailyDelegations := groupDelegationsByDay(delegations)

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"code":    http.StatusOK,
		"message": "Daily delegations retrieved successfully",
		"data": map[string]interface{}{
			"validator_address": validatorAddress,
			"daily_delegations": dailyDelegations,
			"count":             len(dailyDelegations),
		},
	})
}

// GetDelegatorHistory handles GET /api/v1/validators/{validator_address}/delegator/{delegator_address}/history
// Returns historical delegation data for a specific delegator
func (h *DelegationHandler) GetDelegatorHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	validatorAddress := vars["validator_address"]
	delegatorAddress := vars["delegator_address"]

	// Get all delegations for this validator
	delegations, err := h.store.GetDelegations(validatorAddress)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to retrieve delegations",
			"errors":  []string{err.Error()},
		})
		return
	}

	// Filter delegations for the specific delegator
	delegatorHistory := filterDelegationsForDelegator(delegations, delegatorAddress)

	// If no delegations found for this delegator
	if len(delegatorHistory) == 0 {
		respondWithJSON(w, http.StatusNotFound, map[string]interface{}{
			"status":  "error",
			"code":    http.StatusNotFound,
			"message": "No delegations found for this delegator",
			"errors":  []string{fmt.Sprintf("No delegation history found for delegator %s with validator %s", delegatorAddress, validatorAddress)},
		})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"code":    http.StatusOK,
		"message": "Delegator history retrieved successfully",
		"data": map[string]interface{}{
			"validator_address": validatorAddress,
			"delegator_address": delegatorAddress,
			"history":           delegatorHistory,
			"count":             len(delegatorHistory),
		},
	})
}

// Helper function to group delegations by hour
func groupDelegationsByHour(delegations []models.Delegation) map[string][]models.Delegation {
	hourlyMap := make(map[string][]models.Delegation)

	for _, delegation := range delegations {
		// Format time to hourly granularity
		hourKey := delegation.CreatedAt.Format("2006-01-02T15:00:00Z")
		hourlyMap[hourKey] = append(hourlyMap[hourKey], delegation)
	}

	return hourlyMap
}

// Helper function to group delegations by day
func groupDelegationsByDay(delegations []models.Delegation) map[string][]models.Delegation {
	dailyMap := make(map[string][]models.Delegation)

	for _, delegation := range delegations {
		// Format time to daily granularity
		dayKey := delegation.CreatedAt.Format("2006-01-02T00:00:00Z")
		dailyMap[dayKey] = append(dailyMap[dayKey], delegation)
	}

	return dailyMap
}

// Helper function to filter delegations for a specific delegator
func filterDelegationsForDelegator(delegations []models.Delegation, delegatorAddress string) []models.Delegation {
	var filteredDelegations []models.Delegation

	for _, delegation := range delegations {
		if delegation.DelegatorAddress == delegatorAddress {
			filteredDelegations = append(filteredDelegations, delegation)
		}
	}

	// Sort by created_at (most recent first)
	// Since delegations are already ordered by created_at DESC in the store, 
	// we don't need to sort them again

	return filteredDelegations
} 