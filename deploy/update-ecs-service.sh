#!/usr/bin/env bash
set -euo pipefail

if [ "$#" -ne 2 ]; then
  echo "Usage: $0 <service-name> <image-full-ref>"
  exit 2
fi

SERVICE_NAME="$1"   # e.g. user-service
IMAGE="$2"          # e.g. 456309...dkr.ecr.ap-south-1.amazonaws.com/user-service:abcd1234

REGION="${AWS_REGION:-ap-south-1}"
ACCOUNT="${AWS_ACCOUNT_ID:-456309724089}"
CLUSTER="${ECS_CLUSTER:-food-delivery-cluster}"
TASKDEF_TEMPLATE="./deploy/taskdef-template.json"
EXEC_ROLE="${ECS_EXEC_ROLE:-}"
TASK_ROLE="${ECS_TASK_ROLE:-}"

# read template and replace placeholders
TMPFILE=$(mktemp /tmp/taskdef.XXXX.json)
jq --arg family "$SERVICE_NAME" \
   --arg container "$SERVICE_NAME" \
   --arg image "$IMAGE" \
   --arg execRole "${EXEC_ROLE}" \
   --arg taskRole "${TASK_ROLE}" \
   '.family=$family | .containerDefinitions[0].name=$container | .containerDefinitions[0].image=$image |
    (.executionRoleArn = ($execRole | select(. != ""))) |
    (.taskRoleArn = ($taskRole | select(. != "")))' \
   "$TASKDEF_TEMPLATE" > "${TMPFILE}"

echo "Registering task definition for ${SERVICE_NAME}..."
REGISTERED=$(aws ecs register-task-definition --cli-input-json file://"${TMPFILE}" --region "${REGION}")
TASKDEF_ARN=$(echo "${REGISTERED}" | jq -r '.taskDefinition.taskDefinitionArn')

echo "Updating service ${SERVICE_NAME} to use ${TASKDEF_ARN}..."
aws ecs update-service --cluster "${CLUSTER}" --service "${SERVICE_NAME}" --task-definition "${TASKDEF_ARN}" --region "${REGION}" >/dev/null

echo "Updated ${SERVICE_NAME} on cluster ${CLUSTER}."
rm -f "${TMPFILE}"
