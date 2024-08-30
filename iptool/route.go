package iptool

import (
	"encoding/json"
	"fmt"
)

func (c *IPRoute) ListJSON() (string, error) {
	return c.ip.execIPOutput("-details", "-json", "route")
}

type jRoute struct {
	Dest string `json:"dst"`
	Dev  string `json:"dev"`
}

func decodeJSONRoutes(raw string) ([]*jRoute, error) {
	var routes []*jRoute
	if err := json.Unmarshal([]byte(raw), &routes); err != nil {
		return nil, err
	}
	return routes, nil
}

func (c *IPRoute) GetDefaultRoute() (string, error) {
	raw, err := c.ListJSON()
	if err != nil {
		return "", err
	}

	routes, err := decodeJSONRoutes(raw)
	if err != nil {
		return "", err
	}

	for _, route := range routes {
		if route.Dest == "default" {
			return route.Dev, nil
		}
	}
	return "", fmt.Errorf("no public link detected")
}
