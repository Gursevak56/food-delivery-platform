#!/usr/bin/env bash
set -euo pipefail

# === Edit these if needed ===
REGION="ap-south-1"
ACCOUNT_ID="456309724089"
CLUSTER="food-delivery-cluster"
# from your message:
SECURITY_GROUPS="sg-021ff3ef2d071b320"
# list your subnets (no spaces)
SUBNETS="subnet-007728462ef80ac33,subnet-0b43b446e3f82715d,subnet-0129d8e8384225791"
# desired count per service
DESIRED_COUNT=2
# container port used by all services
CONTAINER_PORT=8080

services=(go-geolocate-mongo notification-service order-service payment-service restaurant-service rider-service user-service)

# helper to fetch target group ARN by name (service-tg)
get_tg_arn() {
  svc="$1"
  tg_name="${svc}-tg"
  arn=$(aws elbv2 describe-target-groups \
    --names "${tg_name}" \
    --region "${REGION}" \
    --query 'TargetGroups[0].TargetGroupArn' \
    --output text 2>/dev/null || echo "")
  echo "${arn}"
}

# loop and create services
for svc in "${services[@]}"; do
  echo "=== Creating ECS service for: ${svc} ==="
  TG_ARN=$(get_tg_arn "${svc}")

  if [ -z "${TG_ARN}" ] || [ "${TG_ARN}" == "None" ]; then
    echo "ERROR: Target group '${svc}-tg' not found in region ${REGION}."
    echo "Create it with (example):"
    echo "  aws elbv2 create-target-group \\"
    echo "    --name ${svc}-tg \\"
    echo "    --protocol HTTP \\"
    echo "    --port ${CONTAINER_PORT} \\"
    echo "    --target-type ip \\"
    echo "    --vpc-id <your-vpc-id> \\"
    echo "    --health-check-path /health \\"
    echo "    --region ${REGION}"
    echo "Skipping ${svc} — please create the target group and re-run."
    echo
    continue
  fi

  # create the ECS service (one-time). If it already exists this will fail.
  echo "Using targetGroupArn: ${TG_ARN}"
  echo "Creating service ${svc} in cluster ${CLUSTER}..."

  aws ecs create-service \
    --cluster "${CLUSTER}" \
    --service-name "${svc}" \
    --task-definition "${svc}" \
    --desired-count ${DESIRED_COUNT} \
    --launch-type FARGATE \
    --network-configuration "awsvpcConfiguration={subnets=[${SUBNETS}],securityGroups=[${SECURITY_GROUPS}],assignPublicIp=ENABLED}" \
    --load-balancers "targetGroupArn=${TG_ARN},containerName=${svc},containerPort=${CONTAINER_PORT}" \
    --region "${REGION}"

  echo "Created service ${svc}."
  echo
done

echo "All done. If any services were skipped, create their target groups then re-run."
