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

package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/containernetworking/plugins/pkg/ip"
	"github.com/containernetworking/plugins/pkg/ns"
)

// Host is a host instance
type Host struct {
	nodeType string
	name     string
	ifName   string
	sandbox  string
}

// NewHost for creating a network namespace
func NewHost(name string) (*Host, error) {
	h := new(Host)
	h.nodeType = "Host"
	h.name = name

	// Create a network namespace
	targetNs, err := ns.NewNS()
	if err != nil {
		log.Fatal("failed to open netns: ", err)
	}
	log.Info("netns mouted into the host: ", targetNs.Path())

	h.sandbox = targetNs.Path()

	return h, nil
}

func (h *Host) setupVeth(sandbox string, ifName string, mtu int) (*Host, error) {
	// Get network namespace object
	netns, err := ns.GetNS(sandbox)
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
		h.ifName = containerVeth.Name
		// Host mac address
		// h.mac = containerVeth.HardwareAddr.String()

		// Host name
		h.name = hostVeth.Name

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

func (h *Host) setIfaceIP(sandbox, ifName, address string) error {
	// Get network namespace object
	netns, err := ns.GetNS(sandbox)
	if err != nil {
		log.Fatal("failed to open netns: ", err)
	}
	defer netns.Close()

	err = netns.Do(func(hostNS ns.NetNS) error {
		if err := setIP(ifName, address); err != nil {
			return err
		}
		// ip link set ifName up
		_, err := ifaceUp(ifName)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
