package libovnnorth

import (
	"fmt"
	"io"
	"os"

	"github.com/rboyer/netutil/runner"
)

func DumpDB(w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}
	_, err := runner.Exec("ovsdb-client", []string{
		"dump", "unix:///run/ovn/ovnnb_db.sock",
	}, w, os.Stderr, os.Stdin, "")
	return err
}

func ConnectionSet(conn string) error {
	if conn == "" {
		return fmt.Errorf("conn is required")
	}
	return exec("set-connection", conn)
}

func ConnectionDelete() error {
	return exec("del-connection")
}

func ConnectLogicalRouterToLogicalSwitch(
	router, switchName, routerMAC, routerCIDR string,
) (lrp, lsp string, _ error) {
	var (
		routerPort = router + "-" + switchName
		switchPort = switchName + "-" + router
	)

	// create logical router ports
	if err := LogicalRouterPort.Add(router, routerPort, routerMAC, routerCIDR); err != nil {
		return "", "", err
	}
	if err := LogicalSwitchPort.Add(switchName, switchPort); err != nil {
		return "", "", err
	}
	if err := LogicalSwitchPort.SetType(switchPort, "router"); err != nil {
		return "", "", err
	}
	if err := LogicalSwitchPort.SetAddresses1(switchPort, "router"); err != nil {
		return "", "", err
	}
	if err := LogicalSwitchPort.SetOptions(switchPort,
		"router-port", routerPort,
		"nat-addresses", "router",
	); err != nil {
		return "", "", err
	}
	return routerPort, switchPort, nil
}

func execOutput(args ...string) (string, error) {
	return runner.ExecSimpleOutput("ovn-nbctl", args)
}

func exec(args ...string) error {
	return runner.ExecSimple("ovn-nbctl", args)
}
