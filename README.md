# tinynet
A lightweight instant virtual network for rapid prototyping SDN 

Tinynet provides a simple network library with Go and uses for emulating any network to prototype Software Defined Networks.

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

## API Reference

https://godoc.org/github.com/John-Lin/tinynet

## Example
```go
package main

import (
	tn "github.com/John-Lin/tinynet"
	log "github.com/Sirupsen/logrus"
)

// Custom topology example
// Host1 --- Switch1 --- host2

func main() {
	// add a switch as a Switch1
	sw1, err := tn.AddSwitch("br0")
	if err != nil {
		log.Fatal("failed to AddSwitch:", err)
	}

	// add a host as a Host1
	h1, err := tn.AddHost("10.0.0.102/24")
	if err != nil {
		log.Fatal("failed to AddHost:", err)
	}
	// add a host as a Host2
	h2, err := tn.AddHost("10.0.0.101/24")
	if err != nil {
		log.Fatal("failed to AddHost:", err)
	}

	// add Link for Switch1 - Host1
	if err := tn.AddLink(sw1.Name, h1.Name); err != nil {
		log.Fatal("failed to AddLink:", err)
	}
	// add Link for Switch1 - Host2
	if err := tn.AddLink(sw1.Name, h2.Name); err != nil {
		log.Fatal("failed to AddLink:", err)
	}
}
```
more complicated example that you might see in a [examples](https://github.com/John-Lin/tinynet/tree/master/examples) folder

## Cleanup command

```
$ sudo python clean.py
```