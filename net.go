package muts

import (
	"net"
	"strconv"
	"strings"
)

// LocalIP returns not local IP address (not being a loopback).
// HOSTIP=`ifconfig eth0 | grep "inet addr" | awk -F: '{print $2}' | awk '{print $1}'`
func LocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		Fatalln("unable to get interface addresses", err)
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// FreePort returns a free, usable TCP port (in practice).
func FreePort() int {
	l, err := net.Listen("tcp", "")
	if err != nil {
		Fatalln("unable to get listener for tcp", err)
	}
	defer l.Close()
	t := strings.Split(l.Addr().String(), ":")
	i, _ := strconv.Atoi(t[len(t)-1])
	return i
}

// Port registers and returns a free TCP port from the OS or the argument value if the -local flag was set.
func Port(label string, local int) int {
	if *LocalUse {
		PortRegistry[label] = local
		return local
	}
	if i, ok := PortRegistry[label]; ok {
		return i
	}
	i := FreePort()
	PortRegistry[label] = i
	return i
}

// PortRegistry holds a mapping for resources and their assigned TCP ports
var PortRegistry = map[string]int{}
