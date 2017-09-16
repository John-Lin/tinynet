package main

import (
	tn "github.com/John-Lin/tinynet"
	log "github.com/sirupsen/logrus"
)

// Custom topology example
// leftHost --- leftSwitch --- rightSwitch --- rightHosts

func main() {
	// add a switch as a leftSwitch
	leftSwitch, err := tn.AddSwitch("br0")
	if err != nil {
		log.Fatal("failed to add leftSwitch:", err)
	}
	// add a switch as a Switch1
	rightSwitch, err := tn.AddSwitch("br1")
	if err != nil {
		log.Fatal("failed to add rightSwitch:", err)
	}

	// add a host as a Host1
	leftHost, err := tn.AddHost("10.0.0.102/24")
	if err != nil {
		log.Fatal("failed to add leftHost:", err)
	}
	// add a host as a Host2
	rightHost, err := tn.AddHost("10.0.0.101/24")
	if err != nil {
		log.Fatal("failed to add rightHost:", err)
	}

	// add Link for leftHost - leftSwitch
	if err := tn.AddLink(leftHost, leftSwitch); err != nil {
		log.Fatal("failed to add link between leftSwitch and leftHost: ", err)
	}
	// add Link for leftSwitch - rightSwitch
	if err := tn.AddLink(leftSwitch, rightSwitch); err != nil {
		log.Fatal("failed to add link between leftSwitch and rightSwitch: ", err)
	}
	// add Link for rightSwitch - rightHost
	if err := tn.AddLink(rightSwitch, rightHost); err != nil {
		log.Fatal("failed to add link between rightSwitch and rightHost: ", err)
	}
}
