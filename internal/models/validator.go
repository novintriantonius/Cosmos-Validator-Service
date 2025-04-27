package models

// Validator represents a cosmos validator entity
type Validator struct {
	Name            string `json:"name"`
	Address         string `json:"address"`
	EnabledTracking bool   `json:"enabledTracking"`
} 