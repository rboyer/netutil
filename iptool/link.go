package iptool

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

func (c *IPLink) SetMaster(linkName, masterName string) error {
	// sudo ip set master br-fake dev br-fake-a
	return c.ip.execIP(
		"link", "set", "master", masterName, "dev", linkName,
	)
}

func (c *IPLink) SetUp(linkName string) error {
	return c.ip.execIP(
		"link", "set", "dev", linkName, "up",
	)
}

func (c *IPLink) SetMAC(linkName, mac string) error {
	return c.ip.execIP(
		"link", "set", linkName, "address", mac,
	)
}

func (c *IPLink) SetNetNS(linkName, netns string) error {
	if c.ip.namespace != "" {
		panic("cannot use nested namesapces with SetNetNS")
	}
	return c.ip.execIP(
		"link", "set", linkName, "netns", netns,
	)
}

func (c *IPLink) VethPairExists(name1, name2 string) (bool, error) {
	ok1, err := c.Exists(name1, "veth")
	if err != nil {
		return false, err
	}
	ok2, err := c.Exists(name2, "veth")
	if err != nil {
		return false, err
	}

	return ok1 && ok2, nil
}

func (c *IPLink) CreateVethPair(name1, name2 string) error {
	if ok, err := c.VethPairExists(name1, name2); err != nil {
		return err
	} else if ok {
		return nil
	}
	// sudo ip link add br-provider-a type veth peer name br-provider-b
	return c.ip.execIP(
		"link", "add", name1, "type", "veth", "peer", "name", name2,
	)
}

func isErrNotExist(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "does not exist")
}

func (c *IPLink) BridgeExists(linkName string) (bool, error) {
	return c.Exists(linkName, "bridge")
}

func (c *IPLink) Exists(linkName, kind string) (bool, error) {
	raw, err := c.GetJSON(linkName)
	if err != nil {
		if isErrNotExist(err) {
			return false, nil
		}
		return false, err
	}

	link, err := decodeJSONLink(raw)
	if err != nil {
		return false, err
	}

	var found string
	if link.Info != nil {
		found = link.Info.Kind
	}
	if found == kind {
		return true, nil
	}
	return false, fmt.Errorf("link %q exists but is of kind %q", linkName, found)
}

func (c *IPLink) BridgeCreate(linkName string) error {
	if ok, err := c.BridgeExists(linkName); err != nil {
		return err
	} else if ok {
		return nil
	}

	return c.ip.execIP(
		"link", "add", "name", linkName, "type", "bridge",
	)
}

func (c *IPLink) Delete(linkName string) error {
	return c.ip.execIP("link", "del", linkName)
}

// This will attempt to get the public link interface name
func (c *IPLink) GetPublic() (string, error) {
	// https://serverfault.com/questions/1019363/using-ip-address-show-type-to-display-physical-network-interface
	raw, err := c.ListJSON()
	if err != nil {
		return "", err
	}

	links, err := decodeJSONLinks(raw)
	if err != nil {
		return "", err
	}

	var public []string
OUTER:
	for _, link := range links {
		if link.Info != nil || link.Type == "loopback" {
			continue
		}

		for _, f := range link.Flags {
			if f == "NO-CARRIER" {
				continue OUTER
			}
		}
		public = append(public, link.Name)
	}
	sort.Strings(public)

	if len(public) == 0 {
		return "", fmt.Errorf("no public link detected")
	} else if len(public) != 1 {
		return "", fmt.Errorf("multiple public links detected: %v", public)
	}

	return public[0], nil
}

func (c *IPLink) GetJSON(linkName string) (string, error) {
	return c.ip.execIPOutput(
		"-details", "-json", "link", "show", "dev", linkName,
	)
}

func (c *IPLink) ListJSON() (string, error) {
	return c.ip.execIPOutput(
		"-details", "-json", "link",
	)
}

type jLink struct {
	Name  string     `json:"ifname"`
	Info  *jLinkInfo `json:"linkinfo"`
	Type  string     `json:"link_type"`
	Flags []string   `json:"flags"`
}
type jLinkInfo struct {
	Kind string `json:"info_kind"`
}

func decodeJSONLinks(raw string) ([]*jLink, error) {
	var links []*jLink
	if err := json.Unmarshal([]byte(raw), &links); err != nil {
		return nil, err
	}
	return links, nil
}

func decodeJSONLink(raw string) (*jLink, error) {
	links, err := decodeJSONLinks(raw)
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
