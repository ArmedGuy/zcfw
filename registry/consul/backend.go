package consul

import (
	"fmt"

	cfg "github.com/ArmedGuy/zcfw/config"
	"github.com/hashicorp/consul/api"
)

type ConsulRegistry struct {
	client *api.Client
}

func NewBackend() (*ConsulRegistry, error) {
	client, err := api.NewClient(&api.Config{
		Address: "127.0.0.1:8500",
	})
	if err != nil {
		return nil, err
	}
	return &ConsulRegistry{client: client}, nil
}

func (c *ConsulRegistry) Register() error {
	return nil
}

func (c *ConsulRegistry) Watch(cfg *cfg.Config) (chan []string, error) {
	svc := make(chan []string)
	switch cfg.Mode {
	case "service":
		go watchIntentions(c.client, cfg, svc)
		break
	case "distributed":
		path := "zcfw/zones/" + cfg.Zone + "/firewall"
		go watchKV(c.client, path, svc, true)
		break
	default:
		return nil, fmt.Errorf("Unknown firewall mode %v", cfg.Mode)
	}
	return svc, nil
}
