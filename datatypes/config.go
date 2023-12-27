package datatypes

import (
	"encoding/json"
	"os"
)

type RateLimitConfig struct {
	RateLimitsPerEndpoint []RateLimitsPerEndpoint `json:"rateLimitsPerEndpoint,omitempty"`
}

type RateLimitsPerEndpoint struct {
	Endpoint  string `json:"endpoint,omitempty"`
	Burst     int    `json:"burst,omitempty"`
	Sustained int    `json:"sustained,omitempty"`
}

// LoadConfig reads configuration from given input JSON file
func LoadConfig(path string) (*RateLimitConfig, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg RateLimitConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
