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
	"net"

	"github.com/containernetworking/plugins/pkg/ip"
	"github.com/containernetworking/plugins/pkg/ns"
	log "github.com/sirupsen/logrus"
)

// Host is a host instance
type Host struct {
	NodeType string
	Name     string
	vethName string
	IfName   string
	Sandbox  string
	IP       string
	MAC      string
}

// NewHost for creating a network namespace
func NewHost(name string) (*Host, error) {
	h := new(Host)
	h.NodeType = "Host"
	h.Name = name

	// Create a network namespace
	targetNs, err := ns.NewNS()
	if err != nil {
		log.Fatal("failed to open netns: ", err)
	}
	// log.Info("netns mouted into the host: ", targetNs.Path())
	log.Infof("%s : %s\n", h.Name, targetNs.Path())

	h.Sandbox = targetNs.Path()

	return h, nil
}

func (h *Host) setupVeth(ifName string, mtu int) (*Host, error) {
	// Get network namespace object
	netns, err := ns.GetNS(h.Sandbox)
	if err != nil {
		log.Fatal("failed to open netns: ", err)
	}
	defer netns.Close()

	// attach network namespace and setup veth pair
	err = netns.Do(func(hostNS ns.NetNS) error {
		// create the veth pair in the container and move host end into host netns
		hostVeth, containerVeth, err := ip.SetupVeth(ifName, mtu, hostNS)
		if err != nil {
			return err
		}
		// Host interface name
		h.IfName = containerVeth.Name

		// h.mac = containerVeth.HardwareAddr.String()

		// Host veth name
		h.vethName = hostVeth.Name

		// ip link set lo up
		_, err = ifaceUp("lo")
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (h *Host) setIfaceIP(address string) error {
	// Get network namespace object
	netns, err := ns.GetNS(h.Sandbox)
	if err != nil {
		log.Fatal("failed to open netns: ", err)
	}
	defer netns.Close()

	ipv4Addr, _, err := net.ParseCIDR(address)
	if err != nil {
		return err
	}
	h.IP = ipv4Addr.String()
	err = netns.Do(func(hostNS ns.NetNS) error {
		if err := setIP(h.IfName, address); err != nil {
			return err
		}
		// ip link set ifName up
		_, err := ifaceUp(h.IfName)
		if err != nil {
			return err
		}
		// get Host mac address
		hostIface, err := net.InterfaceByName(h.IfName)
		if err != nil {
			return err
		}
		h.MAC = hostIface.HardwareAddr.String()
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
