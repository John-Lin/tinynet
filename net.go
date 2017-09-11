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
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/plugins/pkg/ns"
)

func AddHost(addr string) (*current.Interface, error) {
	MTU := 1500
	ifName := "eth2"
	// Create a network namespace
	targetNs, err := ns.NewNS()
	log.Info("netns mouted into the host: ", targetNs.Path())

	// Get network namespace object
	netns, err := ns.GetNS(targetNs.Path())
	if err != nil {
		log.Fatal("failed to open netns: ", err)
	}
	defer netns.Close()

	// Create veth pairs, put into network namespace and enable lookback device
	hostInterface, containerInterface, err := setupVeth(netns, ifName, MTU)
	if err != nil {
		log.Fatal("failed to setupVeth: ", err)
	}

	// setup network namespace IP
	err = setupNamespaceIP(netns, containerInterface.Name, addr)
	if err != nil {
		log.Fatal("failed to setupContIP: ", err)
	}

	return hostInterface, nil
}

func AddSwitch(params ...string) (*current.Interface, error) {
	// params[0] for brName
	// params[0] for controller remote IP and port
	// Create a Open vSwitch bridge
	_, brInterface, err := setupOVSBridge(params[0])
	if err != nil {
		log.Fatal("failed to setupOVSBridge: ", err)
	}

	if len(params) == 2 {
		if err := setupCtrlerToOVS(params[0], params[1]); err != nil {
			log.Fatal("failed to setupCtrlerToOVS: ", err)
			return nil, err
		}
	} else if len(params) > 2 {
		return nil, fmt.Errorf("Too many arguments")
	}

	return brInterface, nil
}

func AddLink(brName, contName string) error {
	if !strings.HasPrefix(contName, "veth") && !strings.HasPrefix(brName, "veth") {
		// switch connect to switch
		tap0, tap1, err := setupVethPair()
		if err != nil {
			return fmt.Errorf("failed to setupVethPair: %v", err)
		}

		err = addPortToOVS(brName, tap0.Name)
		if err != nil {
			return fmt.Errorf("failed to addLink: %v", err)
		}

		err = addPortToOVS(contName, tap1.Name)
		if err != nil {
			return fmt.Errorf("failed to addLink: %v", err)
		}

	} else {
		// container connect to switch
		err := addPortToOVS(brName, contName)
		if err != nil {
			log.Fatal("failed to addPortToOVS: ", err)
			return fmt.Errorf("failed to addPortToOVS: %v", err)
		}
	}
	return nil
}

func AddController(brName, hostport string) error {
	if err := setupCtrlerToOVS(brName, hostport); err != nil {
		log.Fatal("failed to setupCtrlerToOVS: ", err)
		return err
	}
	return nil
}
