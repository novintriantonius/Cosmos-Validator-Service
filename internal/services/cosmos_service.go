package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// RetrieveDelegations fetches delegations for a validator from the Cosmos API
// It includes retry mechanism for transient failures
func (s *CosmosService) RetrieveDelegations(ctx context.Context, validatorAddress string) (*models.DelegationsResponse, error) {
	if validatorAddress == "" {
		return nil, fmt.Errorf("validator address cannot be empty")
	}
	
	url := fmt.Sprintf("%s/cosmos/staking/v1beta1/validators/%s/delegations", s.config.BaseURL, validatorAddress)
	
	var (
		resp *http.Response
		err  error
		reqErr error
	)
	
	// Implement retry mechanism
	for attempt := 0; attempt <= s.config.MaxRetries; attempt++ {
		// Create a new request with the context
		var req *http.Request
		req, reqErr = http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if reqErr != nil {
			return nil, fmt.Errorf("failed to create request: %w", reqErr)
		}
		
		// Set appropriate headers
		req.Header.Set("Accept", "application/json")
		
		// Execute the request
		resp, err = s.client.Do(req)
		
		// If no error or non-retryable error, break the loop
		if err == nil && (resp.StatusCode < 500 || resp.StatusCode == http.StatusNotFound) {
			break
		}
		
		// If this was the last attempt, return the error
		if attempt == s.config.MaxRetries {
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve delegations after %d attempts: %w", s.config.MaxRetries+1, err)
			}
			return nil, fmt.Errorf("failed to retrieve delegations after %d attempts: received status code %d", s.config.MaxRetries+1, resp.StatusCode)
		}
		
		// Close the response body if we received a response
		if resp != nil {
			resp.Body.Close()
		}
		
		// Wait before retrying
		select {
		case <-ctx.Done():
			// Context cancelled or timed out
			return nil, ctx.Err()
		case <-time.After(s.config.RetryDelay * time.Duration(attempt+1)):
			// Exponential backoff
		}
	}
	
	// Handle the response
	if resp == nil {
		return nil, fmt.Errorf("unexpected error: nil response after retries")
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyStr := string(bodyBytes)
		
		return nil, fmt.Errorf("API returned non-200 status code: %d, body: %s", resp.StatusCode, bodyStr)
	}
	
	// Parse the response
	var delegationsResp models.DelegationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&delegationsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	// Add timestamp
	delegationsResp.Timestamp = time.Now()
	
	return &delegationsResp, nil
} 