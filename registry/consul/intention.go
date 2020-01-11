package consul

import (
	"log"
	"strconv"
	"strings"
	"time"

	cfg "github.com/ArmedGuy/zcfw/config"
	"github.com/hashicorp/consul/api"
)

func watchIntentions(client *api.Client, cfg *cfg.Config, svc chan []string) {
	var lastIndex uint64
	for {
		value, index, err := listIntentions(client, cfg.Zone, lastIndex)
		if err != nil {
			log.Printf("[WARN] consul: Error fetching config from intentions. %v", err)
			time.Sleep(time.Second)
			continue
		}
		if index != lastIndex {
			log.Printf("[DEBUG] consul: Distributed config changed to #%d.", index)
			svc <- value
			lastIndex = index
		}
	}
}

func listIntentions(client *api.Client, zone string, waitIndex uint64) ([]string, uint64, error) {
	q := &api.QueryOptions{WaitTime: time.Second * 15}
	intentions, meta, err := client.Connect().Intentions(q)
	if err != nil {
		return nil, 0, err
	}
	if len(intentions) == 0 {
		return nil, meta.LastIndex, nil
	}
	var rows []string
	for _, intention := range intentions {
		if strings.HasPrefix(intention.Description, "#zcfw") {
			rules, err := buildRulesFromIntention(client, intention)
			if err != nil {
				return nil, 0, err
			}
			rows = append(rows, rules...)
		}
	}
	return rows, meta.LastIndex, nil

}

func buildRulesFromIntention(client *api.Client, intention *api.Intention) ([]string, error) {
	var rules []string
	var sources []string
	var destinations []string

	args := buildArgs(intention.Description)

	// Build sources map
	if intention.SourceName == "*" {
		if addr, ok := args["src-addr"]; ok {
			if port, ok := args["src-port"]; ok {
				sources = append(sources, addr+" port "+port)
			} else {
				sources = append(sources, addr+"port 0")
			}
		} else {
			sources = append(sources, ":: port 0")
		}
	} else {
		srcServices, err := getServiceInfo(client, intention.SourceName)
		if err != nil {
			return nil, err
		}
		for _, src := range srcServices {
			svcAddr := src.ServiceAddress
			svcPort := strconv.Itoa(src.ServicePort)
			if val, ok := args["src-addr"]; ok {
				svcAddr = val
			}
			if val, ok := args["src-port"]; ok {
				svcPort = val
			}

			sources = append(sources, svcAddr+" port "+svcPort)
		}
	}

	// Build destination map
	if intention.DestinationName == "*" {
		if addr, ok := args["dest-addr"]; ok {
			if port, ok := args["dest-port"]; ok {
				destinations = append(destinations, addr+" port "+port)
			} else {
				destinations = append(destinations, addr+" port 0")
			}
		} else {
			destinations = append(destinations, ":: port 0")
		}
	} else {
		destServices, err := getServiceInfo(client, intention.DestinationName)
		if err != nil {
			return nil, err
		}
		for _, dest := range destServices {
			svcAddr := dest.ServiceAddress
			svcPort := strconv.Itoa(dest.ServicePort)
			if val, ok := args["dest-addr"]; ok {
				svcAddr = val
			}
			if val, ok := args["dest-port"]; ok {
				svcPort = val
			}
			destinations = append(destinations, svcAddr+" port "+svcPort)
		}
	}

	// Build rules
	for _, dest := range destinations {
		for _, src := range sources {
			// Todo append other args to this rule
			rules = append(rules, "rule "+string(intention.Action)+" "+src+" -> "+dest)
		}
	}
	rules = removeDuplicates(rules)
	return rules, nil
}

func buildArgs(description string) map[string]string {
	ret := make(map[string]string)
	if !strings.HasPrefix(description, "#zcfw:") {
		return ret
	}
	noPrefix := string(([]rune(description)[6:]))

	if strings.Contains(noPrefix, ",") {
		parts := strings.Split(noPrefix, ",")
		for _, part := range parts {
			if !strings.Contains(part, "=") {
				continue
			}
			kv := strings.Split(strings.TrimSpace(part), "=")
			ret[kv[0]] = kv[1]
		}
	} else {
		if !strings.Contains(noPrefix, "=") {
			return ret
		}
		kv := strings.Split(strings.TrimSpace(noPrefix), "=")
		ret[kv[0]] = kv[1]
	}
	return ret
}

func removeDuplicates(elements []string) []string {

	found := make(map[string]bool)
	var res []string

	for v := range elements {
		if found[elements[v]] != true {
			found[elements[v]] = true
			res = append(res, elements[v])
		}
	}
	return res
}
