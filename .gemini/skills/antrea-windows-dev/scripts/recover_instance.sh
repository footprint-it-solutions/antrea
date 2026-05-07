#!/usr/bin/env bash
# Antrea Windows Faulty Instance Recoverer

set -e

INSTANCE_ID=$1
REGION="eu-west-2"

if [[ -z "${INSTANCE_ID}" ]]; then
  echo "Usage: $0 <instance-id>"
  exit 1
fi

echo "===> Recovering instance ${INSTANCE_ID}..."

# 1. Discover ASG
echo "===> Discovering ASG..."
ASG_NAME=$(aws ec2 describe-instances --instance-ids "${INSTANCE_ID}" --region "${REGION}" --query 'Reservations[].Instances[].Tags[?Key==`aws:autoscaling:groupName`].Value' --output text)

if [[ -z "${ASG_NAME}" || "${ASG_NAME}" == "None" ]]; then
  echo "!!! Could not find ASG name via standard tag. Attempting tag-based discovery..."
  # Fallback: Find ASG that has the specific EKS tags
  ASG_NAME=$(aws autoscaling describe-auto-scaling-groups --region "${REGION}" --query 'AutoScalingGroups[?Tags[?Key==`eks:cluster-name` && Value==`aws-uk-sbx-audiobroadcast`] && Tags[?Key==`eks:nodegroup-name` && starts_with(Value, `windows-server`)]].AutoScalingGroupName' --output text | head -n 1)
fi

if [[ -z "${ASG_NAME}" ]]; then
  echo "!!! Could not identify ASG for instance ${INSTANCE_ID}. Manual intervention required."
  exit 1
fi

echo "===> ASG identified: ${ASG_NAME}"

# 2. Detach instance (this triggers ASG to launch a replacement if desired capacity is maintained)
echo "===> Detaching instance from ASG (without decrementing capacity)..."
aws autoscaling detach-instances --instance-ids "${INSTANCE_ID}" --auto-scaling-group-name "${ASG_NAME}" --no-should-decrement-desired-capacity --region "${REGION}"

# 3. Terminate instance
echo "===> Terminating faulty instance..."
aws ec2 terminate-instances --instance-ids "${INSTANCE_ID}" --region "${REGION}"

echo "===> Recovery initiated. A new instance should be launched by ASG: ${ASG_NAME}"
