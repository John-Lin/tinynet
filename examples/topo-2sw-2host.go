package main

import (
	tn "github.com/John-Lin/tinynet"
	log "github.com/Sirupsen/logrus"
)

// Custom topology example
// leftHost --- leftSwitch --- rightSwitch --- rightHosts

func main() {
	// add a switch as a leftSwitch
	leftSwitch, err := tn.AddSwitch("br0")
	if err != nil {
		log.Fatal("failed to add leftSwitch: ", err)
	}
	// add a switch as a rightSwitch and set a OpenFlow controller
	rightSwitch, err := tn.AddSwitch("br1", "192.168.8.100:6633")
	if err != nil {
		log.Fatal("failed to add rightSwitch: ", err)
	}

	// add a host as a leftHost
	leftHost, err := tn.AddHost("10.0.0.102/24")
	if err != nil {
		log.Fatal("failed to add leftHost: ", err)
	}
	// add a host as a rightHost
	rightHost, err := tn.AddHost("10.0.0.101/24")
	if err != nil {
		log.Fatal("failed to rightHost: ", err)
	}

	// add Link for leftSwitch - leftHost
	if err := tn.AddLink(leftSwitch.Name, leftHost.Name); err != nil {
		log.Fatal("failed to add link between leftSwitch and leftHost: ", err)
	}
	// add Link for rightSwitch - rightHost
	if err := tn.AddLink(rightSwitch.Name, rightHost.Name); err != nil {
		log.Fatal("failed to add link between rightSwitch and rightHost: ", err)
	}
	// add Link for leftSwitch - rightSwitch
	if err := tn.AddLink(leftSwitch.Name, rightSwitch.Name); err != nil {
		log.Fatal("failed to add link between leftSwitch and rightSwitch: ", err)
	}

}
