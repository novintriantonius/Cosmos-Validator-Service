package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/novintriantonius/cosmos-validator-service/internal/models"
	"github.com/novintriantonius/cosmos-validator-service/internal/store"
)

// ValidatorHandler handles HTTP requests for validator endpoints
type ValidatorHandler struct {
	store store.ValidatorStore
}

// NewValidatorHandler creates a new instance of ValidatorHandler
func NewValidatorHandler(store store.ValidatorStore) *ValidatorHandler {
	return &ValidatorHandler{
		store: store,
	}
}

// GetAll handles GET /validators
func (h *ValidatorHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	validators, err := h.store.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data":  validators,
		"count": len(validators),
	})
}

// GetByAddress handles GET /validators/{address}
func (h *ValidatorHandler) GetByAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	validator, err := h.store.GetByAddress(address)
	if err == store.ErrValidatorNotFound {
		http.Error(w, "Validator not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, validator)
}

// Create handles POST /validators
func (h *ValidatorHandler) Create(w http.ResponseWriter, r *http.Request) {
	var validator models.Validator
	if err := json.NewDecoder(r.Body).Decode(&validator); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if validator.Address == "" {
		http.Error(w, "Address is required", http.StatusBadRequest)
		return
	}

	if validator.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if err := h.store.Add(validator); err == store.ErrValidatorAlreadyExists {
		http.Error(w, "Validator with this address already exists", http.StatusConflict)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, validator)
}

// Update handles PUT /validators/{address}
func (h *ValidatorHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	var validator models.Validator
	if err := json.NewDecoder(r.Body).Decode(&validator); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.store.Update(address, validator); err == store.ErrValidatorNotFound {
		http.Error(w, "Validator not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the updated validator to return in the response
	updatedValidator, _ := h.store.GetByAddress(address)
	respondWithJSON(w, http.StatusOK, updatedValidator)
}

// Delete handles DELETE /validators/{address}
func (h *ValidatorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	if err := h.store.Delete(address); err == store.ErrValidatorNotFound {
		http.Error(w, "Validator not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// respondWithJSON is a helper function to write a JSON response
func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
} 