package models

import (
	"time"
)

// Delegation represents a single delegation entry in our database
type Delegation struct {
	ID                int       `json:"id"`
	ValidatorAddress  string    `json:"validator_address"`
	DelegatorAddress  string    `json:"delegator_address"`
	DelegationShares  string    `json:"delegation_shares"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// DelegationsResponse represents the response from the delegations API
type DelegationsResponse struct {
	DelegationResponses []DelegationResponse `json:"delegation_responses"`
	Pagination          Pagination           `json:"pagination"`
}

// DelegationResponse represents a single delegation response from the API
type DelegationResponse struct {
	Delegation DelegationDetails `json:"delegation"`
	Balance    Balance          `json:"balance"`
}

// DelegationDetails represents the delegation details from the API
type DelegationDetails struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Shares           string `json:"shares"`
}

// Balance represents the token balance
type Balance struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

// Pagination represents pagination information
type Pagination struct {
	NextKey string `json:"next_key"`
	Total   string `json:"total"`
}

// StoredDelegationsData represents delegation data stored for a validator
type StoredDelegationsData struct {
	ValidatorAddress string              `json:"validator_address"`
	Data             DelegationsResponse `json:"data"`
	Timestamp        time.Time           `json:"timestamp"`
	IsEnabled        bool                `json:"is_enabled"`
} 