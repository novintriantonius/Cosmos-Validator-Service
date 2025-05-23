package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/novintriantonius/cosmos-validator-service/internal/models"
	"github.com/novintriantonius/cosmos-validator-service/internal/store"
)

// ValidatorHandler handles validator-related HTTP requests
type ValidatorHandler struct {
	store store.ValidatorStore
}

// NewValidatorHandler creates a new validator handler
func NewValidatorHandler(store store.ValidatorStore) *ValidatorHandler {
	return &ValidatorHandler{store: store}
}

// GetAll handles GET /validators
func (h *ValidatorHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	validators, err := h.store.GetAll()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status": "error",
			"code": http.StatusInternalServerError,
			"message": "Failed to retrieve validators",
			"errors": []string{err.Error()},
		})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"code": http.StatusOK,
		"message": "Validators retrieved successfully",
		"data": map[string]interface{}{
			"validators": validators,
			"count": len(validators),
		},
	})
}

// GetByAddress handles GET /validators/{address}
func (h *ValidatorHandler) GetByAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	validator, err := h.store.GetByAddress(address)
	if err == store.ErrValidatorNotFound {
		respondWithJSON(w, http.StatusNotFound, map[string]interface{}{
			"status": "error",
			"code": http.StatusNotFound,
			"message": "Validator not found",
			"errors": []string{"No validator found with address: " + address},
		})
		return
	} else if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status": "error",
			"code": http.StatusInternalServerError,
			"message": "Failed to retrieve validator",
			"errors": []string{err.Error()},
		})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"code": http.StatusOK,
		"message": "Validator retrieved successfully",
		"data": validator,
	})
}

// Create handles POST /validators
func (h *ValidatorHandler) Create(w http.ResponseWriter, r *http.Request) {
	var validator models.Validator
	if err := json.NewDecoder(r.Body).Decode(&validator); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status": "error",
			"code": http.StatusBadRequest,
			"message": "Invalid request body",
			"errors": []string{err.Error()},
		})
		return
	}

	// Validate required fields
	if validator.Address == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status": "error",
			"code": http.StatusBadRequest,
			"message": "Validation failed",
			"errors": []string{"Address is required"},
		})
		return
	}

	if validator.Name == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status": "error",
			"code": http.StatusBadRequest,
			"message": "Validation failed",
			"errors": []string{"Name is required"},
		})
		return
	}

	// Check if validator already exists
	existingValidator, err := h.store.GetByAddress(validator.Address)
	if err == nil && existingValidator != nil {
		respondWithJSON(w, http.StatusConflict, map[string]interface{}{
			"status": "error",
			"code": http.StatusConflict,
			"message": "Validator already exists",
			"errors": []string{
				fmt.Sprintf("A validator with address '%s' already exists", validator.Address),
				fmt.Sprintf("Existing validator name: '%s'", existingValidator.Name),
			},
		})
		return
	}

	if err := h.store.Add(validator); err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status": "error",
			"code": http.StatusInternalServerError,
			"message": "Failed to create validator",
			"errors": []string{err.Error()},
		})
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"status": "success",
		"code": http.StatusCreated,
		"message": "Validator created successfully",
		"data": validator,
	})
}

// Update handles PUT /validators/{address}
func (h *ValidatorHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	var validator models.Validator
	if err := json.NewDecoder(r.Body).Decode(&validator); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status": "error",
			"code": http.StatusBadRequest,
			"message": "Invalid request body",
			"errors": []string{err.Error()},
		})
		return
	}

	if err := h.store.Update(address, validator); err == store.ErrValidatorNotFound {
		respondWithJSON(w, http.StatusNotFound, map[string]interface{}{
			"status": "error",
			"code": http.StatusNotFound,
			"message": "Validator not found",
			"errors": []string{"No validator found with address: " + address},
		})
		return
	} else if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status": "error",
			"code": http.StatusInternalServerError,
			"message": "Failed to update validator",
			"errors": []string{err.Error()},
		})
		return
	}

	// Get the updated validator to return in the response
	updatedValidator, _ := h.store.GetByAddress(address)
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"code": http.StatusOK,
		"message": "Validator updated successfully",
		"data": updatedValidator,
	})
}

// Delete handles DELETE /validators/{address}
func (h *ValidatorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	if err := h.store.Delete(address); err == store.ErrValidatorNotFound {
		respondWithJSON(w, http.StatusNotFound, map[string]interface{}{
			"status": "error",
			"code": http.StatusNotFound,
			"message": "Validator not found",
			"errors": []string{"No validator found with address: " + address},
		})
		return
	} else if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status": "error",
			"code": http.StatusInternalServerError,
			"message": "Failed to delete validator",
			"errors": []string{err.Error()},
		})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"code": http.StatusOK,
		"message": "Validator deleted successfully",
		"data": nil,
	})
}

// respondWithJSON is a helper function to write a JSON response
func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
} 