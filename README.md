# tinynet
[![Build Status](https://api.travis-ci.org/John-Lin/tinynet.svg?branch=develop)](https://travis-ci.org/John-Lin/tinynet)
[![codecov](https://codecov.io/gh/John-Lin/tinynet/branch/develop/graph/badge.svg)](https://codecov.io/gh/John-Lin/tinynet)
[![GoDoc](https://godoc.org/github.com/John-Lin/tinynet?status.svg)](https://godoc.org/github.com/John-Lin/tinynet)

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

## Features
- Go API for creating any networks topology
- Integrate with docker containers as a host
- A `cleanup` command for removing virtual interfaces, network namespace and bridges

## Example
```go
package main

import (
	tn "github.com/John-Lin/tinynet"
	log "github.com/sirupsen/logrus"
)

// Custom topology example
// Host1 --- Switch1 --- Host2

func main() {
	// add a switch as a Switch1
	sw1, err := tn.AddSwitch("br0")
	if err != nil {
		log.Fatal("failed to AddSwitch:", err)
	}
	// setup Host1 network configuration through NewHostConfig constructor
	// parameters: name, address, interface name, MTU, isDocker, docker image
	hc1 := tn.NewHostConfig("h1", "10.0.0.102/24", "eth2", 1500, true, "library/busybox")
	// add a docker container host as a Host1
	h1, err := tn.AddHostWithConf(hc1)
	if err != nil {
		log.Fatal("failed to AddHost:", err)
	}
	// setup Host2 network configuration
	hc2 := tn.NewHostConfig("h2", "10.0.0.101/24", "eth1", 1500, true, "library/busybox")
	// add a docker container host as a Host2
	h2, err := tn.AddHostWithConf(hc2)
	if err != nil {
		log.Fatal("failed to AddHost:", err)
	}

	// add Link for Switch1 - Host1
	if err := tn.AddLink(sw1, h1); err != nil {
		log.Fatal("failed to AddLink:", err)
	}
	// add Link for Host2 - Switch1
	if err := tn.AddLink(h2, sw1); err != nil {
		log.Fatal("failed to AddLink:", err)
	}
}
```
more complicated example that you might see in a [examples](https://github.com/John-Lin/tinynet/tree/master/examples) folder

## Cleanup command

`cleanup` command to remove interfaces, network namespaces and bridges.

```
$ sudo python clean.py
```
