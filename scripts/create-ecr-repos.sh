#!/usr/bin/env bash
set -euo pipefail
REGION=${AWS_REGION:-ap-south-1}
ACCOUNT_ID=${AWS_ACCOUNT_ID:-456309724089}
services=(go-geolocate-mongo notification-service order-service payment-service restaurant-service rider-service user-service)

for svc in "${services[@]}"; do
  echo "Ensuring ECR repo ${svc}..."
  aws ecr describe-repositories --repository-names "${svc}" --region "${REGION}" >/dev/null 2>&1 || \
    aws ecr create-repository --repository-name "${svc}" --region "${REGION}"
done
echo "All ECR repos ensured."
