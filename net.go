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

	log "github.com/sirupsen/logrus"
)

// AddHost will add a host to topology. set docker to true can enable docker container as a host
func AddHost(name string, addr string, docker bool) (*Host, error) {
	var h *Host
	var err error
	if docker {
		// Create a docker container
		h, err = NewContainer(name, "library/busybox")
		if err != nil {
			log.Fatal("failed to NewContainer: ", err)
		}
	} else {
		// Create a network namespace
		h, err = NewHost(name)
		if err != nil {
			log.Fatal("failed to NewHost: ", err)
		}
	}
	// setup a veth pair
	_, err = h.setupVeth("eth2", 1500)
	if err != nil {
		log.Fatal("failed to setup veth pair: ", err)
	}
	// setup a IP for host
	h.setIfaceIP(addr)
	if err != nil {
		log.Fatalf("failed to setIfaceIP for %s: %v", h.Name, err)
		return nil, err
	}
	return h, nil
}

// AddSwitch will add a switch to topology.
func AddSwitch(params ...string) (*OVSSwitch, error) {
	// params[0] for brName
	// params[1] for controller remote IP and port
	// Create a Open vSwitch bridge
	sw, err := NewOVSSwitch(params[0])
	if err != nil {
		log.Fatal("failed to NewOVSSwitch: ", err)
	}
	if len(params) == 2 {
		if err := sw.setCtrl(params[1]); err != nil {
			log.Warnf("failed to setCtrl for %s: %v", sw.BridgeName, err)
			return nil, err
		}
	} else if len(params) > 2 {
		return nil, fmt.Errorf("Too many arguments")
	}
	return sw, nil
}

// AddLink will add a link between switch to switch or host to switch.
func AddLink(n1, n2 interface{}) error {
	// log.Info(n1.(*OVSSwitch).nodeType)
	// log.Info(n2.(*Host).nodeType)
	var err error
	switch n1.(type) {
	case *OVSSwitch:
		switch n2.(type) {
		case *OVSSwitch:
			tap0, tap1, err := makeVethPair("tap0", "tap1")
			if err != nil {
				log.Fatalf("failed to makeVethPair: %v", err)
			}
			err = n1.(*OVSSwitch).addPort(tap0.Name)
			if err != nil {
				log.Fatalf("failed to addPortswitch - switch: %v", err)
			}
			err = n2.(*OVSSwitch).addPort(tap1.Name)
			if err != nil {
				log.Fatalf("failed to addPort switch - switch: %v", err)
			}
			log.Infof("Adding a link: (%s, %s)", n1.(*OVSSwitch).BridgeName, n2.(*OVSSwitch).BridgeName)
		case *Host:
			err = n1.(*OVSSwitch).addPort(n2.(*Host).vethName)
			if err != nil {
				log.Fatalf("failed to addPort switch - host: %v", err)
			}
			log.Infof("Adding a link: (%s, %s)", n1.(*OVSSwitch).BridgeName, n2.(*Host).Name)
		default:
			log.Fatalf("Type Error")
		}
	case *Host:
		switch n2.(type) {
		case *OVSSwitch:
			err = n2.(*OVSSwitch).addPort(n1.(*Host).vethName)
			if err != nil {
				log.Fatalf("failed to addPort host - switch : %v", err)
			}
			log.Infof("Adding a link: (%s, %s)", n1.(*Host).Name, n2.(*OVSSwitch).BridgeName)
		case *Host:
			log.Fatalf("Type Error: host can not connect to host")
		default:
			log.Fatalf("Type Error")
		}
	default:
		log.Fatalf("Type Error")
	}
	return nil
}

// GetIPs will get IP address from a CIDR notation
// Return an array with IP address. negative num for getting all IP addresses
// (network address and broadcast address are not included)
func GetIPs(cidr string, num int) ([]string, error) {
	ips, err := getAllIPsfromCIDR(cidr)
	if err != nil {
		return nil, err
	}
	if num < 0 {
		return ips, nil
	}
	return ips[0:num], nil
}
