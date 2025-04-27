package store

import (
	"errors"
	"sync"

	"github.com/novintriantonius/cosmos-validator-service/internal/models"
)

var (
	// ErrValidatorNotFound is returned when a validator is not found in the store
	ErrValidatorNotFound = errors.New("validator not found")
	
	// ErrValidatorAlreadyExists is returned when trying to add a validator with an address that already exists
	ErrValidatorAlreadyExists = errors.New("validator with this address already exists")
)

// ValidatorStore provides an interface for validator data operations
type ValidatorStore interface {
	GetAll() ([]models.Validator, error)
	GetByAddress(address string) (models.Validator, error)
	Add(validator models.Validator) error
	Update(address string, validator models.Validator) error
	Delete(address string) error
}

// InMemoryValidatorStore implements ValidatorStore with an in-memory storage
type InMemoryValidatorStore struct {
	validators map[string]models.Validator
	mu         sync.RWMutex
}

// NewInMemoryValidatorStore creates a new instance of InMemoryValidatorStore
func NewInMemoryValidatorStore() *InMemoryValidatorStore {
	return &InMemoryValidatorStore{
		validators: make(map[string]models.Validator),
	}
}

// GetAll returns all validators in the store
func (s *InMemoryValidatorStore) GetAll() ([]models.Validator, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	validators := make([]models.Validator, 0, len(s.validators))
	for _, v := range s.validators {
		validators = append(validators, v)
	}
	
	return validators, nil
}

// GetByAddress returns a validator by its address
func (s *InMemoryValidatorStore) GetByAddress(address string) (models.Validator, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	validator, exists := s.validators[address]
	if !exists {
		return models.Validator{}, ErrValidatorNotFound
	}
	
	return validator, nil
}

// Add adds a new validator to the store
func (s *InMemoryValidatorStore) Add(validator models.Validator) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.validators[validator.Address]; exists {
		return ErrValidatorAlreadyExists
	}
	
	s.validators[validator.Address] = validator
	return nil
}

// Update updates an existing validator in the store
func (s *InMemoryValidatorStore) Update(address string, validator models.Validator) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.validators[address]; !exists {
		return ErrValidatorNotFound
	}
	
	// Ensure the address in the updated validator remains the same
	validator.Address = address
	s.validators[address] = validator
	return nil
}

// Delete removes a validator from the store
func (s *InMemoryValidatorStore) Delete(address string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.validators[address]; !exists {
		return ErrValidatorNotFound
	}
	
	delete(s.validators, address)
	return nil
} 