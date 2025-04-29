package tasks

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/novintriantonius/cosmos-validator-service/internal/services"
	"github.com/novintriantonius/cosmos-validator-service/internal/store"
)

// DelegationSyncTask handles the periodic syncing of delegations for validators
type DelegationSyncTask struct {
	validatorStore      store.ValidatorStore
	delegationStore     store.DelegationStore
	cosmosService       *services.CosmosService
	lastRunStats        SyncStats
	totalDelegationsSynced int
}

// SyncStats contains statistics about the delegation sync
type SyncStats struct {
	TotalRuns           int       // Total number of times the sync has run
	SuccessCount        int       // Total successful validators processed
	ErrorCount          int       // Total validators that had errors
	SkippedCount        int       // Total validators skipped (no changes)
	LastRunTime         time.Time // When the sync last ran
	LastRunDuration     time.Duration
	TotalDelegationsProcessed int
}

// NewDelegationSyncTask creates a new delegation sync task
func NewDelegationSyncTask(
	validatorStore store.ValidatorStore,
	delegationStore store.DelegationStore,
	cosmosService *services.CosmosService,
) *DelegationSyncTask {
	return &DelegationSyncTask{
		validatorStore:  validatorStore,
		delegationStore: delegationStore,
		cosmosService:   cosmosService,
	}
}

// SyncEnabledValidatorDelegations syncs delegations for all enabled validators
func (t *DelegationSyncTask) SyncEnabledValidatorDelegations(ctx context.Context) error {
	log.Printf("[DEBUG] Starting SyncEnabledValidatorDelegations")
	
	// Get all enabled validators
	validators, err := t.validatorStore.GetEnabledValidators()
	if err != nil {
		log.Printf("[ERROR] Failed to get enabled validators: %v", err)
		return fmt.Errorf("error getting enabled validators: %v", err)
	}
	log.Printf("[DEBUG] Found %d enabled validators", len(validators))

	for _, validatorAddress := range validators {
		log.Printf("[DEBUG] Processing validator: %s", validatorAddress)
		
		// Get delegations from API
		delegations, err := t.cosmosService.RetrieveDelegations(ctx, validatorAddress)
		if err != nil {
			log.Printf("[ERROR] Failed to get delegations for validator %s: %v", validatorAddress, err)
			continue
		}
		log.Printf("[DEBUG] Retrieved %d delegations from API for validator %s", 
			len(delegations.DelegationResponses), validatorAddress)
		
		// Log delegation details for debugging
		for i, resp := range delegations.DelegationResponses {
			log.Printf("[DEBUG] Delegation %d: delegator=%s, shares=%s", 
				i+1, resp.Delegation.DelegatorAddress, resp.Delegation.Shares)
		}

		// Save delegations to store
		if err := t.delegationStore.SaveDelegations(validatorAddress, *delegations); err != nil {
			log.Printf("[ERROR] Failed to save delegations for validator %s: %v", validatorAddress, err)
			continue
		}
		log.Printf("[INFO] Successfully processed delegations for validator %s", validatorAddress)
	}

	return nil
}

// GetSyncStats returns the statistics about delegation syncing
func (t *DelegationSyncTask) GetSyncStats() SyncStats {
	return t.lastRunStats
}

// GetTotalDelegationsSynced returns the total number of delegations synced
func (t *DelegationSyncTask) GetTotalDelegationsSynced() int {
	return t.totalDelegationsSynced
} 