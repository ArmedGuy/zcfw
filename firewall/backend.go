package firewall

type Firewall interface {
	// Register to a firewall
	Register() error

	SetRules(rules []string)
}

var Default Firewall
