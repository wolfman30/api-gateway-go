package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// SecretsConfig holds secrets loaded from AWS Secrets Manager
type SecretsConfig struct {
	ApiKey            string `json:"api-key"`
	DatabaseURL       string `json:"database-url"`
	JwtSecret         string `json:"jwt-secret"`
	OAuthClientID     string `json:"oauth-client-id"`
	OAuthClientSecret string `json:"oauth-client-secret"`
}

// LoadFromSecretsManager loads secrets from AWS Secrets Manager with environment-specific names
// For local development only, set USE_LOCAL_SECRETS=true to load from environment variables
func LoadFromSecretsManager(ctx context.Context) (*SecretsConfig, error) {
	// Use local secrets only for local development testing
	// Never use environment variables for production/staging/CI-CD
	if IsLocalDevelopment() {
		return loadLocalSecrets()
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	client := secretsmanager.NewFromConfig(cfg)
	secretsConfig := &SecretsConfig{}

	// Map of base secret names to their config fields
	baseSecrets := map[string]*string{
		"api-key":             &secretsConfig.ApiKey,
		"database-url":        &secretsConfig.DatabaseURL,
		"jwt-secret":          &secretsConfig.JwtSecret,
		"oauth-client-id":     &secretsConfig.OAuthClientID,
		"oauth-client-secret": &secretsConfig.OAuthClientSecret,
	}

	currentEnv := GetCurrentEnvironment()

	for baseName, field := range baseSecrets {
		secretName := GetSecretName(baseName)

		// Try environment-specific secret first, then fallback to base name
		result, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secretName),
		})

		if err != nil {
			// If environment-specific secret not found, try base name
			log.Printf("Secret '%s' not found, trying fallback '%s'", secretName, baseName)
			result, err = client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
				SecretId: aws.String(baseName),
			})

			if err != nil {
				log.Printf("Warning: Secret '%s' (fallback) also not found: %v", baseName, err)
				continue
			}
		}

		secretString := ""
		if result.SecretString != nil {
			secretString = *result.SecretString
		}

		// Try to parse as JSON first
		if secretString != "" {
			var secretData map[string]interface{}
			if err := json.Unmarshal([]byte(secretString), &secretData); err == nil {
				// It's JSON, extract the value
				if val, ok := secretData[baseName]; ok {
					*field = fmt.Sprintf("%v", val)
				} else {
					*field = secretString
				}
			} else {
				// Not JSON, use as-is
				*field = secretString
			}
		}

		log.Printf("Loaded secret '%s' for environment '%s'", secretName, currentEnv)
	}

	return secretsConfig, nil
}

// loadLocalSecrets loads secrets from environment variables (LOCAL DEVELOPMENT ONLY)
// WARNING: This should NEVER be used in production, staging, or CI/CD environments.
// It is only for local testing when USE_LOCAL_SECRETS=true is explicitly set.
func loadLocalSecrets() (*SecretsConfig, error) {
	log.Println("WARNING: Loading secrets from environment variables (LOCAL DEVELOPMENT ONLY)")
	return &SecretsConfig{
		ApiKey:            os.Getenv("LOCAL_API_KEY"),
		DatabaseURL:       os.Getenv("LOCAL_DATABASE_URL"),
		JwtSecret:         os.Getenv("LOCAL_JWT_SECRET"),
		OAuthClientID:     os.Getenv("LOCAL_OAUTH_CLIENT_ID"),
		OAuthClientSecret: os.Getenv("LOCAL_OAUTH_CLIENT_SECRET"),
	}, nil
}
