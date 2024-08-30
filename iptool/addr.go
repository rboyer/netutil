package iptool

import (
	"encoding/json"
	"fmt"
	"net"
)

func (c *IPAddr) GetJSON(linkName string) (string, error) {
	return c.ip.execIPOutput("-details", "-json", "addr", "show", "dev", linkName)
}

func (c *IPAddr) GetIPv4(linkName string) (ipAddr, cidr string, _ error) {
	raw, err := c.GetJSON(linkName)
	if err != nil {
		return "", "", err
	}

	link, err := decodeJSONAddress(raw)
	if err != nil {
		return "", "", err
	} else if link == nil {
		return "", "", nil
	}

	var cidrs []string
	for _, addr := range link.Info {
		if addr.Family == "inet" && addr.Local != "" {
			cidr := addr.CIDR()
			cidrs = append(cidrs, cidr)
		}
	}

	if len(cidrs) == 0 {
		return "", "", nil
	} else if len(cidrs) != 1 {
		return "", "", fmt.Errorf("link %q has too many addresses %d", linkName, len(cidrs))
	}
	gotCIDR := cidrs[0]

	// ip -details -json addr show dev wlp0s20f3 | jq -r '.[].addr_info | map(select(.family == "inet") | .)[] | "\(.local)/\(.prefixlen)"'

	ip, ipnet, err := net.ParseCIDR(gotCIDR)
	if err != nil {
		return "", "", err
	}

	return ip.To4().String(), ipnet.String(), nil
}

func (c *IPAddr) ExistsOnBridge(bridge, cidr string) (bool, error) {
	raw, err := c.GetJSON(bridge)
	if err != nil {
		if isErrNotExist(err) {
			return false, nil
		}
		return false, err
	}

	link, err := decodeJSONAddress(raw)
	if err != nil {
		return false, err
	} else if link == nil {
		return false, nil
	}

	for _, addr := range link.Info {
		if addr.Family != "inet" {
			continue
		}
		found := fmt.Sprintf("%s/%d", addr.Local, addr.PrefixLen)

		if found == cidr {
			return true, nil
		}
		return false, fmt.Errorf("bridge %q has incorrect ipv4 address %q", bridge, found)
	}

	return false, nil
}

func (c *IPAddr) AddToBridge(bridge, cidr string) error {
	if ok, err := c.ExistsOnBridge(bridge, cidr); err != nil {
		return err
	} else if ok {
		return nil
	}

	return c.ip.execIP(
		"addr", "add", cidr, "dev", bridge,
	)
}

type jAddr struct {
	Info []*jAddrInfo `json:"addr_info"`
}
type jAddrInfo struct {
	Family    string `json:"family"`
	Local     string `json:"local"`
	PrefixLen int    `json:"prefixlen"`
}

func (a *jAddrInfo) CIDR() string {
	return fmt.Sprintf("%s/%d", a.Local, a.PrefixLen)
}

func decodeJSONAddresses(raw string) ([]*jAddr, error) {
	var links []*jAddr
	if err := json.Unmarshal([]byte(raw), &links); err != nil {
		return nil, err
	}
	return links, nil
}

func decodeJSONAddress(raw string) (*jAddr, error) {
	links, err := decodeJSONAddresses(raw)
	if err != nil {
		return nil, err
	}
	if len(links) == 0 {
		return nil, nil
	} else if len(links) != 1 {
		return nil, fmt.Errorf("found too many results for unary query")
	}

	return links[0], nil
}
