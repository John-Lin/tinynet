// Copyright (c) 2017 Che Wei, Lin
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

package tinynet

import (
	"crypto/rand"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

func ifaceUp(ifName string) (netlink.Link, error) {
	// finds a link by name and returns a pointer to the object.
	l, err := netlink.LinkByName(ifName)
	if err != nil {
		return nil, fmt.Errorf("could not lookup link on ifaceUp %q: %v", l, err)
	}
	// enables the link device
	if err := netlink.LinkSetUp(l); err != nil {
		return nil, err
	}
	return l, nil
}

func bridgeByName(name string) (*netlink.Bridge, error) {
	l, err := netlink.LinkByName(name)
	if err != nil {
		return nil, fmt.Errorf("could not lookup %q: %v", name, err)
	}
	br, ok := l.(*netlink.Bridge)
	if !ok {
		return nil, fmt.Errorf("%q already exists but is not a bridge", name)
	}
	return br, nil
}

func setIP(ifName, address string) error {
	// finds a link by name and returns a pointer to the object.
	l, err := netlink.LinkByName(ifName)
	if err != nil {
		return fmt.Errorf("could not lookup link on setIP %q: %v", l, err)
	}
	addr, _ := netlink.ParseAddr(address)
	if err := netlink.AddrAdd(l, addr); err != nil {
		return fmt.Errorf("could not add IP address to %q: %v", ifName, err)
	}
	return nil
}

func ifaceFromNetlinkLink(l netlink.Link) net.Interface {
	a := l.Attrs()
	return net.Interface{
		Index:        a.Index,
		MTU:          a.MTU,
		Name:         a.Name,
		HardwareAddr: a.HardwareAddr,
		Flags:        a.Flags,
	}
}

// Create a veth pair
func makeVethPair(ifName1, ifName2 string) (net.Interface, net.Interface, error) {
	peerName, err := pseudoRandomName()
	if err != nil {
		return net.Interface{}, net.Interface{}, err
	}
	v0 := ifName1 + peerName
	v1 := ifName2 + peerName

	vethLink := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{
			Name:  v0,
			Flags: net.FlagUp,
		},
		PeerName: v1,
	}
	if err := netlink.LinkAdd(vethLink); err != nil {
		return net.Interface{}, net.Interface{}, err
	}

	// ip link set v0 up
	veth1, err := ifaceUp(v0)
	if err != nil {
		return net.Interface{}, net.Interface{}, err
	}
	// ip link set v1 up
	veth0, err := ifaceUp(v1)
	if err != nil {
		return net.Interface{}, net.Interface{}, err
	}
	log.Infof("Created a veth pair: %s : %s", v0, v1)
	return ifaceFromNetlinkLink(veth1), ifaceFromNetlinkLink(veth0), nil
}

// Generate pseudo random string with length 8, but not Universally unique.
func pseudoRandomName() (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random veth name: %v", err)
	}
	return fmt.Sprintf("%x", b), nil
}
