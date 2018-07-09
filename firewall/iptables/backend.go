package iptables

import "github.com/ArmedGuy/zcfw/firewall"

type IptablesFirewall struct {
	table string
}

func NewBackend() (firewall.Firewall, error) {
	return &IptablesFirewall{table: ""}, nil
}

func (ip *IptablesFirewall) Register() error {
	return nil
}

func (ip *IptablesFirewall) AddRule(rule *firewall.Rule) error {
	return nil
}

func (ip *IptablesFirewall) RemoveRule(rule *firewall.Rule) error {
	return nil
}
