package firewall

type Rule struct {
	SourceService      string
	DestinationService string
	SourceIP           string
	DestinationIP      string
	SourcePort         int
	DestinationPort    int
	Rule               string
}
