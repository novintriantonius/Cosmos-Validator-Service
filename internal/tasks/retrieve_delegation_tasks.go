package tasks

import (
	"context"
	"fmt"
	"log"
	"sync"
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

// SyncEnabledValidatorDelegations updates delegations for all enabled validators
func (t *DelegationSyncTask) SyncEnabledValidatorDelegations(ctx context.Context) error {
	// Start time for metrics
	startTime := time.Now()
	
	stats := SyncStats{
		LastRunTime: startTime,
	}
	
	// Get all validators that have tracking enabled
	enabledValidators, err := t.delegationStore.GetEnabledValidators()
	if err != nil {
		return fmt.Errorf("failed to get enabled validators: %w", err)
	}
	
	// If we don't have any enabled validators, check all validators in the store
	if len(enabledValidators) == 0 {
		validators, err := t.validatorStore.GetAll()
		if err != nil {
			return fmt.Errorf("failed to get validators: %w", err)
		}
		
		// Enable tracking for all validators by default
		for _, validator := range validators {
			if err := t.delegationStore.EnableDelegationTracking(validator.Address); err != nil {
				log.Printf("Failed to enable tracking for validator %s: %v", validator.Address, err)
				continue
			}
			enabledValidators = append(enabledValidators, validator.Address)
		}
	}
	
	if len(enabledValidators) == 0 {
		log.Println("No validators to sync delegations for")
		return nil
	}
	
	log.Printf("Syncing delegations for %d enabled validators", len(enabledValidators))
	
	// Use a wait group to process validators concurrently
	var wg sync.WaitGroup
	
	// Create a semaphore to limit concurrent requests
	sem := make(chan struct{}, 5) // Max 5 concurrent requests
	
	// Track previous data state
	type syncResult struct {
		validatorAddress string
		success          bool
		skipped          bool
		error            error
		delegationCount  int
	}
	
	resultCh := make(chan syncResult, len(enabledValidators))
	
	for _, validatorAddress := range enabledValidators {
		wg.Add(1)
		
		// Add to semaphore (blocks if at capacity)
		sem <- struct{}{}
		
		go func(address string) {
			defer wg.Done()
			defer func() { <-sem }() // Release semaphore slot when done
			
			result := syncResult{
				validatorAddress: address,
			}
			
			// Create a timeout context for this specific request
			reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			
			// Check if the validator exists
			_, err := t.validatorStore.GetByAddress(address)
			if err != nil {
				result.error = fmt.Errorf("validator not found: %w", err)
				resultCh <- result
				log.Printf("Validator %s not found in store: %v", address, err)
				return
			}
			
			// Get current stored delegations to check for existing data
			existingData, _ := t.delegationStore.GetDelegations(address)
			
			// Retrieve delegations for this validator
			delegationsResp, err := t.cosmosService.RetrieveDelegations(reqCtx, address)
			if err != nil {
				result.error = fmt.Errorf("failed to retrieve delegations: %w", err)
				resultCh <- result
				log.Printf("Error retrieving delegations for validator %s: %v", address, err)
				return
			}
			
			result.delegationCount = len(delegationsResp.DelegationResponses)
			
			// Save delegations to the store
			if err := t.delegationStore.SaveDelegations(address, *delegationsResp); err != nil {
				result.error = fmt.Errorf("failed to save delegations: %w", err)
				resultCh <- result
				log.Printf("Error saving delegations for validator %s: %v", address, err)
				return
			}
			
			// Get updated data to check if it was actually updated
			updatedData, _ := t.delegationStore.GetDelegations(address)
			
			// Check if the data was actually updated or just the timestamp
			if existingData != nil && updatedData != nil {
				// If timestamps differ by more than just the update time difference, data was actually updated
				dataWasUpdated := !existingData.Data.Timestamp.Equal(updatedData.Data.Timestamp)
				
				if !dataWasUpdated {
					result.skipped = true
					log.Printf("No changes for validator %s - skipped updating %d delegations", 
						address, result.delegationCount)
				} else {
					result.success = true
					log.Printf("Updated %d delegations for validator %s", 
						result.delegationCount, address)
				}
			} else {
				// First time data was saved
				result.success = true
				log.Printf("Initial sync of %d delegations for validator %s", 
					result.delegationCount, address)
			}
			
			resultCh <- result
		}(validatorAddress)
	}
	
	// Wait for all goroutines to complete
	wg.Wait()
	close(resultCh)
	
	// Process all results
	for result := range resultCh {
		if result.error != nil {
			stats.ErrorCount++
		} else if result.skipped {
			stats.SkippedCount++
		} else if result.success {
			stats.SuccessCount++
		}
		stats.TotalDelegationsProcessed += result.delegationCount
	}
	
	stats.LastRunDuration = time.Since(startTime)
	stats.TotalRuns = t.lastRunStats.TotalRuns + 1
	t.lastRunStats = stats
	t.totalDelegationsSynced += stats.TotalDelegationsProcessed
	
	log.Printf("Delegations sync completed in %v. Success: %d, Skipped: %d, Errors: %d, Total delegations: %d", 
		stats.LastRunDuration, stats.SuccessCount, stats.SkippedCount, stats.ErrorCount, stats.TotalDelegationsProcessed)
	
	if stats.ErrorCount > 0 {
		return fmt.Errorf("completed with %d errors (and %d successes, %d skipped)", 
			stats.ErrorCount, stats.SuccessCount, stats.SkippedCount)
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