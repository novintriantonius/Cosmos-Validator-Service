package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

// DelegationStoreImpl implements the DelegationStore interface with PostgreSQL storage
type DelegationStoreImpl struct {
	db *sql.DB
	mu sync.RWMutex
}

// NewDelegationStore creates a new instance of DelegationStoreImpl
func NewDelegationStore(db *sql.DB) *DelegationStoreImpl {
	return &DelegationStoreImpl{
		db: db,
	}
}

// SaveDelegations saves delegations for a validator
func (s *DelegationStoreImpl) SaveDelegations(validatorAddress string, data models.DelegationsResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling delegations data: %v", err)
	}

	query := `
		INSERT INTO delegations (validator_address, data, timestamp, is_enabled)
		VALUES ($1, $2, $3, true)
		ON CONFLICT (validator_address)
		DO UPDATE SET
			data = $2,
			timestamp = $3,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err = s.db.Exec(query, validatorAddress, jsonData, data.Timestamp)
	if err != nil {
		return fmt.Errorf("error saving delegations: %v", err)
	}

	return nil
}

// GetDelegations retrieves delegations for a validator
func (s *DelegationStoreImpl) GetDelegations(validatorAddress string) (*models.StoredDelegationsData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `
		SELECT data, timestamp, is_enabled
		FROM delegations
		WHERE validator_address = $1
	`

	var (
		jsonData  []byte
		stored    models.StoredDelegationsData
		timestamp sql.NullTime
	)

	err := s.db.QueryRow(query, validatorAddress).Scan(&jsonData, &timestamp, &stored.IsEnabled)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no delegations found for validator %s", validatorAddress)
	}
	if err != nil {
		return nil, fmt.Errorf("error querying delegations: %v", err)
	}

	// Parse JSON data
	if err := json.Unmarshal(jsonData, &stored.Data); err != nil {
		return nil, fmt.Errorf("error unmarshaling delegations data: %v", err)
	}

	stored.ValidatorAddress = validatorAddress
	if timestamp.Valid {
		stored.Timestamp = timestamp.Time
	}

	return &stored, nil
}

// GetAllDelegations retrieves all stored delegations
func (s *DelegationStoreImpl) GetAllDelegations() (map[string]*models.StoredDelegationsData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `
		SELECT validator_address, data, timestamp, is_enabled
		FROM delegations
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying all delegations: %v", err)
	}
	defer rows.Close()

	delegations := make(map[string]*models.StoredDelegationsData)
	for rows.Next() {
		var (
			stored    models.StoredDelegationsData
			jsonData  []byte
			timestamp sql.NullTime
		)

		err := rows.Scan(&stored.ValidatorAddress, &jsonData, &timestamp, &stored.IsEnabled)
		if err != nil {
			return nil, fmt.Errorf("error scanning delegation row: %v", err)
		}

		if err := json.Unmarshal(jsonData, &stored.Data); err != nil {
			return nil, fmt.Errorf("error unmarshaling delegations data: %v", err)
		}

		if timestamp.Valid {
			stored.Timestamp = timestamp.Time
		}

		delegations[stored.ValidatorAddress] = &stored
	}

	return delegations, nil
}

// EnableDelegationTracking enables tracking for a validator
func (s *DelegationStoreImpl) EnableDelegationTracking(validatorAddress string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `
		INSERT INTO delegations (validator_address, is_enabled, data, timestamp)
		VALUES ($1, true, '{}', $2)
		ON CONFLICT (validator_address)
		DO UPDATE SET
			is_enabled = true,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := s.db.Exec(query, validatorAddress, time.Now())
	if err != nil {
		return fmt.Errorf("error enabling delegation tracking: %v", err)
	}

	return nil
}

// DisableDelegationTracking disables tracking for a validator
func (s *DelegationStoreImpl) DisableDelegationTracking(validatorAddress string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `
		UPDATE delegations
		SET is_enabled = false, updated_at = CURRENT_TIMESTAMP
		WHERE validator_address = $1
	`

	result, err := s.db.Exec(query, validatorAddress)
	if err != nil {
		return fmt.Errorf("error disabling delegation tracking: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no validator found with address %s", validatorAddress)
	}

	return nil
}

// GetEnabledValidators gets all validators with enabled tracking
func (s *DelegationStoreImpl) GetEnabledValidators() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `
		SELECT validator_address
		FROM delegations
		WHERE is_enabled = true
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying enabled validators: %v", err)
	}
	defer rows.Close()

	var validators []string
	for rows.Next() {
		var address string
		if err := rows.Scan(&address); err != nil {
			return nil, fmt.Errorf("error scanning validator address: %v", err)
		}
		validators = append(validators, address)
	}

	return validators, nil
} 