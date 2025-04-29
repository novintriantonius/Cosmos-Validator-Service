package store_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/novintriantonius/cosmos-validator-service/internal/models"
	"github.com/novintriantonius/cosmos-validator-service/internal/store"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestValidatorStore_GetAll(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	store := store.NewValidatorStore(db)

	// Mock rows
	rows := sqlmock.NewRows([]string{"address", "name", "enabled_tracking"}).
		AddRow("val1", "Validator 1", true).
		AddRow("val2", "Validator 2", false)

	mock.ExpectQuery("SELECT address, name, enabled_tracking FROM validators").
		WillReturnRows(rows)

	validators, err := store.GetAll()
	assert.NoError(t, err)
	assert.Len(t, validators, 2)
	assert.Equal(t, "val1", validators[0].Address)
	assert.Equal(t, "Validator 1", validators[0].Name)
	assert.True(t, validators[0].EnabledTracking)
}

func TestValidatorStore_GetByAddress(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	store := store.NewValidatorStore(db)

	// Test case: Validator found
	rows := sqlmock.NewRows([]string{"address", "name", "enabled_tracking"}).
		AddRow("val1", "Validator 1", true)

	mock.ExpectQuery("SELECT address, name, enabled_tracking FROM validators WHERE address = \\$1").
		WithArgs("val1").
		WillReturnRows(rows)

	validator, err := store.GetByAddress("val1")
	assert.NoError(t, err)
	assert.Equal(t, "val1", validator.Address)
	assert.Equal(t, "Validator 1", validator.Name)
	assert.True(t, validator.EnabledTracking)

	// Test case: Validator not found
	mock.ExpectQuery("SELECT address, name, enabled_tracking FROM validators WHERE address = \\$1").
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	_, err = store.GetByAddress("nonexistent")
	assert.Error(t, err)
	assert.Equal(t, "validator not found", err.Error())
}

func TestDelegationStore_SaveDelegations(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	delegationStore := store.NewDelegationStore(db)

	// Prepare test data for updated model structure
	delegationsResponse := models.DelegationsResponse{
		DelegationResponses: []models.DelegationResponse{
			{
				Delegation: models.DelegationDetails{
					DelegatorAddress: "delegator1",
					ValidatorAddress: "validator1",
					Shares:           "100.0",
				},
				Balance: models.Balance{
					Denom:  "uatom",
					Amount: "100",
				},
			},
		},
		Pagination: models.Pagination{
			NextKey: "next",
			Total:   "1",
		},
	}

	// Mock the database behavior
	mock.ExpectBegin()
	
	// Prepare statement mock
	mock.ExpectPrepare("INSERT INTO delegations").WillBeClosed()
	
	// Query for existing delegations
	rows := sqlmock.NewRows([]string{"delegator_address", "delegation_shares"})
	mock.ExpectQuery("SELECT DISTINCT ON \\(delegator_address\\)").WithArgs("validator1").WillReturnRows(rows)
	
	// Execute the insert
	mock.ExpectExec("INSERT INTO delegations").WithArgs("validator1", "delegator1", "100.0").WillReturnResult(sqlmock.NewResult(1, 1))
	
	// Commit transaction
	mock.ExpectCommit()

	err := delegationStore.SaveDelegations("validator1", delegationsResponse)
	assert.NoError(t, err)
}

func TestDelegationStore_GetDelegations(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	delegationStore := store.NewDelegationStore(db)

	// Setup current time for test
	createdAt := time.Now()
	updatedAt := createdAt

	// Mock the rows returned by the query with actual delegation model structure
	rows := sqlmock.NewRows([]string{"id", "validator_address", "delegator_address", "delegation_shares", "created_at", "updated_at"}).
		AddRow(1, "validator1", "delegator1", "100.0", createdAt, updatedAt)

	mock.ExpectQuery("SELECT id, validator_address, delegator_address, delegation_shares, created_at, updated_at FROM delegations WHERE validator_address = \\$1").
		WithArgs("validator1").
		WillReturnRows(rows)

	delegations, err := delegationStore.GetDelegations("validator1")
	assert.NoError(t, err)
	assert.Len(t, delegations, 1)
	assert.Equal(t, "validator1", delegations[0].ValidatorAddress)
	assert.Equal(t, "delegator1", delegations[0].DelegatorAddress)
	assert.Equal(t, "100.0", delegations[0].DelegationShares)
}

func TestDelegationStore_GetEnabledValidators(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	delegationStore := store.NewDelegationStore(db)

	rows := sqlmock.NewRows([]string{"address"}).
		AddRow("validator1").
		AddRow("validator2")

	mock.ExpectQuery("SELECT address FROM validators WHERE enabled_tracking = true").
		WillReturnRows(rows)

	validators, err := delegationStore.GetEnabledValidators()
	assert.NoError(t, err)
	assert.Len(t, validators, 2)
	assert.Contains(t, validators, "validator1")
	assert.Contains(t, validators, "validator2")
}

func TestDelegationStore_DelegationExists(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	delegationStore := store.NewDelegationStore(db)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("validator1", "delegator1", "100.0").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := delegationStore.DelegationExists("validator1", "delegator1", "100.0")
	assert.NoError(t, err)
	assert.True(t, exists)
} 