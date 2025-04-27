package delegations_test

import (
	"context"
	"testing"
	"time"

	"github.com/novintriantonius/cosmos-validator-service/internal/services"
)

// TestRetrieveDelegationsE2E is an e2e test that retrieves delegations from the actual Cosmos API
// Skip this test by default as it requires an internet connection and makes actual API calls
func TestRetrieveDelegationsE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	// Create a service with slightly longer timeouts for e2e testing
	config := services.CosmosServiceConfig{
		MaxRetries: 3,
		RetryDelay: 1000 * time.Millisecond,
		Timeout:    15 * time.Second,
	}
	service := services.NewCosmosServiceWithConfig(config)
	
	// Set a timeout for the entire test
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Use the Binance validator address as an example
	validatorAddress := "cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s"
	
	// Fetch delegations
	resp, err := service.RetrieveDelegations(ctx, validatorAddress)
	
	// Check for errors
	if err != nil {
		t.Fatalf("Failed to retrieve delegations: %v", err)
	}
	
	// Verify that we got some responses
	if resp == nil {
		t.Fatal("Expected non-nil response")
	}
	
	if len(resp.DelegationResponses) == 0 {
		t.Error("Expected at least one delegation response")
	}
	
	// Verify that the response is properly structured
	for i, delegation := range resp.DelegationResponses {
		// Check validator address
		if delegation.Delegation.ValidatorAddress != validatorAddress {
			t.Errorf("Delegation %d has incorrect validator address: expected %s, got %s", 
				i, validatorAddress, delegation.Delegation.ValidatorAddress)
		}
		
		// Check that delegator address is present
		if delegation.Delegation.DelegatorAddress == "" {
			t.Errorf("Delegation %d has empty delegator address", i)
		}
		
		// Check that shares is present
		if delegation.Delegation.Shares == "" {
			t.Errorf("Delegation %d has empty shares", i)
		}
		
		// Check that balance has proper denomination
		if delegation.Balance.Denom != "uatom" {
			t.Errorf("Delegation %d has unexpected denom: expected uatom, got %s", 
				i, delegation.Balance.Denom)
		}
	}
	
	// Check pagination info
	if resp.Pagination.Total == "" {
		t.Error("Expected pagination total to be non-empty")
	}
} 