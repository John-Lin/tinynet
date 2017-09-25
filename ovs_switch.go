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
	"strconv"
	"time"

	"github.com/John-Lin/ovsdb"
	log "github.com/sirupsen/logrus"
)

// OVSSwitch is a bridge instance
type OVSSwitch struct {
	NodeType     string
	BridgeName   string
	CtrlHostPort string
	ovsdb        *ovsdb.OvsDriver
}

// NewOVSSwitch for creating a ovs bridge
func NewOVSSwitch(bridgeName string) (*OVSSwitch, error) {
	sw := new(OVSSwitch)
	sw.NodeType = "OVSSwitch"
	sw.BridgeName = bridgeName

	sw.ovsdb = ovsdb.NewOvsDriverWithUnix(bridgeName)
	log.Infoln(sw.BridgeName)

	// Check if port is already part of the OVS and add it
	if !sw.ovsdb.IsPortNamePresent(bridgeName) {
		// Create an internal port in OVS
		err := sw.ovsdb.CreatePort(bridgeName, "internal", 0)
		if err != nil {
			log.Fatalf("Error creating the internal port. Err: %v", err)
			return nil, err
		}
	}

	time.Sleep(300 * time.Millisecond)
	log.Infof("Waiting for OVS bridge %s setup", bridgeName)

	// ip link set ovs up
	_, err := ifaceUp(bridgeName)
	if err != nil {
		return nil, err
	}

	return sw, nil
}

// addPort for asking OVSDB driver to add the port
func (sw *OVSSwitch) addPort(ifName string) error {
	if !sw.ovsdb.IsPortNamePresent(ifName) {
		err := sw.ovsdb.CreatePort(ifName, "", 0)
		if err != nil {
			log.Fatalf("Error creating the port. Err: %v", err)
			return err
		}
	}
	return nil
}

// setCtrl for seting up OpenFlow controller for ovs bridge
func (sw *OVSSwitch) setCtrl(hostport string) error {
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		log.Fatalf("Invalid controller IP and port. Err: %v", err)
		return err
	}
	uPort, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		log.Fatalf("Invalid controller port number. Err: %v", err)
		return err
	}
	err = sw.ovsdb.AddController(host, uint16(uPort))
	if err != nil {
		log.Fatalf("Error adding controller to OVS. Err: %v", err)
		return err
	}
	sw.CtrlHostPort = hostport
	return nil
}
