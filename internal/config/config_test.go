package config

import (
	"os"
	"testing"
)

func TestGetCurrentEnvironment_Default(t *testing.T) {
	// Ensure ENVIRONMENT is unset or invalid
	os.Unsetenv("ENVIRONMENT")
	env := GetCurrentEnvironment()
	if env != Dev {
		t.Errorf("Expected default environment Dev, got %s", env)
	}
}

func TestGetCurrentEnvironment_StagingProd(t *testing.T) {
	os.Setenv("ENVIRONMENT", "Staging")
	if GetCurrentEnvironment() != Staging {
		t.Errorf("Expected Staging, got %s", GetCurrentEnvironment())
	}

	os.Setenv("ENVIRONMENT", "PROD")
	if GetCurrentEnvironment() != Prod {
		t.Errorf("Expected Prod, got %s", GetCurrentEnvironment())
	}
}

func TestLoadEnvironmentConfig_Fallbacks(t *testing.T) {
	// Test environment-specific and fallback values
	os.Setenv("ENVIRONMENT", "dev")
	os.Setenv("SQS_QUEUE_URL_DEV", "url-dev")
	os.Setenv("S3_BUCKET_DEV", "bucket-dev")
	os.Setenv("API_PORT", "9090")
	os.Setenv("LOG_LEVEL", "debug")

	cfg := LoadEnvironmentConfig()
	if cfg.SqsQueueURL != "url-dev" {
		t.Errorf("Expected SqsQueueURL 'url-dev', got %s", cfg.SqsQueueURL)
	}
	if cfg.S3Bucket != "bucket-dev" {
		t.Errorf("Expected S3Bucket 'bucket-dev', got %s", cfg.S3Bucket)
	}
	if cfg.ApiPort != "9090" {
		t.Errorf("Expected ApiPort '9090', got %s", cfg.ApiPort)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("Expected LogLevel 'debug', got %s", cfg.LogLevel)
	}
}

func TestGetSecretName(t *testing.T) {
	os.Setenv("ENVIRONMENT", "dev")
	if name := GetSecretName("api-key"); name != "api-key-dev" {
		t.Errorf("Expected api-key-dev for dev, got %s", name)
	}

	os.Setenv("ENVIRONMENT", "prod")
	if name := GetSecretName("api-key"); name != "api-key-prod" {
		t.Errorf("Expected api-key-prod for prod, got %s", name)
	}
}

func TestIsLocalDevelopment(t *testing.T) {
	os.Unsetenv("USE_LOCAL_SECRETS")
	if IsLocalDevelopment() {
		t.Error("Expected false when USE_LOCAL_SECRETS is unset")
	}

	os.Setenv("USE_LOCAL_SECRETS", "true")
	if !IsLocalDevelopment() {
		t.Error("Expected true when USE_LOCAL_SECRETS is 'true'")
	}
}
