package config

import (
	"fmt"
	"os"
	"strconv"

	"deepComparator/pkg/models"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	Database1    models.DatabaseConfig `json:"database1"`
	Database2    models.DatabaseConfig `json:"database2"`
	OutputFormat string                `json:"output_format"`
	OutputFile   string                `json:"output_file"`
	LogLevel     string                `json:"log_level"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig(envFile string) (*Config, error) {
	// Load .env file if provided
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	config := &Config{}

	// Database 1 configuration
	db1Port, err := strconv.Atoi(getEnvOrDefault("DB1_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB1_PORT: %w", err)
	}

	config.Database1 = models.DatabaseConfig{
		Host:     getEnvOrDefault("DB1_HOST", "localhost"),
		Port:     db1Port,
		Database: getEnvOrDefault("DB1_DATABASE", ""),
		Username: getEnvOrDefault("DB1_USERNAME", ""),
		Password: getEnvOrDefault("DB1_PASSWORD", ""),
		SSLMode:  getEnvOrDefault("DB1_SSL_MODE", "disable"),
	}

	// Database 2 configuration
	db2Port, err := strconv.Atoi(getEnvOrDefault("DB2_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB2_PORT: %w", err)
	}

	config.Database2 = models.DatabaseConfig{
		Host:     getEnvOrDefault("DB2_HOST", "localhost"),
		Port:     db2Port,
		Database: getEnvOrDefault("DB2_DATABASE", ""),
		Username: getEnvOrDefault("DB2_USERNAME", ""),
		Password: getEnvOrDefault("DB2_PASSWORD", ""),
		SSLMode:  getEnvOrDefault("DB2_SSL_MODE", "disable"),
	}

	// Application configuration
	config.OutputFormat = getEnvOrDefault("OUTPUT_FORMAT", "json")
	config.OutputFile = getEnvOrDefault("OUTPUT_FILE", "comparison_result.json")
	config.LogLevel = getEnvOrDefault("LOG_LEVEL", "info")

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Database1.Database == "" {
		return fmt.Errorf("DB1_DATABASE is required")
	}
	if c.Database1.Username == "" {
		return fmt.Errorf("DB1_USERNAME is required")
	}
	if c.Database2.Database == "" {
		return fmt.Errorf("DB2_DATABASE is required")
	}
	if c.Database2.Username == "" {
		return fmt.Errorf("DB2_USERNAME is required")
	}
	return nil
}

// getEnvOrDefault returns the environment variable value or the default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
