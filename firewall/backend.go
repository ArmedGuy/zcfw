package firewall

type Firewall interface {
	// Register to a firewall
	Register() error

	// AddRule adds a new rule in the firewall
	AddRule(*Rule) error

	// RemoveRule removes a rule from the firewall
	RemoveRule(*Rule) error
}

var Default Firewall
