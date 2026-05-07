#!/usr/bin/env bash
# Antrea Windows Build & Deploy Orchestrator

set -e

MANIFEST_PATH="../bmg-ace-k8s-infra/environments/aws-uk-sbx-audiobroadcast/antrea-windows-with-ovs.yml"
ECR_REPO="346273507914.dkr.ecr.eu-west-2.amazonaws.com/antrea-windows-test"
REGION="eu-west-2"

# 1. Commit changes if any
SOURCE_PATHS="pkg/ cmd/ build/ hack/ Makefile go.mod go.sum"
if [[ -n $(git status -s $SOURCE_PATHS) ]]; then
  echo "===> Committing Antrea source changes..."
  git add $SOURCE_PATHS
  git commit -m "Antrea Windows dev update" || true
fi

# 2. Get 7-char SHA
GIT_SHA=$(git rev-parse --short=7 HEAD)
TAG="pr1-${GIT_SHA}"
FULL_IMAGE="${ECR_REPO}:${TAG}"

echo "===> Target Image: ${FULL_IMAGE}"

# 3. Build and push
echo "===> Building and pushing image..."
export ANTREA_WINDOWS_IMAGE="${ECR_REPO}"
./build/images/build-windows.sh --push --agent-tag "${TAG}"

# 4. Update manifest
echo "===> Updating manifest..."
# Use sed to replace the image tag. We look for the ECR repo pattern and replace the tag.
sed -i "s|image: ${ECR_REPO}:.*|image: ${FULL_IMAGE}|g" "${MANIFEST_PATH}"

# 5. Apply manifest
echo "===> Applying manifest..."
kubectl apply -f "${MANIFEST_PATH}"

echo "===> Done! New image deployed: ${TAG}"
echo "===> Run 'kubectl get pods -n kube-system -l app=antrea,component=antrea-agent -w' to monitor."
