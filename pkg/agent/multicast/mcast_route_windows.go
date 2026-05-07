//go:build windows
// +build windows

// Copyright 2026 Antrea Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package multicast

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"

	"antrea.io/antrea/v2/pkg/agent/util"
)

func (c *MRouteClient) Initialize() error {
	// On Windows, the network transformation (HNS/OVS) can take some time to settle.
	// We wait for the multicast interfaces (especially antrea-gw0) to be ready.
	if len(c.multicastInterfaces) > 0 {
		klog.InfoS("Waiting for multicast interfaces to settle", "interfaces", c.multicastInterfaces)
		err := wait.PollUntilContextTimeout(context.TODO(), 2*time.Second, 20*time.Second, true, func(ctx context.Context) (bool, error) {
			for _, ifaceName := range c.multicastInterfaces {
				ipv4Addr, _, _, err := util.GetIPNetDeviceByName(ifaceName)
				if err != nil || ipv4Addr == nil {
					klog.V(2).InfoS("Multicast interface not ready yet", "interface", ifaceName, "err", err)
					return false, nil
				}
			}
			return true, nil
		})
		if err != nil {
			klog.ErrorS(err, "Multicast interfaces did not settle in time", "interfaces", c.multicastInterfaces)
			// On Windows, we proceed even if interfaces aren't ready to avoid blocking agent startup.
			// The interfaces might be undergoing recreation by boot scripts.
		} else {
			klog.InfoS("Multicast interfaces settled")
		}
	}

	c.setMulticastInterfaces()
	// On Windows, multicast routing is handled by OpenFlow, so we don't need to
	// allocate VIFs for the gateway interface or external interfaces in the host
	// kernel. We still initialize these fields to safe values to avoid any
	// potential issues with other shared code.
	c.internalInterfaceVIF = 0
	c.externalInterfaceVIFs = make([]uint16, len(c.multicastInterfaceConfigs))
	for i := range c.multicastInterfaceConfigs {
		c.externalInterfaceVIFs[i] = uint16(i + 1)
	}
	return nil
}

func (c *MRouteClient) parseIGMPMsg(msg []byte) (*parsedIGMPMsg, error) {
	// Multicast routing messages from raw sockets are not supported on Windows.
	return nil, nil
}

func (c *MRouteClient) run(stopCh <-chan struct{}) {
	// Multicast routing worker is not required on Windows.
}
