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
	"strconv"

	"github.com/John-Lin/ovsdb"
	log "github.com/Sirupsen/logrus"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/plugins/pkg/ip"
	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/vishvananda/netlink"
)

func setupVeth(netns ns.NetNS, ifName string, mtu int) (*current.Interface, *current.Interface, error) {
	contIface := &current.Interface{}
	hostIface := &current.Interface{}

	err := netns.Do(func(hostNS ns.NetNS) error {
		// create the veth pair in the container and move host end into host netns
		hostVeth, containerVeth, err := ip.SetupVeth(ifName, mtu, hostNS)
		if err != nil {
			return err
		}
		contIface.Name = containerVeth.Name
		contIface.Mac = containerVeth.HardwareAddr.String()
		contIface.Sandbox = netns.Path()
		hostIface.Name = hostVeth.Name

		// ip link set lo up
		if err := ifaceUp("lo"); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	// need to lookup hostVeth again as its index has changed during ns move
	hostVeth, err := netlink.LinkByName(hostIface.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to lookup %q: %v", hostIface.Name, err)
	}
	hostIface.Mac = hostVeth.Attrs().HardwareAddr.String()

	return hostIface, contIface, nil
}

func addPortToOVS(OVSBrName, hostVethName string) error {
	var ovsDriver *ovsdb.OvsDriver

	ovsDriver = ovsdb.NewOvsDriverWithUnix(OVSBrName)

	// Ask OVSDB driver to add the port
	if !ovsDriver.IsPortNamePresent(hostVethName) {
		err := ovsDriver.CreatePort(hostVethName, "", 0)
		if err != nil {
			return fmt.Errorf("Error creating the port. Err: %v", err)
		}
	}
	return nil
}

func setupOVSBridge(OVSBrName string) (*netlink.Bridge, *current.Interface, error) {
	// create ovs bridge
	ovsbr, err := ensureOVSBridge(OVSBrName)
	if err != nil {
		return nil, nil, err
	}

	return ovsbr, &current.Interface{
		Name: OVSBrName,
	}, nil

	// return ovsbr, nil
}

func setupNamespaceIP(netns ns.NetNS, ifName, address string) error {
	err := netns.Do(func(hostNS ns.NetNS) error {
		if err := ensureifNameAddr(ifName, address); err != nil {
			return err
		}

		// ip link set ifName up
		if err := ifaceUp(ifName); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil

}

func setupCtrlerToOVS(OVSBrName, hostport string) error {
	var ovsDriver *ovsdb.OvsDriver
	ovsDriver = ovsdb.NewOvsDriverWithUnix(OVSBrName)
	// setup OpenFlow controller for ovs bridge

	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		return fmt.Errorf("Invalid controller IP and port. Err: %v", err)
	}

	uPort, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		return fmt.Errorf("Invalid controller port number. Err: %v", err)
	}

	err = ovsDriver.AddController(host, uint16(uPort))
	if err != nil {
		return fmt.Errorf("Error adding controller to OVS. Err: %v", err)
	}
	return nil
}

func setupVethPair() (net.Interface, net.Interface, error) {
	// Create a veth pair
	peerName, err := ip.RandomVethName()
	if err != nil {
		return net.Interface{}, net.Interface{}, err
	}
	v0 := "tap0" + peerName[4:]
	v1 := "tap1" + peerName[4:]
	log.Infof("Create a veth pair: %s : %s", v0, v1)

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

	veth0, err := netlink.LinkByName(v0)
	if err != nil {
		return net.Interface{}, net.Interface{}, err
	}
	if err := netlink.LinkSetUp(veth0); err != nil {
		return net.Interface{}, net.Interface{}, err
	}

	veth1, err := netlink.LinkByName(v1)
	if err != nil {
		return net.Interface{}, net.Interface{}, err
	}
	if err := netlink.LinkSetUp(veth1); err != nil {
		return net.Interface{}, net.Interface{}, err
	}

	return ifaceFromNetlinkLink(veth1), ifaceFromNetlinkLink(veth0), nil
}
