# Secrets and Configuration Management

This document explains how the API Gateway handles secrets and configuration across different environments.

## Overview

The API Gateway uses a **two-tier configuration system**:

1. **Secrets** - Sensitive values (API keys, database URLs, JWT secrets) loaded from **AWS Secrets Manager**
2. **Environment Configuration** - Infrastructure settings (queue URLs, S3 buckets, cluster names) loaded from **environment variables**

### Key Principle

- **Secrets**: Always use AWS Secrets Manager (except for local development)
- **Config**: Environment variables (DNS-resolvable infrastructure URLs)
- **Local Dev**: Use `USE_LOCAL_SECRETS=true` flag to temporarily load from env vars

---

## Secrets Management

### Production, Staging, CI/CD Environments

Secrets are **automatically loaded from AWS Secrets Manager** on startup:

```go
// main.go
secrets, err := config.LoadFromSecretsManager(ctx)
if err != nil {
    log.Fatalf("Failed to load secrets: %v", err)
}
```

**Secret Names** follow a pattern based on environment:

| Environment | Secret Name Pattern |
|---|---|
| dev | `{baseName}-dev` (e.g., `api-key-dev`) |
| staging | `{baseName}-staging` |
| prod | `{baseName}-prod` |

**Fallback Behavior**: If an environment-specific secret is not found, the loader attempts the base name.

**Supported Secrets:**
- `api-key` → `SecretsConfig.ApiKey`
- `database-url` → `SecretsConfig.DatabaseURL`
- `jwt-secret` → `SecretsConfig.JwtSecret`
- `oauth-client-id` → `SecretsConfig.OAuthClientID`
- `oauth-client-secret` → `SecretsConfig.OAuthClientSecret`

**AWS Region**: Automatically determined by AWS SDK (uses `~/.aws/credentials` or IAM role)

### Local Development Only

For local testing, you can load secrets from environment variables:

```bash
# Enable local development mode
export USE_LOCAL_SECRETS=true

# Set local secret environment variables
export LOCAL_API_KEY="your-api-key"
export LOCAL_DATABASE_URL="postgresql://localhost:5432/db"
export LOCAL_JWT_SECRET="dev-secret"
export LOCAL_OAUTH_CLIENT_ID="dev-client"
export LOCAL_OAUTH_CLIENT_SECRET="dev-secret"

# Run the server
go run ./cmd/server
```

**⚠️ WARNING**: This mode is **ONLY for local development** and should never be enabled in:
- CI/CD pipelines
- Deployed environments
- Production systems

The code will log a warning when using local secrets:
```
WARNING: Loading secrets from environment variables (LOCAL DEVELOPMENT ONLY)
```

---

## Environment Configuration

Environment-specific infrastructure settings are loaded from environment variables:

```go
// main.go
envConfig := config.LoadEnvironmentConfig()
log.Printf("Running in environment: %s", envConfig.Environment)
```

### Detected Environment

The current environment is determined by the `ENVIRONMENT` variable:

```bash
export ENVIRONMENT=dev      # Development
export ENVIRONMENT=staging  # Staging
export ENVIRONMENT=prod     # Production
```

Defaults to `dev` if not specified.

### Configuration Variables

#### SQS Queue URL
```bash
# Environment-specific
export SQS_QUEUE_URL_DEV=https://sqs.us-east-1.amazonaws.com/123456789012/dev-queue
export SQS_QUEUE_URL_STAGING=https://sqs.us-east-1.amazonaws.com/123456789012/staging-queue
export SQS_QUEUE_URL_PROD=https://sqs.us-east-1.amazonaws.com/123456789012/prod-queue

# Or base name (used if environment-specific not found)
export SQS_QUEUE_URL=https://sqs.us-east-1.amazonaws.com/123456789012/queue
```

#### S3 Bucket
```bash
export S3_BUCKET_DEV=my-app-dev
export S3_BUCKET_STAGING=my-app-staging
export S3_BUCKET_PROD=my-app-prod
```

#### ECS Cluster
```bash
export ECS_CLUSTER_DEV=dev-cluster
export ECS_CLUSTER_STAGING=staging-cluster
export ECS_CLUSTER_PROD=prod-cluster
```

#### Cluster Name
```bash
export CLUSTER_NAME_DEV=dev-k8s-cluster
export CLUSTER_NAME_STAGING=staging-k8s-cluster
export CLUSTER_NAME_PROD=prod-k8s-cluster
```

#### API Port (Same for all environments)
```bash
export API_PORT=8080  # Defaults to 8080 if not set
```

#### Log Level (Same for all environments)
```bash
export LOG_LEVEL=info  # Defaults to info if not set
```

---

## Startup Flow

When the API Gateway starts, it:

1. **Load Secrets** → From AWS Secrets Manager (or local env vars if `USE_LOCAL_SECRETS=true`)
2. **Load Environment Config** → From environment variables
3. **Initialize AWS Clients** → SQS, S3, etc. (using IAM role or credentials)
4. **Start Server** → Listen on configured port

```
main.go
├── config.LoadFromSecretsManager(ctx)     ← AWS Secrets Manager
├── config.LoadEnvironmentConfig()         ← Environment Variables
├── awsconfig.LoadDefaultConfig(ctx)       ← AWS SDK
├── sqs.NewFromConfig(awsCfg)              ← Create SQS Client
└── http.ListenAndServe(addr, mux)         ← Start Server
```

---

## Setup for Each Environment

### Development (Local Machine)

```bash
# Enable local development secrets
export USE_LOCAL_SECRETS=true
export ENVIRONMENT=dev

# Set local secrets
export LOCAL_API_KEY="test-key"
export LOCAL_DATABASE_URL="postgresql://localhost:5432/db"
export LOCAL_JWT_SECRET="test-secret"
export LOCAL_OAUTH_CLIENT_ID="test-client"
export LOCAL_OAUTH_CLIENT_SECRET="test-secret"

# Set environment config
export SQS_QUEUE_URL=http://localhost:9324/000000000000/local-queue
export S3_BUCKET=local-bucket
export API_PORT=8080

# Run
go run ./cmd/server
```

### Staging (AWS)

Set up in GitHub Actions secrets or AWS Systems Manager Parameter Store:

```bash
ENVIRONMENT=staging
# Secrets are loaded from AWS Secrets Manager (api-key-staging, database-url-staging, etc.)
# Config comes from environment variables in deployment script
```

### Production (AWS)

```bash
ENVIRONMENT=prod
# Secrets are loaded from AWS Secrets Manager (api-key-prod, database-url-prod, etc.)
# Config comes from environment variables in deployment script
```

---

## AWS Secrets Manager Setup

To set up secrets in AWS Secrets Manager:

```bash
# Use AWS CLI to create secrets
aws secretsmanager create-secret \
  --name api-key-dev \
  --secret-string "your-api-key" \
  --region us-east-1

aws secretsmanager create-secret \
  --name database-url-dev \
  --secret-string "postgresql://user:password@host:5432/db" \
  --region us-east-1

# Or store as JSON
aws secretsmanager create-secret \
  --name api-gateway-secrets-dev \
  --secret-string '{
    "api-key": "your-key",
    "database-url": "postgresql://...",
    "jwt-secret": "secret",
    "oauth-client-id": "id",
    "oauth-client-secret": "secret"
  }' \
  --region us-east-1
```

The loader automatically parses JSON secrets and extracts individual values.

---

## Code Examples

### Loading Secrets and Config in Your Code

```go
package main

import (
    "context"
    "log"
    "github.com/wolfman30/api-gateway-go/internal/config"
)

func main() {
    ctx := context.Background()
    
    // Load secrets
    secrets, err := config.LoadFromSecretsManager(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    // Load config
    envConfig := config.LoadEnvironmentConfig()
    
    // Use them
    log.Printf("Connected to queue: %s", envConfig.SqsQueueURL)
    log.Printf("Using bucket: %s", envConfig.S3Bucket)
    
    // secrets.ApiKey, secrets.DatabaseURL, etc. are ready to use
}
```

### Checking if in Local Development Mode

```go
if config.IsLocalDevelopment() {
    log.Println("Running in local development mode")
} else {
    log.Println("Running in production mode")
}
```

---

## Troubleshooting

### Secrets Not Found

Check AWS Secrets Manager for the expected secret name:

```bash
aws secretsmanager list-secrets --region us-east-1
aws secretsmanager get-secret-value --secret-id api-key-dev --region us-east-1
```

The loader logs warnings for missing secrets:
```
Warning: Secret 'api-key-dev' (fallback) also not found: ...
```

### Wrong Environment Detected

Verify the `ENVIRONMENT` variable:

```bash
echo $ENVIRONMENT
# Should output: dev, staging, or prod
```

### Local Secrets Not Loading

Ensure `USE_LOCAL_SECRETS=true` is set:

```bash
export USE_LOCAL_SECRETS=true
go run ./cmd/server
# Should log: WARNING: Loading secrets from environment variables (LOCAL DEVELOPMENT ONLY)
```

### AWS Credentials Not Found

Ensure AWS credentials are available:

```bash
# Check credentials file
cat ~/.aws/credentials

# Or use IAM role (recommended for deployments)
# Or set via environment variables
export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...
```

---

## Security Best Practices

1. ✅ **Never commit secrets** to version control
2. ✅ **Use AWS Secrets Manager** for production/staging
3. ✅ **Rotate secrets regularly** (AWS Secrets Manager supports automatic rotation)
4. ✅ **Use IAM roles** instead of static credentials in deployments
5. ✅ **Enable encryption** for secrets at rest (KMS)
6. ✅ **Audit secret access** using CloudTrail
7. ⚠️ **Local development only** for `USE_LOCAL_SECRETS=true`
8. ✅ **Use environment variables** for non-sensitive infrastructure config

---

## Related Files

- `internal/config/secrets.go` - Secrets loading logic
- `internal/config/environment.go` - Environment configuration logic
- `cmd/server/main.go` - Application startup
