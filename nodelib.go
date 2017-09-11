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
	"fmt"
	"net"
	"time"

	"github.com/John-Lin/ovsdb"
	"github.com/vishvananda/netlink"
)

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

func ifaceUp(ifName string) error {
	// finds a link by name and returns a pointer to the object.
	deviceLink, err := netlink.LinkByName(ifName)
	if err != nil {
		return fmt.Errorf("could not lookup link on ifaceUp %q: %v", deviceLink, err)
	}

	// enables the link device
	if err := netlink.LinkSetUp(deviceLink); err != nil {
		return err
	}

	return nil
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

func ensureOVSBridge(OVSBrName string) (*netlink.Bridge, error) {
	var ovsDriver *ovsdb.OvsDriver
	// create a ovs bridge
	ovsDriver = ovsdb.NewOvsDriverWithUnix(OVSBrName)

	// Check if port is already part of the OVS and add it
	if !ovsDriver.IsPortNamePresent(OVSBrName) {
		// Create an internal port in OVS
		err := ovsDriver.CreatePort(OVSBrName, "internal", 0)
		if err != nil {
			return nil, fmt.Errorf("Error creating the internal port. Err: %v", err)
		}
	}

	time.Sleep(300 * time.Millisecond)

	// finds a link by name and returns a pointer to the object.
	// ovsbr, _ := netlink.LinkByName(OVSBrName)
	ovsbrLink, err := netlink.LinkByName(OVSBrName)
	if err != nil {
		return nil, fmt.Errorf("could not lookup link on ensureOVSBridge %q: %v", OVSBrName, err)
	}

	// enables the link device
	if err := netlink.LinkSetUp(ovsbrLink); err != nil {
		return nil, err
	}

	ovsbr, _ := bridgeByName(OVSBrName)

	return ovsbr, nil
}

func ensureifNameAddr(ifName, address string) error {
	// finds a link by name and returns a pointer to the object.
	deviceLink, err := netlink.LinkByName(ifName)
	if err != nil {
		return fmt.Errorf("could not lookup link on ensureifNameAddr %q: %v", deviceLink, err)
	}
	addr, _ := netlink.ParseAddr(address)

	if err := netlink.AddrAdd(deviceLink, addr); err != nil {
		return fmt.Errorf("could not add IP address to %q: %v", ifName, err)
	}
	return nil
}
