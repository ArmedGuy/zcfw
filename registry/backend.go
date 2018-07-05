package registry

import (
	cfg "github.com/ArmedGuy/zcfw/config"
)

type Registry interface {
	Register() error
	Watch(*cfg.Config) (chan string, error)
}

var Default Registry
