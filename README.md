# tinynet
A lightweight instant virtual network for rapid prototyping SDN 


## Prerequisite
[**Open vSwitch**](http://openvswitch.org/)

Installation from Debian Packages
```
$ sudo apt-get install openvswitch-switch
```

or [installation from source](http://docs.openvswitch.org/en/latest/intro/install/#installation-from-source)

## Install 

```
$ go get -u github.com/John-Lin/tinynet
```

## Example
```go
package main

import (
	tn "github.com/John-Lin/tinynet"
	log "github.com/Sirupsen/logrus"
)

func main() {
	// add Switch1
	sw1, err := tn.AddSwitch("br0")
	if err != nil {
		log.Fatal("failed to AddSwitch:", err)
	}

	// add Switch2 and set a OpenFlow controller
	sw2, err := tn.AddSwitch("br1", "128.199.225.69:6633")
	if err != nil {
		log.Fatal("failed to AddSwitch:", err)
	}

	// add Hosts1
	h1, err := tn.AddHost("10.0.0.102/24")
	if err != nil {
		log.Fatal("failed to AddHost:", err)
	}

	// add Hosts2
	h2, err := tn.AddHost("10.0.0.101/24")
	if err != nil {
		log.Fatal("failed to AddHost:", err)
	}

	// add Link sw1 - h1
	if err := tn.AddLink(sw1.Name, h1.Name); err != nil {
		log.Fatal("failed to AddLink:", err)
	}

	// add Link sw2 - h2
	if err := tn.AddLink(sw2.Name, h2.Name); err != nil {
		log.Fatal("failed to AddLink:", err)
	}

	// add Link sw1 - sw2
	if err := tn.AddLink(sw1.Name, sw2.Name); err != nil {
		log.Fatal("failed to AddLink:", err)
	}

}
```
