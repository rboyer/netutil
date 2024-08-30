package iptool

import (
	"github.com/rboyer/netutil/runner"
)

type IPIP struct {
	namespace string

	NetNS IPNetNS
	Link  IPLink
	Addr  IPAddr
	Route IPRoute
}

var IP = create("")

func Namespaced(netns string) *IPIP {
	if netns == "" {
		return IP
	}
	return create(netns)
}

func create(netns string) *IPIP {
	ip := &IPIP{namespace: netns}
	ip.setPointers()
	if netns != "" {
		ip.NetNS.disabled = true
	}
	return ip
}

func (i *IPIP) setPointers() {
	i.NetNS.ip = i
	i.Link.ip = i
	i.Addr.ip = i
	i.Route.ip = i
}

type IPNetNS struct {
	ip       *IPIP
	disabled bool
}

type IPLink struct {
	ip *IPIP
}

type IPAddr struct {
	ip *IPIP
}

type IPRoute struct {
	ip *IPIP
}

// ================= Hybrid =================

func (i *IPIP) CreateLoopback(setUP bool) error {
	// SIM: create loopback
	if err := i.execIP("addr", "add", "127.0.0.1/8", "dev", "lo"); err != nil {
		return err
	}

	if setUP {
		// SIM: Bring up loopback
		if err := i.execIP("link", "set", "dev", "lo", "up"); err != nil {
			return err
		}
	}
	return nil
}

// ================= Util =================

func (i *IPIP) Exec(args ...string) error {
	if i.namespace == "" {
		panic("Exec only works in a namespace")
	}
	return runner.ExecSimple("ip",
		append([]string{"netns", "exec", i.namespace}, args...),
	)
}

func (i *IPIP) execIP(args ...string) error {
	if i.namespace != "" {
		args = append([]string{"-n", i.namespace}, args...)
	}
	return runner.ExecSimple("ip", args)
}

func (i *IPIP) execIPOutput(args ...string) (string, error) {
	if i.namespace != "" {
		args = append([]string{"-n", i.namespace}, args...)
	}
	return runner.ExecSimpleOutput("ip", args)
}
