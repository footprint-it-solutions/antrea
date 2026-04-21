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
	"fmt"
	"net"
	"syscall"

	"k8s.io/klog/v2"
)

const (
	IGMPMsgNocache = 0
	MaxVIFs        = 32
	SizeofIgmpmsg  = 0
)

type Socket struct {
	sockFD syscall.Handle
}

func (s *Socket) AddMrouteEntry(src net.IP, group net.IP, iif uint16, oifs []uint16) error {
	return nil
}

func (s *Socket) GetMroutePacketCount(src net.IP, group net.IP) (uint32, error) {
	return 0, nil
}

func (s *Socket) DelMrouteEntry(src net.IP, group net.IP, iif uint16) error {
	return nil
}

func (s *Socket) FlushMRoute() {
}

func CreateMulticastSocket() (*Socket, error) {
	// We use a UDP socket to join multicast groups on the host.
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		return nil, fmt.Errorf("failed to create multicast socket: %v", err)
	}
	return &Socket{sockFD: fd}, nil
}

func (s *Socket) AllocateVIFs(interfaceNames []string, startVIF uint16) ([]uint16, error) {
	vifs := make([]uint16, len(interfaceNames))
	for i := range interfaceNames {
		vifs[i] = startVIF + uint16(i)
	}
	return vifs, nil
}

func (s *Socket) MulticastInterfaceJoinMgroup(mgroup net.IP, ifaceIP net.IP, ifaceName string) error {
	klog.V(2).InfoS("Joining multicast group", "group", mgroup, "interface", ifaceName, "ip", ifaceIP)
	mreq := &syscall.IPMreq{}
	mIP := mgroup.To4()
	if mIP == nil {
		return fmt.Errorf("multicast group %s is not a valid IPv4 address", mgroup)
	}
	copy(mreq.Multiaddr[:], mIP)
	iIP := ifaceIP.To4()
	if iIP == nil {
		return fmt.Errorf("interface IP %s is not a valid IPv4 address", ifaceIP)
	}
	copy(mreq.Interface[:], iIP)
	err := syscall.SetsockoptIPMreq(s.sockFD, syscall.IPPROTO_IP, syscall.IP_ADD_MEMBERSHIP, mreq)
	if err != nil {
		return fmt.Errorf("failed to join multicast group %s on %s: %v", mgroup, ifaceName, err)
	}
	return nil
}

func (s *Socket) MulticastInterfaceLeaveMgroup(mgroup net.IP, ifaceIP net.IP, ifaceName string) error {
	klog.V(2).InfoS("Leaving multicast group", "group", mgroup, "interface", ifaceName)
	mreq := &syscall.IPMreq{}
	mIP := mgroup.To4()
	if mIP == nil {
		return fmt.Errorf("multicast group %s is not a valid IPv4 address", mgroup)
	}
	copy(mreq.Multiaddr[:], mIP)
	iIP := ifaceIP.To4()
	if iIP == nil {
		return fmt.Errorf("interface IP %s is not a valid IPv4 address", ifaceIP)
	}
	copy(mreq.Interface[:], iIP)
	err := syscall.SetsockoptIPMreq(s.sockFD, syscall.IPPROTO_IP, syscall.IP_DROP_MEMBERSHIP, mreq)
	if err != nil {
		return fmt.Errorf("failed to leave multicast group %s on %s: %v", mgroup, ifaceName, err)
	}
	return nil
}

func (s *Socket) GetFD() int {
	return int(s.sockFD)
}
