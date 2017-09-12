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
	"fmt"

	log "github.com/Sirupsen/logrus"
)

func AddHost(addr string) (*Host, error) {
	// Create a network namespace
	h, err := NewHost("h1")
	if err != nil {
		log.Fatal("failed to NewHost: ", err)
	}
	// setup a veth pair
	_, err = h.setupVeth(h.sandbox, "eth2", 1500)
	if err != nil {
		log.Fatal("failed to open netns: ", err)
	}
	// setup a IP for host
	h.setIfaceIP(h.sandbox, h.ifName, addr)
	if err != nil {
		log.Fatalf("failed to setIfaceIP for %s: %v", h.name, err)
		return nil, err
	}
	return h, nil
}

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
			log.Warnf("failed to setCtrl for %s: %v", sw.bridgeName, err)
			return nil, err
		}
	} else if len(params) > 2 {
		return nil, fmt.Errorf("Too many arguments")
	}
	return sw, nil
}

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
		case *Host:
			err = n1.(*OVSSwitch).addPort(n2.(*Host).name)
			if err != nil {
				fmt.Printf("failed to addPort switch - host: %v", err)
			}
		default:
			log.Fatalf("Type Error")
		}
	case *Host:
		switch n2.(type) {
		case *OVSSwitch:
			err = n2.(*OVSSwitch).addPort(n1.(*Host).name)
			if err != nil {
				fmt.Printf("failed to addPort host - switch : %v", err)
			}
		case *Host:
			log.Fatalf("Type Error: host can not connect to host")
		default:
			log.Fatalf("Type Error")
		}
	default:
		log.Fatalf("Type Error")
	}

	// if n1.NodeType == "OVSSwitch" && n2.NodeType == "OVSSwitch" {
	// 	// switch - switch
	// 	tap0, tap1, err := makeVethPair("tap0", "tap1")
	// 	if err != nil {
	// 		log.Fatalf("failed to makeVethPair: %v", err)
	// 	}
	// 	err = n1.addPort(tap0.Name)
	// 	if err != nil {
	// 		log.Fatalf("failed to addPortswitch - switch: %v", err)
	// 	}
	// 	err = n2.addPort(tap1.Name)
	// 	if err != nil {
	// 		log.Fatalf("failed to addPort switch - switch: %v", err)
	// 	}
	// } else {
	// 	// switch - host || host - switch
	// 	err = n1.addPort(n2.name)
	// 	if err != nil {
	// 		fmt.Printf("failed to addPort switch - host: %v", err)
	// 	}
	// }
	return nil
}

// func switch2switch(n1, n2 *OVSSwitch) {
// 	var err error
// 	var tap [2]net.Interface
// 	var n [2]*OVSSwitch

// 	n[0] = n1
// 	n[1] = n2

// 	tap[0], tap[1], err = makeVethPair("tap0", "tap1")
// 	if err != nil {
// 		log.Fatalf("failed to makeVethPair: %v", err)
// 	}
// 	for i := 0; i < 1; i++ {
// 		err = n[i].addPort(tap[i].Name)
// 		if err != nil {
// 			log.Fatalf("failed to addPort switch - switch: %v", err)
// 		}
// 	}
// 	// err = n2.addPort(tap1.Name)
// 	// if err != nil {
// 	// 	log.Fatalf("failed to addPort switch - switch: %v", err)
// 	// }
// 	return nil
// }

// func switch2host() {
// 	return n1.addPort(n2.name)
// }
