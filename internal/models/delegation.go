package models

import (
	"time"
)

// DelegationsResponse represents the response from the delegations API
type DelegationsResponse struct {
	DelegationResponses []DelegationResponse `json:"delegation_responses"`
	Pagination          Pagination           `json:"pagination"`
	Timestamp           time.Time            `json:"timestamp"`
}

// DelegationResponse represents a single delegation response
type DelegationResponse struct {
	Delegation Delegation `json:"delegation"`
	Balance    Balance    `json:"balance"`
}

// Delegation represents the delegation details
type Delegation struct {
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