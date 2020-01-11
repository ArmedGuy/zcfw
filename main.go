package main

import (
	"log"
	"os"

	"github.com/ArmedGuy/zcfw/config"
	"github.com/ArmedGuy/zcfw/firewall"
	"github.com/ArmedGuy/zcfw/firewall/nftables"
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
	firewall.Default, err = nftables.NewBackend()
	if err != nil {
		log.Printf("[ERROR] zcfw: Failed to initialize firewall backend. %v", err)
		os.Exit(1)
	}
}

func ConfigEquals(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func watchRegistry() {
	log.Printf("[DEBUG] zcfw: Starting to watch registry for changes")
	cfg := &config.Config{Zone: "default", Mode: "service"}
	svc, err := registry.Default.Watch(cfg)
	if err != nil {
		log.Printf("[ERROR] zcfw: Failed to start watching registry. %v", err)
	}
	var (
		last   []string
		config []string
	)

	for {
		config = <-svc
		log.Printf("[DEBUG] zcfw: Read from service")
		if ConfigEquals(config, last) {
			log.Printf("[DEBUG] zcfw: Same as last time, continue")
			continue
		}
		firewall.Default.SetRules(config)
		last = config
	}
}
