package firewall

import (
	"strconv"
	"strings"
)


type Rule struct {
	Allow bool
	SrcAddr string 
	SrcPort int
	DestAddr string
	DestPort int
}

func NewRule(raw string) Rule {
	parts := strings.Split(raw, " ")
	r := Rule{}
	r.Allow = parts[1] == "allow"
	r.SrcAddr = parts[2]
	r.SrcPort, _ = strconv.Atoi(parts[3])
	r.DestAddr = parts[5]
	r.DestPort, _ = strconv.Atoi(parts[6])
	return r
}
