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

func (c *MRouteClient) parseIGMPMsg(msg []byte) (*parsedIGMPMsg, error) {
	// Multicast routing messages from raw sockets are not supported on Windows.
	return nil, nil
}

func (c *MRouteClient) run(stopCh <-chan struct{}) {
	// Multicast routing worker is not required on Windows.
}
