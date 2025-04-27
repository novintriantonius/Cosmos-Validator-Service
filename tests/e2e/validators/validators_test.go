package validators_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/novintriantonius/cosmos-validator-service/internal/handlers"
	"github.com/novintriantonius/cosmos-validator-service/internal/models"
	"github.com/novintriantonius/cosmos-validator-service/internal/store"
)

// TestValidatorCRUDE2E tests the CRUD operations for validators in an e2e fashion
func TestValidatorCRUDE2E(t *testing.T) {
	// Setup dependencies
	validatorStore := store.NewInMemoryValidatorStore()
	router := handlers.SetupRouter(validatorStore)
	
	// Create an HTTP test server
	server := httptest.NewServer(router)
	defer server.Close()
	
	baseURL := server.URL
	
	// 1. Create a validator
	validator := models.Validator{
		Name:            "Binance Node",
		Address:         "cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s",
		EnabledTracking: true,
	}
	
	// Convert to JSON
	validatorJSON, err := json.Marshal(validator)
	if err != nil {
		t.Fatalf("Failed to marshal validator: %v", err)
	}
	
	// Make the create request
	createURL := fmt.Sprintf("%s/validators", baseURL)
	createResp, err := http.Post(createURL, "application/json", bytes.NewBuffer(validatorJSON))
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}
	defer createResp.Body.Close()
	
	// Check status code
	if createResp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, createResp.StatusCode)
	}
	
	// Parse response
	var createdValidator models.Validator
	if err := json.NewDecoder(createResp.Body).Decode(&createdValidator); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	// Check response data
	if createdValidator.Name != validator.Name || 
		createdValidator.Address != validator.Address || 
		createdValidator.EnabledTracking != validator.EnabledTracking {
		t.Errorf("Created validator doesn't match: expected %+v, got %+v", validator, createdValidator)
	}
	
	// 2. Get all validators
	getAllURL := fmt.Sprintf("%s/validators", baseURL)
	getAllResp, err := http.Get(getAllURL)
	if err != nil {
		t.Fatalf("Failed to get all validators: %v", err)
	}
	defer getAllResp.Body.Close()
	
	// Check status code
	if getAllResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, getAllResp.StatusCode)
	}
	
	// Parse response
	var getAllResponse map[string]interface{}
	if err := json.NewDecoder(getAllResp.Body).Decode(&getAllResponse); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	// Check data
	if count, ok := getAllResponse["count"].(float64); !ok || count != 1 {
		t.Errorf("Expected count 1, got %v", getAllResponse["count"])
	}
	
	// 3. Get validator by address
	getByAddressURL := fmt.Sprintf("%s/validators/%s", baseURL, validator.Address)
	getByAddressResp, err := http.Get(getByAddressURL)
	if err != nil {
		t.Fatalf("Failed to get validator by address: %v", err)
	}
	defer getByAddressResp.Body.Close()
	
	// Check status code
	if getByAddressResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, getByAddressResp.StatusCode)
	}
	
	// Parse response
	var retrievedValidator models.Validator
	if err := json.NewDecoder(getByAddressResp.Body).Decode(&retrievedValidator); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	// Check retrieved data
	if retrievedValidator.Name != validator.Name || 
		retrievedValidator.Address != validator.Address || 
		retrievedValidator.EnabledTracking != validator.EnabledTracking {
		t.Errorf("Retrieved validator doesn't match: expected %+v, got %+v", validator, retrievedValidator)
	}
	
	// 4. Update validator
	updatedValidator := models.Validator{
		Name:            "Updated Binance Node",
		EnabledTracking: false,
	}
	
	// Convert to JSON
	updatedValidatorJSON, err := json.Marshal(updatedValidator)
	if err != nil {
		t.Fatalf("Failed to marshal updated validator: %v", err)
	}
	
	// Create PUT request
	updateURL := fmt.Sprintf("%s/validators/%s", baseURL, validator.Address)
	updateReq, err := http.NewRequest(http.MethodPut, updateURL, bytes.NewBuffer(updatedValidatorJSON))
	if err != nil {
		t.Fatalf("Failed to create update request: %v", err)
	}
	updateReq.Header.Set("Content-Type", "application/json")
	
	// Send update request
	client := &http.Client{}
	updateResp, err := client.Do(updateReq)
	if err != nil {
		t.Fatalf("Failed to update validator: %v", err)
	}
	defer updateResp.Body.Close()
	
	// Check status code
	if updateResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, updateResp.StatusCode)
	}
	
	// Parse response
	var finalValidator models.Validator
	if err := json.NewDecoder(updateResp.Body).Decode(&finalValidator); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	// Check updated data
	if finalValidator.Name != updatedValidator.Name || 
		finalValidator.EnabledTracking != updatedValidator.EnabledTracking {
		t.Errorf("Updated validator doesn't match: expected name=%s, enabledTracking=%v; got name=%s, enabledTracking=%v", 
			updatedValidator.Name, updatedValidator.EnabledTracking, 
			finalValidator.Name, finalValidator.EnabledTracking)
	}
	
	// Check that address didn't change
	if finalValidator.Address != validator.Address {
		t.Errorf("Address changed: expected %s, got %s", validator.Address, finalValidator.Address)
	}
	
	// 5. Delete validator
	deleteURL := fmt.Sprintf("%s/validators/%s", baseURL, validator.Address)
	deleteReq, err := http.NewRequest(http.MethodDelete, deleteURL, nil)
	if err != nil {
		t.Fatalf("Failed to create delete request: %v", err)
	}
	
	// Send delete request
	deleteResp, err := client.Do(deleteReq)
	if err != nil {
		t.Fatalf("Failed to delete validator: %v", err)
	}
	defer deleteResp.Body.Close()
	
	// Check status code
	if deleteResp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, deleteResp.StatusCode)
	}
	
	// 6. Verify validator is deleted
	checkDeletedURL := fmt.Sprintf("%s/validators/%s", baseURL, validator.Address)
	checkDeletedResp, err := http.Get(checkDeletedURL)
	if err != nil {
		t.Fatalf("Failed to check deleted validator: %v", err)
	}
	defer checkDeletedResp.Body.Close()
	
	// Check status code - should be not found
	if checkDeletedResp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, checkDeletedResp.StatusCode)
	}
} 