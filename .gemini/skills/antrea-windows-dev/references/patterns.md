# Antrea Windows Diagnostic Patterns

## Success Patterns
- `Multicast interfaces settled`: Indicates 'antrea-gw0' is ready.
- `Agent initialized NodeConfig`: Node configuration is complete.
- `Installed Pod network forwarding rules`: CNI was able to set up pod networking.
- `OFSwitch is connected`: OpenFlow connection to OVS is established.

## Failure Patterns
- `timeout waiting for group`: Usually indicates `antrea-agent` is stuck waiting for initialization or networking is broken.
- `ovs_client.go:90] Not connected yet`: OVS services are likely not running or the bridge hasn't been created.
- `BUGCHECK CODE 0000003B`: Kernel crash (BSOD) in `ovsext.sys`.
- `dial tcp <IP>:10250: i/o timeout`: Kubelet is unreachable, often due to networking instability or crash.

## Recovery Triggers
- SSM command returns `Timeout`.
- EC2 status check `reachability` is `failed`.
- Pod status `Terminating` for more than 5 minutes.
