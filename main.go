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
	log.Printf("[DEBUG] zcfw: Starting to watch registry for changes")
	cfg := &config.Config{Zone: "default", Mode: "service"}
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
		log.Printf("[DEBUG] zcfw: Read from service")
		if config == last {
			log.Printf("[DEBUG] zcfw: Same as last time, contine")
			continue
		}
		log.Println(config)
		last = config
	}
}
