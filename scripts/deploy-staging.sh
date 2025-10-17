#!/bin/bash
set -e

# Deploy API Gateway to Staging Environment
# Usage: ./deploy-staging.sh

ENVIRONMENT="staging"
AWS_REGION="${AWS_REGION:-us-east-1}"
CLUSTER_NAME="${CLUSTER_NAME:-api-gateway-staging}"
SERVICE_NAME="${SERVICE_NAME:-api-gateway}"

echo "=========================================="
echo "Deploying API Gateway to Staging"
echo "=========================================="
echo "Environment: $ENVIRONMENT"
echo "Region: $AWS_REGION"
echo "Cluster: $CLUSTER_NAME"
echo "Service: $SERVICE_NAME"
echo ""

# Pre-deployment checks
echo "Running pre-deployment checks..."
if ! command -v curl &> /dev/null; then
    echo "⚠️  Warning: curl not found, skipping pre-deployment smoke test"
else
    echo "✓ curl is available"
fi

# Get current task definition
TASK_DEF=$(aws ecs describe-services \
  --cluster "$CLUSTER_NAME" \
  --services "$SERVICE_NAME" \
  --region "$AWS_REGION" \
  --query 'services[0].taskDefinition' \
  --output text)

echo "Current task definition: $TASK_DEF"

# Update service with force new deployment
echo "Updating service..."
aws ecs update-service \
  --cluster "$CLUSTER_NAME" \
  --service "$SERVICE_NAME" \
  --task-definition "$TASK_DEF" \
  --force-new-deployment \
  --region "$AWS_REGION"

# Wait for service to stabilize
echo "Waiting for service to stabilize..."
aws ecs wait services-stable \
  --cluster "$CLUSTER_NAME" \
  --services "$SERVICE_NAME" \
  --region "$AWS_REGION"

# Get service status
echo ""
echo "Service status after deployment:"
aws ecs describe-services \
  --cluster "$CLUSTER_NAME" \
  --services "$SERVICE_NAME" \
  --region "$AWS_REGION" \
  --query 'services[0].[serviceName, status, desiredCount, runningCount]' \
  --output table

# Post-deployment smoke tests
echo ""
echo "Running post-deployment smoke tests..."
# Add your smoke tests here
# Example: curl http://staging-api-gateway.internal:8080/health

echo ""
echo "✅ API Gateway successfully deployed to Staging!"
echo ""
echo "Next steps:"
echo "  - Check logs: aws logs tail /ecs/api-gateway-staging --follow"
echo "  - Run integration tests: npm run test:integration"
echo "  - Request staging approval in GitHub"
