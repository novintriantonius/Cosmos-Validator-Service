package store

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/novintriantonius/cosmos-validator-service/internal/models"
)

// DelegationStore defines the interface for delegation storage
type DelegationStore interface {
	// SaveDelegations saves delegations for a validator
	SaveDelegations(validatorAddress string, data models.DelegationsResponse) error
	
	// GetDelegations retrieves delegations for a validator
	GetDelegations(validatorAddress string) (*models.StoredDelegationsData, error)
	
	// GetAllDelegations retrieves all stored delegations
	GetAllDelegations() (map[string]*models.StoredDelegationsData, error)
	
	// EnableDelegationTracking enables tracking for a validator
	EnableDelegationTracking(validatorAddress string) error
	
	// DisableDelegationTracking disables tracking for a validator
	DisableDelegationTracking(validatorAddress string) error
	
	// GetEnabledValidators gets all validators with enabled tracking
	GetEnabledValidators() ([]string, error)
}

// InMemoryDelegationStore implements the DelegationStore interface with in-memory storage
type InMemoryDelegationStore struct {
	delegations map[string]*models.StoredDelegationsData
	mu          sync.RWMutex
}

// NewInMemoryDelegationStore creates a new in-memory delegation store
func NewInMemoryDelegationStore() *InMemoryDelegationStore {
	return &InMemoryDelegationStore{
		delegations: make(map[string]*models.StoredDelegationsData),
	}
}

// SaveDelegations saves delegations for a validator if the data has changed
func (s *InMemoryDelegationStore) SaveDelegations(validatorAddress string, data models.DelegationsResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Check if we already have delegations for this validator
	stored, exists := s.delegations[validatorAddress]
	
	if exists {
		// Check if the data has actually changed
		if !hasDelegationsChanged(stored.Data, data) {
			// Data hasn't changed, only update the last checked timestamp
			stored.Timestamp = time.Now()
			return nil
		}
		
		// Update existing record because data has changed
		stored.Data = data
		stored.Timestamp = data.Timestamp
	} else {
		// Create new record with tracking enabled by default
		stored = &models.StoredDelegationsData{
			ValidatorAddress: validatorAddress,
			Data:             data,
			Timestamp:        data.Timestamp,
			IsEnabled:        true,
		}
		s.delegations[validatorAddress] = stored
	}
	
	return nil
}

// hasDelegationsChanged compares two DelegationsResponse objects to check if there are actual changes
// Returns true if the delegations data has changed
func hasDelegationsChanged(old, new models.DelegationsResponse) bool {
	// If the number of delegations has changed, data has changed
	if len(old.DelegationResponses) != len(new.DelegationResponses) {
		return true
	}
	
	// If pagination data has changed, consider it changed
	if old.Pagination.Total != new.Pagination.Total || old.Pagination.NextKey != new.Pagination.NextKey {
		return true
	}
	
	// Convert old delegations to a map for faster lookup
	oldDelegations := make(map[string]models.DelegationResponse, len(old.DelegationResponses))
	for _, delegation := range old.DelegationResponses {
		// Use delegator address as a key since it should be unique for a validator
		oldDelegations[delegation.Delegation.DelegatorAddress] = delegation
	}
	
	// Check if any delegation has changed
	for _, newDelegation := range new.DelegationResponses {
		delegatorAddr := newDelegation.Delegation.DelegatorAddress
		oldDelegation, exists := oldDelegations[delegatorAddr]
		
		if !exists {
			// This is a new delegation
			return true
		}
		
		// Check if any field has changed
		if !reflect.DeepEqual(oldDelegation, newDelegation) {
			return true
		}
	}
	
	return false
}

// GetDelegations retrieves delegations for a validator
func (s *InMemoryDelegationStore) GetDelegations(validatorAddress string) (*models.StoredDelegationsData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	delegations, exists := s.delegations[validatorAddress]
	if !exists {
		return nil, fmt.Errorf("no delegations found for validator %s", validatorAddress)
	}
	
	return delegations, nil
}

// GetAllDelegations retrieves all stored delegations
func (s *InMemoryDelegationStore) GetAllDelegations() (map[string]*models.StoredDelegationsData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Create a copy to avoid race conditions
	result := make(map[string]*models.StoredDelegationsData, len(s.delegations))
	for k, v := range s.delegations {
		result[k] = v
	}
	
	return result, nil
}

// EnableDelegationTracking enables tracking for a validator
func (s *InMemoryDelegationStore) EnableDelegationTracking(validatorAddress string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	delegations, exists := s.delegations[validatorAddress]
	if !exists {
		// Create an empty record with tracking enabled
		delegations = &models.StoredDelegationsData{
			ValidatorAddress: validatorAddress,
			IsEnabled:        true,
		}
		s.delegations[validatorAddress] = delegations
	} else {
		delegations.IsEnabled = true
	}
	
	return nil
}

// DisableDelegationTracking disables tracking for a validator
func (s *InMemoryDelegationStore) DisableDelegationTracking(validatorAddress string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	delegations, exists := s.delegations[validatorAddress]
	if !exists {
		return fmt.Errorf("no delegations found for validator %s", validatorAddress)
	}
	
	delegations.IsEnabled = false
	return nil
}

// GetEnabledValidators gets all validators with enabled tracking
func (s *InMemoryDelegationStore) GetEnabledValidators() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var enabledValidators []string
	for address, delegations := range s.delegations {
		if delegations.IsEnabled {
			enabledValidators = append(enabledValidators, address)
		}
	}
	
	return enabledValidators, nil
} 