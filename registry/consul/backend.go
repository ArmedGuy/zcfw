package consul

import (
	cfg "github.com/ArmedGuy/zcfw/config"
	"github.com/hashicorp/consul/api"
)

type consul struct {
	client *api.Client
}

func NewBackend() (*consul, error) {
	client, err := api.NewClient(&api.Config{
		Address: "127.0.0.1",
	})
	if err != nil {
		return nil, err
	}
	return &consul{client: client}, nil
}

func (c *consul) Register() error {
	return nil
}

func (c *consul) Watch(cfg *cfg.Config) (chan string, error) {
	svc := make(chan string)
	if cfg.Mode == "service" {
		watchIntentions(c.client, cfg, svc)
	} else {
		path := "zcfw/zones/" + cfg.Zone + "/firewall"
		go watchKV(c.client, path, svc, true)
	}
	return svc, nil
}
