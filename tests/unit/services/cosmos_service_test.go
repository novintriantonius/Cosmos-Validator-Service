package services_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/novintriantonius/cosmos-validator-service/internal/services"
)

func TestNewCosmosService(t *testing.T) {
	service := services.NewCosmosService()
	
	if service.GetConfig().BaseURL != services.DefaultBaseURL {
		t.Errorf("Expected BaseURL to be %s, got %s", services.DefaultBaseURL, service.GetConfig().BaseURL)
	}
	
	if service.GetConfig().MaxRetries != services.DefaultMaxRetries {
		t.Errorf("Expected MaxRetries to be %d, got %d", services.DefaultMaxRetries, service.GetConfig().MaxRetries)
	}
	
	if service.GetConfig().RetryDelay != services.DefaultRetryDelay*time.Millisecond {
		t.Errorf("Expected RetryDelay to be %d, got %d", services.DefaultRetryDelay*time.Millisecond, service.GetConfig().RetryDelay)
	}
	
	if service.GetConfig().Timeout != services.DefaultTimeout*time.Second {
		t.Errorf("Expected Timeout to be %d, got %d", services.DefaultTimeout*time.Second, service.GetConfig().Timeout)
	}
}

func TestNewCosmosServiceWithConfig(t *testing.T) {
	customClient := &http.Client{Timeout: 20 * time.Second}
	config := services.CosmosServiceConfig{
		BaseURL:    "https://custom-api.example.com",
		MaxRetries: 5,
		RetryDelay: 1000 * time.Millisecond,
		Timeout:    15 * time.Second,
		HTTPClient: customClient,
	}
	
	service := services.NewCosmosServiceWithConfig(config)
	
	if service.GetConfig().BaseURL != config.BaseURL {
		t.Errorf("Expected BaseURL to be %s, got %s", config.BaseURL, service.GetConfig().BaseURL)
	}
	
	if service.GetConfig().MaxRetries != config.MaxRetries {
		t.Errorf("Expected MaxRetries to be %d, got %d", config.MaxRetries, service.GetConfig().MaxRetries)
	}
	
	if service.GetConfig().RetryDelay != config.RetryDelay {
		t.Errorf("Expected RetryDelay to be %d, got %d", config.RetryDelay, service.GetConfig().RetryDelay)
	}
	
	if service.GetConfig().Timeout != config.Timeout {
		t.Errorf("Expected Timeout to be %d, got %d", config.Timeout, service.GetConfig().Timeout)
	}
}

func TestRetrieveDelegations_EmptyAddress(t *testing.T) {
	service := services.NewCosmosService()
	_, err := service.RetrieveDelegations(context.Background(), "")
	
	if err == nil {
		t.Error("Expected error for empty validator address, got nil")
	}
}

func TestRetrieveDelegations_SuccessfulRequest(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the path is correct
		expectedPath := "/cosmos/staking/v1beta1/validators/cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s/delegations"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		
		// Send a sample response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"delegation_responses": [
				{
					"delegation": {
						"delegator_address": "cosmos1qqqc64yy3qkvxrmp5mlpattr6wpnxvt4qrdtlg",
						"validator_address": "cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s",
						"shares": "14001399.971292884114943026"
					},
					"balance": {
						"denom": "uatom",
						"amount": "14000000"
					}
				}
			],
			"pagination": {
				"next_key": "AokxoW+kv3CwnEI4DGW35C9REtY=",
				"total": "10546"
			}
		}`))
	}))
	defer server.Close()
	
	// Create a service with the test server URL
	config := services.CosmosServiceConfig{
		BaseURL:    server.URL,
		MaxRetries: 1,
		RetryDelay: 100 * time.Millisecond,
		Timeout:    5 * time.Second,
	}
	service := services.NewCosmosServiceWithConfig(config)
	
	// Test the RetrieveDelegations method
	resp, err := service.RetrieveDelegations(context.Background(), "cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s")
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if resp == nil {
		t.Fatal("Expected response, got nil")
	}
	
	if len(resp.DelegationResponses) != 1 {
		t.Errorf("Expected 1 delegation response, got %d", len(resp.DelegationResponses))
	}
	
	if resp.DelegationResponses[0].Delegation.ValidatorAddress != "cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s" {
		t.Errorf("Expected validator address cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s, got %s", 
			resp.DelegationResponses[0].Delegation.ValidatorAddress)
	}
	
	if resp.Pagination.Total != "10546" {
		t.Errorf("Expected pagination total 10546, got %s", resp.Pagination.Total)
	}
} 