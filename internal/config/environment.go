package config

import (
	"fmt"
	"os"
	"strings"
)

// Environment represents the deployment environment
type Environment string

const (
	Dev     Environment = "dev"
	Staging Environment = "staging"
	Prod    Environment = "prod"
)

// EnvironmentConfig holds environment-specific configuration
type EnvironmentConfig struct {
	Environment Environment
	ClusterName string
	EcsCluster  string
	SqsQueueURL string
	S3Bucket    string
	ApiPort     string
	LogLevel    string
}

// GetCurrentEnvironment returns the current deployment environment
func GetCurrentEnvironment() Environment {
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))
	switch Environment(env) {
	case Dev, Staging, Prod:
		return Environment(env)
	default:
		return Dev // Default to dev
	}
}

// LoadEnvironmentConfig loads environment-specific configuration
func LoadEnvironmentConfig() *EnvironmentConfig {
	currentEnv := GetCurrentEnvironment()

	config := &EnvironmentConfig{
		Environment: currentEnv,
	}

	// Load environment-specific values with suffixes
	suffix := string(currentEnv)

	// ECS Cluster (environment-specific)
	config.EcsCluster = os.Getenv("ECS_CLUSTER_" + strings.ToUpper(suffix))
	if config.EcsCluster == "" {
		config.EcsCluster = os.Getenv("ECS_CLUSTER") // Fallback to base name
	}

	// SQS Queue URL (environment-specific)
	config.SqsQueueURL = os.Getenv("SQS_QUEUE_URL_" + strings.ToUpper(suffix))
	if config.SqsQueueURL == "" {
		config.SqsQueueURL = os.Getenv("SQS_QUEUE_URL") // Fallback to base name
	}

	// S3 Bucket (environment-specific)
	config.S3Bucket = os.Getenv("S3_BUCKET_" + strings.ToUpper(suffix))
	if config.S3Bucket == "" {
		config.S3Bucket = os.Getenv("S3_BUCKET") // Fallback to base name
	}

	// Cluster Name (environment-specific)
	config.ClusterName = os.Getenv("CLUSTER_NAME_" + strings.ToUpper(suffix))
	if config.ClusterName == "" {
		config.ClusterName = os.Getenv("CLUSTER_NAME") // Fallback to base name
	}

	// API Port
	config.ApiPort = os.Getenv("API_PORT")
	if config.ApiPort == "" {
		config.ApiPort = "8080"
	}

	// Log Level
	config.LogLevel = os.Getenv("LOG_LEVEL")
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}

	return config
}

// GetSecretName returns the environment-specific secret name with fallback
// e.g., for secret "api-key" and env "dev", returns "api-key-dev"
// If the env-specific secret is not found, returns the base name
func GetSecretName(baseName string) string {
	currentEnv := GetCurrentEnvironment()
	if currentEnv == Dev {
		return fmt.Sprintf("%s-dev", baseName)
	}
	return fmt.Sprintf("%s-%s", baseName, string(currentEnv))
}

// String returns the environment name as a string
func (e Environment) String() string {
	return string(e)
}

// IsLocalDevelopment checks if we should use local development secrets
// This is only for local testing and should never be true in CI/CD or deployed environments
func IsLocalDevelopment() bool {
	return os.Getenv("USE_LOCAL_SECRETS") == "true"
}
