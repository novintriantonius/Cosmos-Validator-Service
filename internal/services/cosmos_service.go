package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/novintriantonius/cosmos-validator-service/internal/models"
)

const (
	// DefaultBaseURL is the default base URL for the Cosmos API
	DefaultBaseURL = "https://cosmos-api.polkachu.com"
	
	// DefaultMaxRetries is the default number of retries for API calls
	DefaultMaxRetries = 3
	
	// DefaultRetryDelay is the default delay between retries in milliseconds
	DefaultRetryDelay = 500
	
	// DefaultTimeout is the default timeout for API calls in seconds
	DefaultTimeout = 10
)

// CosmosServiceConfig holds configuration for the Cosmos service
type CosmosServiceConfig struct {
	BaseURL    string
	MaxRetries int
	RetryDelay time.Duration
	Timeout    time.Duration
	HTTPClient *http.Client
}

// CosmosService provides methods to interact with the Cosmos API
type CosmosService struct {
	config CosmosServiceConfig
	client *http.Client
}

// NewCosmosService creates a new instance of CosmosService with default configurations
func NewCosmosService() *CosmosService {
	return &CosmosService{
		config: CosmosServiceConfig{
			BaseURL:    DefaultBaseURL,
			MaxRetries: DefaultMaxRetries,
			RetryDelay: DefaultRetryDelay * time.Millisecond,
			Timeout:    DefaultTimeout * time.Second,
		},
		client: &http.Client{
			Timeout: DefaultTimeout * time.Second,
		},
	}
}

// NewCosmosServiceWithConfig creates a new instance of CosmosService with custom configurations
func NewCosmosServiceWithConfig(config CosmosServiceConfig) *CosmosService {
	// Apply defaults for empty values
	if config.BaseURL == "" {
		config.BaseURL = DefaultBaseURL
	}
	
	if config.MaxRetries <= 0 {
		config.MaxRetries = DefaultMaxRetries
	}
	
	if config.RetryDelay <= 0 {
		config.RetryDelay = DefaultRetryDelay * time.Millisecond
	}
	
	if config.Timeout <= 0 {
		config.Timeout = DefaultTimeout * time.Second
	}
	
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{
			Timeout: config.Timeout,
		}
	}
	
	return &CosmosService{
		config: config,
		client: config.HTTPClient,
	}
}

// GetConfig returns the service configuration (for testing purposes)
func (s *CosmosService) GetConfig() CosmosServiceConfig {
	return s.config
}

// RetrieveDelegations retrieves delegations for a validator
func (s *CosmosService) RetrieveDelegations(ctx context.Context, validatorAddress string) (*models.DelegationsResponse, error) {
	log.Printf("[DEBUG] Starting RetrieveDelegations for validator %s", validatorAddress)
	
	// Build the URL
	url := fmt.Sprintf("%s/cosmos/staking/v1beta1/validators/%s/delegations", s.config.BaseURL, validatorAddress)
	log.Printf("[DEBUG] Making request to URL: %s", url)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to create request: %v", err)
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Failed to send request: %v", err)
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[ERROR] Unexpected status code: %d, body: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read response body: %v", err)
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	log.Printf("[DEBUG] Received response body: %s", string(body))

	// Parse response
	var delegationsResp models.DelegationsResponse
	if err := json.Unmarshal(body, &delegationsResp); err != nil {
		log.Printf("[ERROR] Failed to parse response: %v", err)
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	log.Printf("[DEBUG] Successfully parsed %d delegations", len(delegationsResp.DelegationResponses))
	return &delegationsResp, nil
} 