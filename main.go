package main

import (
	"github.com/ArmedGuy/zcfw/firewall"
	"github.com/ArmedGuy/zcfw/firewall/iptables"
)

func main() {
	initBackend()
}

func initBackend() {
	firewall.Default = iptables.NewBackend()
}
