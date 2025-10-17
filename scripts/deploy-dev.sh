#!/bin/bash
set -e

# Deploy API Gateway to Dev Environment
# Usage: ./deploy-dev.sh

ENVIRONMENT="dev"
AWS_REGION="${AWS_REGION:-us-east-1}"
CLUSTER_NAME="${CLUSTER_NAME:-api-gateway-dev}"
SERVICE_NAME="${SERVICE_NAME:-api-gateway}"

echo "=========================================="
echo "Deploying API Gateway to Dev"
echo "=========================================="
echo "Environment: $ENVIRONMENT"
echo "Region: $AWS_REGION"
echo "Cluster: $CLUSTER_NAME"
echo "Service: $SERVICE_NAME"
echo ""

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
  --query 'services[0].[serviceName, status, desiredCount, runningCount, deployments[0].[taskDefinition, desiredCount, runningCount]]' \
  --output table | jq .

echo ""
echo "✅ API Gateway successfully deployed to Dev!"
echo ""
echo "Next steps:"
echo "  - Check logs: aws logs tail /ecs/api-gateway-dev --follow"
echo "  - Health check: curl http://dev-api-gateway.internal:8080/health"
echo "  - Rollback: aws ecs update-service --cluster api-gateway-dev --service api-gateway --task-definition <previous-revision>"
