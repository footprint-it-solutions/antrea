---
name: antrea-windows-dev
description: Automates iterative development, building, testing, and deployment of Antrea on Windows. Includes automatic recovery of faulty EC2 instances via ASG detachment and termination.
---

# Antrea Windows Dev

## Overview

This skill enables an automated "inner loop" for developing and testing Antrea features on Windows nodes. It handles the complexity of cross-compiling Windows images, updating Kubernetes manifests, and managing faulty cloud infrastructure.

## Inner Loop Workflow

1.  **Code Changes**: Perform surgical edits to Antrea source code (e.g., multicast logic in `pkg/agent/openflow/`).
2.  **Build & Deploy**: Use `scripts/build_push_update.sh` to:
    *   Commit changes with "Antrea Windows dev update".
    *   Build and push the image to ECR with a `pr1-<SHA>` tag.
    *   Update the `antrea-agent-windows` DaemonSet manifest.
    *   Apply the manifest to the cluster.
3.  **Observation**: Monitor the new agent pod and node status.
    *   Check logs for success/failure patterns (see `references/patterns.md`).
    *   Verify RDP and SSM connectivity.

### Example Request
"I've finished the multicast fixes. Build a new image and deploy it."

## Fault Recovery Workflow

If a Windows node becomes "Unrecoverable" (e.g., OVS driver crash leading to BSOD/SSM timeout), follow this procedure:

**CRITICAL NOTE**: New Windows nodes take approximately **15 minutes** to complete the multiple reboots and network configuration required by the boot scripts (`PrepareAntrea.ps1`). The node will **not** appear in the `kubectl get nodes` list until the final reboot is finished and the `kubelet` service starts.

1.  **Identify Instance**: Get the EC2 Instance ID for the `NotReady` node.
2.  **Recover**: Execute `scripts/recover_instance.sh <instance-id>` to:
    *   Dynamically discover the ASG (tagged with `eks:cluster-name=aws-uk-sbx-audiobroadcast`).
    *   Detach the instance from the ASG (without decrementing capacity).
    *   Terminate the instance.
3.  **Verification**: Wait for the ASG to launch a replacement node and for the Antrea agent to initialize.

### Example Request
"Node ip-10-231-229-80 is unreachable and SSM is timing out. Replace it."

## Resources

### scripts/
- `build_push_update.sh`: Orchestrates Git, Docker, and Kubectl for deployment.
- `recover_instance.sh`: Automates ASG detachment and EC2 termination.

### references/
- `patterns.md`: Catalog of success and failure log patterns for Antrea on Windows.
