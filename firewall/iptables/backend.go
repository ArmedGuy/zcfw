package iptables

import "github.com/ArmedGuy/zcfw/firewall"

type iptables struct {
	table string
}

func NewBackend() (firewall.Firewall, error) {
	return &iptables{table: ""}, nil
}

func (ip *iptables) Register() error {
	return nil
}

func (ip *iptables) AddRule(rule *firewall.Rule) error {
	return nil
}

func (ip *iptables) RemoveRule(rule *firewall.Rule) error {
	return nil
}
