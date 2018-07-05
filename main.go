package main

import (
	"log"
	"os"

	"github.com/ArmedGuy/zcfw/config"
	"github.com/ArmedGuy/zcfw/firewall"
	"github.com/ArmedGuy/zcfw/firewall/iptables"
	"github.com/ArmedGuy/zcfw/registry"
	"github.com/ArmedGuy/zcfw/registry/consul"
)

func main() {
	initRegistry()
	initFirewall()

	watchRegistry()
}

func initRegistry() {
	var err error
	registry.Default, err = consul.NewBackend()
	if err != nil {
		log.Printf("[ERROR] zcfw: Failed to initialize consul backend. %v", err)
		os.Exit(1)
	}
}
func initFirewall() {
	var err error
	firewall.Default, err = iptables.NewBackend()
	if err != nil {
		log.Printf("[ERROR] zcfw: Failed to initialize firewall backend. %v", err)
		os.Exit(1)
	}
}

func watchRegistry() {
	cfg := &config.Config{Zone: "default", Mode: "distributed"}
	svc, err := registry.Default.Watch(cfg)
	if err != nil {
		log.Printf("[ERROR] zcfw: Failed to start watching registry. %v", err)
	}
	var (
		last   string
		config string
	)
	for {
		config = <-svc
		if config == last {
			continue
		}
		log.Printf(config)
		last = config
	}
}
