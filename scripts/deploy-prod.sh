#!/bin/bash
set -e

# Deploy API Gateway to Production Environment
# Usage: ./deploy-prod.sh

ENVIRONMENT="prod"
AWS_REGION="${AWS_REGION:-us-east-1}"
CLUSTER_NAME="${CLUSTER_NAME:-api-gateway-prod}"
SERVICE_NAME="${SERVICE_NAME:-api-gateway}"

read -p "Are you sure you want to deploy to PRODUCTION? (yes/no): " CONFIRM
if [ "$CONFIRM" != "yes" ]; then
  echo "Aborting deployment"
  exit 1
fi

echo "=========================================="
echo "Deploying API Gateway to Production"
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

# Blue-green / rolling approach: update service with force new deployment
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

echo ""
echo "✅ API Gateway successfully deployed to Production!"
echo ""
echo "Next steps:"
echo "  - Verify health: curl https://api.example.com/health"
echo "  - Monitor logs: aws logs tail /ecs/api-gateway-prod --follow"
echo "  - If needed, rollback using previous task definition"
