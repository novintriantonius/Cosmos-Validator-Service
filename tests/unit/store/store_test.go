package store_test

import (
	"database/sql"
	"encoding/json"
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

	store := store.NewDelegationStore(db)

	// Prepare test data
	delegations := models.DelegationsResponse{
		DelegationResponses: []models.DelegationResponse{
			{
				Delegation: models.Delegation{
					DelegatorAddress: "delegator1",
					ValidatorAddress: "validator1",
					Shares:          "100.0",
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

	jsonData, _ := json.Marshal(delegations)

	mock.ExpectExec("INSERT INTO delegations").
		WithArgs("validator1", jsonData).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := store.SaveDelegations("validator1", delegations)
	assert.NoError(t, err)
}

func TestDelegationStore_GetDelegations(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	store := store.NewDelegationStore(db)

	// Prepare test data
	delegations := models.DelegationsResponse{
		DelegationResponses: []models.DelegationResponse{
			{
				Delegation: models.Delegation{
					DelegatorAddress: "delegator1",
					ValidatorAddress: "validator1",
					Shares:          "100.0",
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

	jsonData, _ := json.Marshal(delegations)

	rows := sqlmock.NewRows([]string{"data", "is_enabled"}).
		AddRow(jsonData, true)

	mock.ExpectQuery("SELECT data, is_enabled FROM delegations WHERE validator_address = \\$1").
		WithArgs("validator1").
		WillReturnRows(rows)

	stored, err := store.GetDelegations("validator1")
	assert.NoError(t, err)
	assert.Equal(t, "validator1", stored.ValidatorAddress)
	assert.True(t, stored.IsEnabled)
	assert.Len(t, stored.Data.DelegationResponses, 1)
}

func TestDelegationStore_EnableDelegationTracking(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	store := store.NewDelegationStore(db)

	// The query expects 4 parameters: validator_address, is_enabled, data, and timestamp
	mock.ExpectExec("^INSERT INTO delegations \\(validator_address, is_enabled, data, timestamp\\) VALUES \\(\\$1, true, '\\{\\}', \\$2\\) ON CONFLICT \\(validator_address\\) DO UPDATE SET is_enabled = true, updated_at = CURRENT_TIMESTAMP$").
		WithArgs("validator1", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := store.EnableDelegationTracking("validator1")
	assert.NoError(t, err)
}

func TestDelegationStore_DisableDelegationTracking(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	store := store.NewDelegationStore(db)

	mock.ExpectExec("UPDATE delegations").
		WithArgs("validator1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := store.DisableDelegationTracking("validator1")
	assert.NoError(t, err)
}

func TestDelegationStore_GetEnabledValidators(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	store := store.NewDelegationStore(db)

	rows := sqlmock.NewRows([]string{"validator_address"}).
		AddRow("validator1").
		AddRow("validator2")

	mock.ExpectQuery("SELECT validator_address FROM delegations WHERE is_enabled = true").
		WillReturnRows(rows)

	validators, err := store.GetEnabledValidators()
	assert.NoError(t, err)
	assert.Len(t, validators, 2)
	assert.Contains(t, validators, "validator1")
	assert.Contains(t, validators, "validator2")
} 