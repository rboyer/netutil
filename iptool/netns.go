package iptool

import (
	"bufio"
	"strings"

	"github.com/rboyer/netutil/runner"
)

func (c *IPNetNS) Exists(netns string) (bool, error) {
	if c.disabled {
		panic("cannot execute nested namespace commands")
	}
	out, err := runner.ExecSimpleOutput("ip", []string{"netns", "list"})
	if err != nil {
		return false, err
	}

	scan := bufio.NewScanner(strings.NewReader(out))
	for scan.Scan() {
		if netns == strings.TrimSpace(scan.Text()) {
			return true, nil
		}
	}
	if scan.Err() != nil {
		return false, scan.Err()
	}

	return false, nil
}

func (c *IPNetNS) Create(netns string) error {
	if c.disabled {
		panic("cannot execute nested namespace commands")
	}
	// Create the netns
	exist, err := c.Exists(netns)
	if err != nil {
		return err
	} else if exist {
		return nil
	}
	return c.ip.execIP("netns", "add", netns)
}

func (c *IPNetNS) Delete(netns string) error {
	if c.disabled {
		panic("cannot execute nested namespace commands")
	}
	// return IPCommand("-all", "netns", "delete", netns)
	return c.ip.execIP("netns", "delete", netns)
}
