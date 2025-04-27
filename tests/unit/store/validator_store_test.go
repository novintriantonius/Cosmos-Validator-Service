package store_test

import (
	"testing"

	"github.com/novintriantonius/cosmos-validator-service/internal/models"
	"github.com/novintriantonius/cosmos-validator-service/internal/store"
)

func TestNewInMemoryValidatorStore(t *testing.T) {
	s := store.NewInMemoryValidatorStore()
	if s == nil {
		t.Fatal("Expected non-nil store")
	}
}

func TestAddValidator(t *testing.T) {
	s := store.NewInMemoryValidatorStore()
	
	validator := models.Validator{
		Name:            "Test Validator",
		Address:         "cosmosvaloper1testaddress",
		EnabledTracking: true,
	}
	
	// Test adding a validator
	err := s.Add(validator)
	if err != nil {
		t.Fatalf("Failed to add validator: %v", err)
	}
	
	// Test adding the same validator again (should error)
	err = s.Add(validator)
	if err != store.ErrValidatorAlreadyExists {
		t.Errorf("Expected ErrValidatorAlreadyExists, got %v", err)
	}
}

func TestGetAllValidators(t *testing.T) {
	s := store.NewInMemoryValidatorStore()
	
	// Test empty store
	validators, err := s.GetAll()
	if err != nil {
		t.Fatalf("Failed to get validators: %v", err)
	}
	if len(validators) != 0 {
		t.Errorf("Expected 0 validators, got %d", len(validators))
	}
	
	// Add some validators
	validator1 := models.Validator{
		Name:            "Test Validator 1",
		Address:         "cosmosvaloper1testaddress1",
		EnabledTracking: true,
	}
	
	validator2 := models.Validator{
		Name:            "Test Validator 2",
		Address:         "cosmosvaloper1testaddress2",
		EnabledTracking: false,
	}
	
	s.Add(validator1)
	s.Add(validator2)
	
	// Test getting all validators
	validators, err = s.GetAll()
	if err != nil {
		t.Fatalf("Failed to get validators: %v", err)
	}
	if len(validators) != 2 {
		t.Errorf("Expected 2 validators, got %d", len(validators))
	}
}

func TestGetValidatorByAddress(t *testing.T) {
	s := store.NewInMemoryValidatorStore()
	
	// Test getting a non-existent validator
	_, err := s.GetByAddress("nonexistent")
	if err != store.ErrValidatorNotFound {
		t.Errorf("Expected ErrValidatorNotFound, got %v", err)
	}
	
	// Add a validator
	validator := models.Validator{
		Name:            "Test Validator",
		Address:         "cosmosvaloper1testaddress",
		EnabledTracking: true,
	}
	
	s.Add(validator)
	
	// Test getting the validator by address
	retrieved, err := s.GetByAddress(validator.Address)
	if err != nil {
		t.Fatalf("Failed to get validator: %v", err)
	}
	
	// Check that the retrieved validator matches the original
	if retrieved.Name != validator.Name || retrieved.Address != validator.Address || retrieved.EnabledTracking != validator.EnabledTracking {
		t.Errorf("Retrieved validator does not match original: %+v vs %+v", retrieved, validator)
	}
}

func TestUpdateValidator(t *testing.T) {
	s := store.NewInMemoryValidatorStore()
	
	// Test updating a non-existent validator
	err := s.Update("nonexistent", models.Validator{})
	if err != store.ErrValidatorNotFound {
		t.Errorf("Expected ErrValidatorNotFound, got %v", err)
	}
	
	// Add a validator
	validator := models.Validator{
		Name:            "Test Validator",
		Address:         "cosmosvaloper1testaddress",
		EnabledTracking: true,
	}
	
	s.Add(validator)
	
	// Update the validator
	updatedValidator := models.Validator{
		Name:            "Updated Test Validator",
		Address:         "cosmosvaloper1testaddress", // Same address
		EnabledTracking: false,
	}
	
	err = s.Update(validator.Address, updatedValidator)
	if err != nil {
		t.Fatalf("Failed to update validator: %v", err)
	}
	
	// Check that the validator was updated
	retrieved, _ := s.GetByAddress(validator.Address)
	if retrieved.Name != updatedValidator.Name || retrieved.EnabledTracking != updatedValidator.EnabledTracking {
		t.Errorf("Validator was not updated correctly: %+v", retrieved)
	}
	
	// Check that the address didn't change
	if retrieved.Address != validator.Address {
		t.Errorf("Validator address changed from %s to %s", validator.Address, retrieved.Address)
	}
}

func TestDeleteValidator(t *testing.T) {
	s := store.NewInMemoryValidatorStore()
	
	// Test deleting a non-existent validator
	err := s.Delete("nonexistent")
	if err != store.ErrValidatorNotFound {
		t.Errorf("Expected ErrValidatorNotFound, got %v", err)
	}
	
	// Add a validator
	validator := models.Validator{
		Name:            "Test Validator",
		Address:         "cosmosvaloper1testaddress",
		EnabledTracking: true,
	}
	
	s.Add(validator)
	
	// Delete the validator
	err = s.Delete(validator.Address)
	if err != nil {
		t.Fatalf("Failed to delete validator: %v", err)
	}
	
	// Check that the validator was deleted
	_, err = s.GetByAddress(validator.Address)
	if err != store.ErrValidatorNotFound {
		t.Errorf("Expected validator to be deleted, but got %v", err)
	}
}
