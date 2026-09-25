package keystone

import "fmt"

func (c *Client) Endpoint(serviceType, iface, region string) (string, error) {
	if c.token == nil {
		return "", fmt.Errorf("keystone: no token")
	}

	for _, svc := range c.token.ServiceCatalog {
		if svc.Type != serviceType {
			continue
		}
		for _, ep := range svc.Endpoints {
			if ep.Interface == iface && (region == "" || ep.Region == region) {
				return ep.URL, nil
			}
		}
	}

	return "", fmt.Errorf("keystone: endpoint not found for type=%s interface=%s region=%s", serviceType, iface, region)
}
