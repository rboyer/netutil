package libovs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rboyer/netutil/internal/util"
	"github.com/rboyer/netutil/runner"
)

const (
	TableOVS       = "open_vswitch" // also "open"
	TableOVSBridge = "bridge"

	ExternalIDs = "external_ids"
)

func DumpDB(w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}
	_, err := runner.Exec("ovsdb-client", []string{
		"dump",
	}, w, os.Stderr, os.Stdin, "")
	return err
}

func BridgeHasPort(bridge, link string) (bool, error) {
	ports, err := BridgeListPorts(bridge)
	if err != nil {
		return false, err
	}
	for _, port := range ports {
		if port == link {
			return true, nil
		}
	}
	return false, nil
}

func BridgeListPorts(bridge string) ([]string, error) {
	ports, err := execOutput("list-ports", bridge)
	if err != nil {
		return nil, err
	}

	var found []string

	scan := bufio.NewScanner(strings.NewReader(ports))
	for scan.Scan() {
		found = append(found, strings.TrimSpace(scan.Text()))
	}
	if scan.Err() != nil {
		return nil, scan.Err()
	}

	return found, nil
	// ovs-vsctl list-ports "$bridge" | grep -q "$link"
}

func SetExternalIDs(kvs ...string) error {
	opts, err := util.PairsToOptions(ExternalIDs, kvs)
	if err != nil {
		return err
	} else if len(opts) == 0 {
		return fmt.Errorf("must explicitly delete ids")
	}

	args := append([]string{
		"set", TableOVS, ".",
	}, opts...)

	return exec(args...)
}

func DeleteExternalIDs() error {
	return exec("set", TableOVS, ".", ExternalIDs+"={}")
}

func BridgeCreate(bridge string) error {
	if bridge == "" {
		return fmt.Errorf("bridge is empty")
	}
	return exec("--may-exist", "add-br", bridge)
}

func BridgeSet(bridge string, kvs ...string) error {
	if bridge == "" {
		return fmt.Errorf("bridge is empty")
	}

	opts, err := util.PairsToOptions("", kvs)
	if err != nil {
		return err
	} else if len(opts) == 0 {
		return fmt.Errorf("must explicitly delete options")
	}

	args := append([]string{
		"set", TableOVSBridge, bridge,
	}, opts...)

	return exec(args...)
}

func BridgeDelete(bridge string) error {
	if bridge == "" {
		return nil
	}
	return exec("del-br", bridge)
}

func BridgeAddPort(bridge, port string, extraCommand ...string) error {
	if bridge == "" {
		return fmt.Errorf("bridge is empty")
	}
	if port == "" {
		return fmt.Errorf("port is empty")
	}

	args := []string{"--may-exist", "add-port", bridge, port}
	if len(extraCommand) > 0 {
		args = append(args, "--")
		args = append(args, extraCommand...)
	}

	return exec(args...)
}

func BridgeDeletePort(bridge, port string) error {
	if bridge == "" {
		return nil
	}
	if port == "" {
		return nil
	}
	return exec("--if-exists", "del-port", bridge, port)
}

func InterfaceGetID(netdev string) (string, error) {
	// iface_name="$(sudo ovs-vsctl get interface "$netdev" external_ids:iface-id | jq -r .)"
	ifaceNameQuoted, err := execOutput(
		"get", "interface", netdev, ExternalIDs+":iface-id",
	)
	if err != nil {
		return "", err
	}
	ifaceName := strings.Trim(ifaceNameQuoted, `"`)
	return ifaceName, nil
}

func InterfaceGetMAC(netdev string) (string, error) {
	// mac="$(virsh dumpxml "${vm_name}" | xmlstarlet sel -B -I -t -v '/domain/devices/interface[@type="bridge"]/mac/@address')"
	// mac="$(sudo ovs-vsctl get interface "$netdev" external_ids:attached-mac | jq -r .)"
	macQuoted, err := execOutput(
		"get", "interface", netdev, ExternalIDs+":attached-mac",
	)
	if err != nil {
		return "", err
	}
	mac := strings.Trim(macQuoted, `"`)
	return mac, nil
}

func execOutput(args ...string) (string, error) {
	return runner.ExecSimpleOutput("ovs-vsctl", args)
}

func exec(args ...string) error {
	return runner.ExecSimple("ovs-vsctl", args)
}
