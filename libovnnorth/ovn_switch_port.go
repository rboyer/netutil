package libovnnorth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rboyer/netutil/internal/util"
)

type OVN_LogicalSwitchPort struct{}

var (
	LogicalSwitchPort = &OVN_LogicalSwitchPort{}
)

func (o *OVN_LogicalSwitchPort) Add(switchName, switchPortName string) error {
	if switchName == "" {
		return errors.New("no LS name")
	}
	if switchPortName == "" {
		return errors.New("no LSP name")
	}
	return exec("--may-exist", "lsp-add", switchName, switchPortName)
}

func (o *OVN_LogicalSwitchPort) AddWithAddresses(
	switchName, switchPortName string,
	mac, ipv4, dhcpUUID string,
) error {
	if err := o.Add(switchName, switchPortName); err != nil {
		return err
	}
	if mac != "" && ipv4 != "" {
		if err := o.SetAddresses2(switchPortName, mac, ipv4); err != nil {
			return err
		}
	}
	if dhcpUUID != "" {
		if err := o.SetDHCPv4Options(switchPortName, dhcpUUID); err != nil {
			return err
		}
	}
	return nil
}

func (o *OVN_LogicalSwitchPort) SetType(switchPortName string, typ string) error {
	if switchPortName == "" {
		return nil
	}
	if typ == "" {
		return errors.New("no switch port type name")
	}
	return exec("lsp-set-type", switchPortName, typ)
}

// unary address field like "router"
func (o *OVN_LogicalSwitchPort) SetAddresses1(switchPortName, addr string) error {
	if switchPortName == "" {
		return nil
	}
	return exec("lsp-set-addresses", switchPortName, addr)
}

// binary address field like "${MAC} dynamic"
func (o *OVN_LogicalSwitchPort) SetAddresses2(switchPortName, mac, addr string) error {
	if switchPortName == "" {
		return nil
	}
	if mac == "" {
		return fmt.Errorf("mac is required")
	}
	if addr == "" {
		return fmt.Errorf("addr is required")
	}
	return exec("lsp-set-addresses", switchPortName, mac+" "+addr)
}

func (o *OVN_LogicalSwitchPort) SetDHCPv4Options(switchPortName, optionsUUID string) error {
	if switchPortName == "" {
		return nil
	}
	if optionsUUID == "" {
		return nil
	}
	return exec("lsp-set-dhcpv4-options", switchPortName, optionsUUID)
}

func (o *OVN_LogicalSwitchPort) SetOptions(switchPortName string, kvs ...string) error {
	if switchPortName == "" {
		return nil
	}

	opts, err := util.PairsToOptions("", kvs)
	if err != nil {
		return err
	} else if len(opts) == 0 {
		return nil // TODO: erase the options?
	}

	args := append([]string{"lsp-set-options", switchPortName}, opts...)

	return exec(args...)
}

func (o *OVN_LogicalSwitchPort) DebugListAddresses() (string, error) {
	// udo ovn-nbctl get logical_switch_port 4d90c2a3-a6e5-4424-882b-fb513b196bbb dynamic_addresses
	return execOutput(
		"--columns", "dynamic_addresses", "list", "logical_switch_port",
	)
}

func (o *OVN_LogicalSwitchPort) GetAddresses(switchPortName string) (mac, ip string, _ error) {
	if switchPortName == "" {
		return "", "", nil
	}

	raw, err := execOutput(
		"get", "logical_switch_port", switchPortName, "dynamic_addresses",
	)
	if err != nil {
		return "", "", err
	}

	macIP := strings.Trim(raw, `"`)

	mac, ip, found := strings.Cut(macIP, " ")
	if !found {
		return "", "", fmt.Errorf("lsp %q lacks IP", switchPortName)
	}

	return mac, ip, nil
}

func (o *OVN_LogicalSwitchPort) Delete(switchPortName string) error {
	if switchPortName == "" {
		return nil
	}
	return exec("--if-exists", "lsp-del", switchPortName)
}
