API Gateway - CI/CD

This repository contains GitHub Actions workflows and deployment scripts to support a 3-environment pipeline: dev -> staging -> prod.

Files added:
- .github/workflows/ci-cd.yml: Tests, build, and deployment pipeline (dev/staging/prod)
- internal/config/environment.go: Environment helpers and environment config loader
- internal/config/secrets.go: AWS Secrets Manager loader (with local fallback)
- scripts/deploy-dev.sh: Dev deployment script
- scripts/deploy-staging.sh: Staging deployment script
- scripts/deploy-prod.sh: Production deployment script (manual confirmation)

How it works:
- Push to develop: runs tests and builds; if push to develop triggers deploy-dev job
- Push to staging: runs build + deploy to staging (runs staging tests)
- Push to main: runs build + deploy to production (requires confirmation in script)

Secrets and prerequisites:
- AWS credentials must be available via GitHub OIDC or secrets
- Required env vars: ECS_CLUSTER, SQS_QUEUE_URL, S3_BUCKET, CLUSTER_NAME
- Optional: USE_LOCAL_SECRETS=true for local development

Notes:
- Update cluster and service names in the scripts to match your AWS account
- For production, prefer manual approval on GitHub environment and OIDC-based deployment
