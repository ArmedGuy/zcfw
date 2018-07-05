package consul

import "github.com/hashicorp/consul/api"

func getServiceInfo(client *api.Client, service string) ([]*api.CatalogService, error) {
	q := &api.QueryOptions{RequireConsistent: true}
	svcs, _, err := client.Catalog().Service(service, "", q)
	if err != nil {
		return nil, err
	}
	return svcs, nil
}
