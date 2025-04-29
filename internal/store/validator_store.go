package store

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/novintriantonius/cosmos-validator-service/internal/models"
)

var (
	// ErrValidatorNotFound is returned when a validator is not found in the store
	ErrValidatorNotFound = errors.New("validator not found")
	
	// ErrValidatorAlreadyExists is returned when trying to add a validator with an address that already exists
	ErrValidatorAlreadyExists = errors.New("validator with this address already exists")
)

// ValidatorStore defines the interface for validator storage operations
type ValidatorStore interface {
	GetAll() ([]models.Validator, error)
	GetByAddress(address string) (*models.Validator, error)
	GetEnabledValidators() ([]string, error)
	Add(validator models.Validator) error
	Update(address string, validator models.Validator) error
	Delete(address string) error
}

// ValidatorStoreImpl implements ValidatorStore with PostgreSQL storage
type ValidatorStoreImpl struct {
	db *sql.DB
	mu sync.RWMutex
}

// NewValidatorStore creates a new instance of ValidatorStoreImpl
func NewValidatorStore(db *sql.DB) *ValidatorStoreImpl {
	return &ValidatorStoreImpl{
		db: db,
	}
}

// GetAll returns all validators from the database
func (s *ValidatorStoreImpl) GetAll() ([]models.Validator, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT address, name, enabled_tracking FROM validators`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying validators: %v", err)
	}
	defer rows.Close()

	var validators []models.Validator
	for rows.Next() {
		var v models.Validator
		if err := rows.Scan(&v.Address, &v.Name, &v.EnabledTracking); err != nil {
			return nil, fmt.Errorf("error scanning validator row: %v", err)
		}
		validators = append(validators, v)
	}

	return validators, nil
}

// GetByAddress returns a validator by its address
func (s *ValidatorStoreImpl) GetByAddress(address string) (*models.Validator, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT address, name, enabled_tracking FROM validators WHERE address = $1`
	var v models.Validator
	err := s.db.QueryRow(query, address).Scan(&v.Address, &v.Name, &v.EnabledTracking)
	if err == sql.ErrNoRows {
		return nil, ErrValidatorNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error querying validator: %v", err)
	}
	return &v, nil
}

// Add adds a new validator to the database
func (s *ValidatorStoreImpl) Add(validator models.Validator) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `
		INSERT INTO validators (address, name, enabled_tracking)
		VALUES ($1, $2, $3)
	`
	_, err := s.db.Exec(query, validator.Address, validator.Name, validator.EnabledTracking)
	if err != nil {
		return fmt.Errorf("error inserting validator: %v", err)
	}
	return nil
}

// Update updates an existing validator in the database
func (s *ValidatorStoreImpl) Update(address string, validator models.Validator) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `
		UPDATE validators
		SET name = $1, enabled_tracking = $2, updated_at = CURRENT_TIMESTAMP
		WHERE address = $3
	`
	result, err := s.db.Exec(query, validator.Name, validator.EnabledTracking, address)
	if err != nil {
		return fmt.Errorf("error updating validator: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return ErrValidatorNotFound
	}

	return nil
}

// Delete removes a validator from the database
func (s *ValidatorStoreImpl) Delete(address string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `DELETE FROM validators WHERE address = $1`
	result, err := s.db.Exec(query, address)
	if err != nil {
		return fmt.Errorf("error deleting validator: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return ErrValidatorNotFound
	}

	return nil
}

// GetEnabledValidators returns a list of validator addresses that have enabled tracking
func (s *ValidatorStoreImpl) GetEnabledValidators() ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `
		SELECT address 
		FROM validators 
		WHERE enabled_tracking = true
	`

	rows, err := s.db.Query(query)
	if err != nil {
		log.Printf("[ERROR] Failed to query enabled validators: %v", err)
		return nil, fmt.Errorf("error querying enabled validators: %v", err)
	}
	defer rows.Close()

	var addresses []string
	for rows.Next() {
		var address string
		if err := rows.Scan(&address); err != nil {
			log.Printf("[ERROR] Failed to scan validator row: %v", err)
			return nil, fmt.Errorf("error scanning validator row: %v", err)
		}
		addresses = append(addresses, address)
	}

	log.Printf("[DEBUG] Found %d enabled validators", len(addresses))
	return addresses, nil
} 