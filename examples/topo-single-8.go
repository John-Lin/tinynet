package main

import (
	tn "github.com/John-Lin/tinynet"
	log "github.com/sirupsen/logrus"
	"strconv"
)

// Custom topology example
// Single switch with 8 hosts

func main() {
	// add a switch as a Switch
	Switch, err := tn.AddSwitch("br0")
	if err != nil {
		log.Fatal("failed to add Switch:", err)
	}

	// add 8 IPs with CIDR 192.168.1.0/24
	ips, _ := tn.GetIPs("192.168.1.0/24", 8)

	// add 8 hosts and link to switch
	for (i := 0; i < 7; i++) {
		hostName := "h" + strconv.Itoa(i)
		tn.AddHost(hostName, ips[i])
		if err != nil {
			log.Fatal("failed to add host:", err)
		}

		// add Link for leftHost - leftSwitch
		if err := tn.AddLink(hostName, Switch); err != nil {
			log.Fatal("failed to add link between Switch and host: ", err)
		}
	}
}
