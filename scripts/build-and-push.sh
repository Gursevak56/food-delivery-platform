#!/usr/bin/env bash
set -euo pipefail
ACCOUNT_ID=456309724089
REGION=ap-south-1
COMMIT=$(git rev-parse --short HEAD)
services=(go-geolocate-mongo notification-service order-service payment-service restaurant-service rider-service user-service)

# login
aws ecr get-login-password --region ${REGION} | docker login --username AWS --password-stdin ${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com

for svc in "${services[@]}"; do
  echo "Building $svc..."
  docker build -t ${svc}:${COMMIT} ./services/${svc}
  REPO=${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/${svc}
  docker tag ${svc}:${COMMIT} ${REPO}:${COMMIT}
  docker tag ${svc}:${COMMIT} ${REPO}:latest
  echo "Pushing $svc..."
  docker push ${REPO}:${COMMIT}
  docker push ${REPO}:latest
  echo "Pushed ${REPO}:${COMMIT}"
done
echo "All services built and pushed successfully."
