package validators_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/novintriantonius/cosmos-validator-service/internal/database"
	"github.com/novintriantonius/cosmos-validator-service/internal/models"
	"github.com/novintriantonius/cosmos-validator-service/internal/routes"
	"github.com/novintriantonius/cosmos-validator-service/internal/services"
	"github.com/novintriantonius/cosmos-validator-service/internal/store"
)

func setupTestDB() (*sql.DB, error) {
	// Use an in-memory SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	// Run migrations
	err = database.RunMigrations(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// TestValidatorCRUDE2E tests the CRUD operations for validators in an e2e fashion
func TestValidatorCRUDE2E(t *testing.T) {
	// Skip this test for now until we can set up proper test infrastructure
	t.Skip("Skipping E2E test until proper test database infrastructure is set up")

	// Setup test database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Close()

	// Setup dependencies
	validatorStore := store.NewValidatorStore(db)
	delegationStore := store.NewDelegationStore(db)
	cosmosService := services.NewCosmosService()
	router := routes.SetupRouter(validatorStore, delegationStore, cosmosService)
	
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
	createURL := fmt.Sprintf("%s/api/v1/validators", baseURL)
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
	var createResponse map[string]interface{}
	if err := json.NewDecoder(createResp.Body).Decode(&createResponse); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	// Check response status
	if status, ok := createResponse["status"].(string); !ok || status != "success" {
		t.Errorf("Expected status 'success', got %v", createResponse["status"])
	}
	
	// 2. Get all validators
	getAllURL := fmt.Sprintf("%s/api/v1/validators", baseURL)
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
	
	// Check response status
	if status, ok := getAllResponse["status"].(string); !ok || status != "success" {
		t.Errorf("Expected status 'success', got %v", getAllResponse["status"])
	}
	
	// 3. Get validator by address
	getByAddressURL := fmt.Sprintf("%s/api/v1/validators/%s", baseURL, validator.Address)
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
	var getByAddressResponse map[string]interface{}
	if err := json.NewDecoder(getByAddressResp.Body).Decode(&getByAddressResponse); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	// Check response status
	if status, ok := getByAddressResponse["status"].(string); !ok || status != "success" {
		t.Errorf("Expected status 'success', got %v", getByAddressResponse["status"])
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
	updateURL := fmt.Sprintf("%s/api/v1/validators/%s", baseURL, validator.Address)
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
	var updateResponse map[string]interface{}
	if err := json.NewDecoder(updateResp.Body).Decode(&updateResponse); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	// Check response status
	if status, ok := updateResponse["status"].(string); !ok || status != "success" {
		t.Errorf("Expected status 'success', got %v", updateResponse["status"])
	}
	
	// 5. Delete validator
	deleteURL := fmt.Sprintf("%s/api/v1/validators/%s", baseURL, validator.Address)
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
	if deleteResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, deleteResp.StatusCode)
	}
	
	// Parse response
	var deleteResponse map[string]interface{}
	if err := json.NewDecoder(deleteResp.Body).Decode(&deleteResponse); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	// Check response status
	if status, ok := deleteResponse["status"].(string); !ok || status != "success" {
		t.Errorf("Expected status 'success', got %v", deleteResponse["status"])
	}
	
	// 6. Verify validator is deleted
	checkDeletedURL := fmt.Sprintf("%s/api/v1/validators/%s", baseURL, validator.Address)
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