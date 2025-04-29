package store

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/novintriantonius/cosmos-validator-service/internal/models"
)

// DelegationStore defines the interface for delegation storage
type DelegationStore interface {
	// SaveDelegations saves delegations for a validator
	SaveDelegations(validatorAddress string, data models.DelegationsResponse) error
	
	// GetDelegations retrieves delegations for a validator
	GetDelegations(validatorAddress string) ([]models.Delegation, error)
	
	// GetAllDelegations retrieves all stored delegations
	GetAllDelegations() (map[string][]models.Delegation, error)
	
	// EnableDelegationTracking enables tracking for a validator
	EnableDelegationTracking(validatorAddress string) error
	
	// DisableDelegationTracking disables tracking for a validator
	DisableDelegationTracking(validatorAddress string) error
	
	// GetEnabledValidators gets all validators with enabled tracking
	GetEnabledValidators() ([]string, error)

	// DelegationExists checks if a delegation exists for the given validator, delegator, and shares
	DelegationExists(validatorAddress, delegatorAddress, delegationShares string) (bool, error)
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

	log.Printf("[DEBUG] Starting SaveDelegations for validator %s with %d delegations", 
		validatorAddress, len(data.DelegationResponses))

	// Start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("[ERROR] Failed to start transaction: %v", err)
		return fmt.Errorf("error starting transaction: %v", err)
	}
	log.Printf("[DEBUG] Transaction started successfully")

	// Prepare the insert statement
	stmt, err := tx.Prepare(`
		INSERT INTO delegations (validator_address, delegator_address, delegation_shares)
		VALUES ($1, $2, $3)
	`)
	if err != nil {
		log.Printf("[ERROR] Failed to prepare insert statement: %v", err)
		tx.Rollback()
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()
	log.Printf("[DEBUG] Insert statement prepared successfully")

	// Get latest delegations for this validator
	log.Printf("[DEBUG] Querying latest delegations for validator %s", validatorAddress)
	latestDelegations := make(map[string]string) // delegator_address -> shares
	rows, err := tx.Query(`
		SELECT DISTINCT ON (delegator_address) delegator_address, delegation_shares
		FROM delegations
		WHERE validator_address = $1
		ORDER BY delegator_address, created_at DESC
	`, validatorAddress)
	if err != nil {
		log.Printf("[ERROR] Failed to query latest delegations: %v", err)
		tx.Rollback()
		return fmt.Errorf("error querying latest delegations: %v", err)
	}
	defer rows.Close()

	existingCount := 0
	for rows.Next() {
		var delegatorAddress, shares string
		if err := rows.Scan(&delegatorAddress, &shares); err != nil {
			log.Printf("[ERROR] Failed to scan delegation row: %v", err)
			tx.Rollback()
			return fmt.Errorf("error scanning delegation row: %v", err)
		}
		latestDelegations[delegatorAddress] = shares
		existingCount++
	}
	log.Printf("[DEBUG] Found %d existing delegations for validator %s", existingCount, validatorAddress)

	// Insert each delegation if shares have changed
	successCount := 0
	skippedCount := 0
	for i, resp := range data.DelegationResponses {
		delegatorAddress := resp.Delegation.DelegatorAddress
		newShares := resp.Delegation.Shares

		log.Printf("[DEBUG] Processing delegation %d/%d: delegator=%s, shares=%s", 
			i+1, len(data.DelegationResponses), delegatorAddress, newShares)

		// Check if we have a previous delegation for this delegator
		if existingShares, exists := latestDelegations[delegatorAddress]; exists {
			// Skip if shares haven't changed
			if existingShares == newShares {
				log.Printf("[DEBUG] Skipping delegation for delegator %s - shares unchanged (existing=%s, new=%s)", 
					delegatorAddress, existingShares, newShares)
				skippedCount++
				continue
			}
			log.Printf("[DEBUG] Shares changed for delegator %s: old=%s, new=%s", 
				delegatorAddress, existingShares, newShares)
		} else {
			log.Printf("[DEBUG] New delegator %s with shares %s", delegatorAddress, newShares)
		}

		// Insert new delegation
		_, err := stmt.Exec(
			validatorAddress,
			delegatorAddress,
			newShares,
		)
		if err != nil {
			log.Printf("[ERROR] Failed to insert delegation %d for delegator %s: %v", 
				i, delegatorAddress, err)
			tx.Rollback()
			return fmt.Errorf("error inserting delegation: %v", err)
		}
		successCount++
		log.Printf("[DEBUG] Successfully inserted delegation for delegator %s", delegatorAddress)
	}

	log.Printf("[INFO] Delegation processing complete for validator %s: %d processed, %d skipped, %d successful", 
		validatorAddress, len(data.DelegationResponses), skippedCount, successCount)

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		log.Printf("[ERROR] Failed to commit transaction: %v", err)
		return fmt.Errorf("error committing transaction: %v", err)
	}

	log.Printf("[INFO] Successfully committed transaction for validator %s", validatorAddress)
	return nil
}

// GetDelegations retrieves delegations for a validator
func (s *DelegationStoreImpl) GetDelegations(validatorAddress string) ([]models.Delegation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `
		SELECT id, validator_address, delegator_address, delegation_shares, created_at, updated_at
		FROM delegations
		WHERE validator_address = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, validatorAddress)
	if err != nil {
		return nil, fmt.Errorf("error querying delegations: %v", err)
	}
	defer rows.Close()

	var delegations []models.Delegation
	for rows.Next() {
		var d models.Delegation
		err := rows.Scan(
			&d.ID,
			&d.ValidatorAddress,
			&d.DelegatorAddress,
			&d.DelegationShares,
			&d.CreatedAt,
			&d.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning delegation row: %v", err)
		}
		delegations = append(delegations, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating delegation rows: %v", err)
	}

	return delegations, nil
}

// GetAllDelegations retrieves all stored delegations
func (s *DelegationStoreImpl) GetAllDelegations() (map[string][]models.Delegation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `
		SELECT id, validator_address, delegator_address, delegation_shares, created_at, updated_at
		FROM delegations
		ORDER BY validator_address, created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying all delegations: %v", err)
	}
	defer rows.Close()

	delegations := make(map[string][]models.Delegation)
	for rows.Next() {
		var d models.Delegation
		err := rows.Scan(
			&d.ID,
			&d.ValidatorAddress,
			&d.DelegatorAddress,
			&d.DelegationShares,
			&d.CreatedAt,
			&d.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning delegation row: %v", err)
		}
		delegations[d.ValidatorAddress] = append(delegations[d.ValidatorAddress], d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating delegation rows: %v", err)
	}

	return delegations, nil
}

// EnableDelegationTracking enables tracking for a validator
func (s *DelegationStoreImpl) EnableDelegationTracking(validatorAddress string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// No need to do anything as tracking is implicit in the presence of delegations
	return nil
}

// DisableDelegationTracking disables tracking for a validator
func (s *DelegationStoreImpl) DisableDelegationTracking(validatorAddress string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Delete all delegations for the validator
	query := `DELETE FROM delegations WHERE validator_address = $1`
	result, err := s.db.Exec(query, validatorAddress)
	if err != nil {
		return fmt.Errorf("error deleting delegations: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no delegations found for validator %s", validatorAddress)
	}

	return nil
}

// GetEnabledValidators gets all validators with enabled tracking
func (s *DelegationStoreImpl) GetEnabledValidators() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `
		SELECT address
		FROM validators
		WHERE enabled_tracking = true
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

// DelegationExists checks if a delegation exists for the given validator, delegator, and shares
func (s *DelegationStoreImpl) DelegationExists(validatorAddress, delegatorAddress, delegationShares string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM delegations
			WHERE validator_address = $1
			AND delegator_address = $2
			AND delegation_shares = $3
		)
	`

	var exists bool
	err := s.db.QueryRow(query, validatorAddress, delegatorAddress, delegationShares).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking delegation existence: %v", err)
	}

	return exists, nil
} 