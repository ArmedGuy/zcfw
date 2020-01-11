package nftables

import (
	"log"
	"os"
	"os/exec"
	"github.com/ArmedGuy/zcfw/firewall"
	"text/template"
)

type NftablesFirewall struct {
	template string
	Rules []firewall.Rule
}

func NewBackend() (firewall.Firewall, error) {
	return &NftablesFirewall{template: "/etc/zcfw/nftables.template"}, nil
}

func (nf *NftablesFirewall) Register() error {
	return nil
}

func (nf *NftablesFirewall) SetRules(rules []string) {
	nf.Rules = make([]firewall.Rule, len(rules))
	for i, rule := range rules {
		nf.Rules[i] =  firewall.NewRule(rule)
	}

	tmpl, err := template.ParseFiles(nf.template)
	if err != nil {
		log.Printf("[ERROR] nftables: Unable to parse template for firewall: %v", err)
		return
	}

	f, err := os.Create("/tmp/zcfw_nftables")
	if err != nil {
		log.Printf("[ERROR] nftables: Unable to open firewall file for writing: %v", err)
		return
	}
	err = tmpl.Execute(f, nf)
	if err != nil {
		log.Printf("[ERROR] nftables: Unable to write firewall config: %v", err)
		return
	}
	f.Close()

	err = exec.Command("nft", "-c", "-f", "/tmp/zcfw_nftables").Run()
	if err != nil {
		log.Printf("[ERROR] nftables: Unable to update firewall config: %v", err)
		return
	}

}