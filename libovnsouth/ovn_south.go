package libovnsouth

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rboyer/netutil/runner"
)

func DumpDB(w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}
	_, err := runner.Exec("ovsdb-client", []string{
		"dump", "unix:///run/ovn/ovnsb_db.sock",
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

func ChassisUUIDGetByHostname(hostname string) (string, error) {
	if hostname == "" {
		return "", fmt.Errorf("hostname is empty")
	}
	raw, err := execOutput(
		"--bare", "--columns", "name", "find", "Chassis", "hostname="+hostname,
	)
	if err != nil {
		return "", err
	}
	scan := bufio.NewScanner(strings.NewReader(raw))
	var last string
	for scan.Scan() {
		last = scan.Text()
	}
	if scan.Err() != nil {
		return "", scan.Err()
	}
	return last, nil
	// uid=$(sudo ovn-sbctl --bare --columns name find Chassis hostname=mamashark | tail -n 1)
}

func execOutput(args ...string) (string, error) {
	return runner.ExecSimpleOutput("ovn-sbctl", args)
}

func exec(args ...string) error {
	return runner.ExecSimple("ovn-sbctl", args)
}
