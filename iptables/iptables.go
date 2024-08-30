package iptables

import (
	"bufio"
	"log"
	"strings"

	"github.com/rboyer/netutil/runner"
)

func InsertChain(tableName, chain, baseChain string) error {
	found, err := FindChain(tableName, chain)
	if err != nil {
		return err
	}
	if !found {
		// iptables -t $tableName -N $chain
		if err := iptablesCommand("-t", tableName, "-N", chain); err != nil {
			return err
		}
	}

	return InsertRule(tableName, baseChain, "-j", chain)
}

func DeleteChain(tableName, chain, baseChain string) error {
	if err := DeleteRule(tableName, baseChain, "-j", chain); err != nil {
		return err
	}

	found, err := FindChain(tableName, chain)
	if err != nil {
		return err
	} else if !found {
		return nil
	}

	// iptables -t $tableName -X $chain
	return iptablesCommand("-t", tableName, "-X", chain)
}

func FindChain(tableName, chain string) (bool, error) {
	chains, err := iptablesCommandOutput("-t", tableName, "-L", "-n")
	if err != nil {
		return false, err
	}
	scan := bufio.NewScanner(strings.NewReader(chains))
	for scan.Scan() {
		if strings.HasPrefix(scan.Text(), "Chain "+chain+" (") {
			log.Printf("chain [iptables -t %s -N %s] found", tableName, chain)
			return true, nil
		}
	}
	if scan.Err() != nil {
		return false, scan.Err()
	}
	log.Printf("chain [iptables -t %s -N %s] not found", tableName, chain)
	return false, nil
}

func InsertRule(tableName string, suffix ...string) error {
	found, err := FindRule(tableName, suffix...)
	if err != nil {
		return err
	} else if found {
		return nil
	}

	args := append([]string{"-t", tableName, "-I"}, suffix...)
	return iptablesCommand(args...)
}

func DeleteRule(tableName string, suffix ...string) error {
	found, err := FindRule(tableName, suffix...)
	if err != nil {
		return err
	} else if !found {
		return nil
	}
	args := append([]string{"-t", tableName, "-D"}, suffix...)
	return iptablesCommand(args...)
}

func FindRule(tableName string, suffix ...string) (bool, error) {
	args := append([]string{"-t", tableName, "-C"}, suffix...)
	code, err := iptablesCommandExitCode(args...)
	if err != nil {
		return false, err
	} else if code == 0 {
		log.Printf("rule [iptables -t %s -I %s] found", tableName, strings.Join(suffix, " "))
		return true, nil
	}
	log.Printf("rule [iptables -t %s -I %s] not found", tableName, strings.Join(suffix, " "))
	return false, nil
}

func iptablesCommandExitCode(args ...string) (int, error) {
	return runner.ExecSimpleExitCode("iptables", args)
}

func iptablesCommand(args ...string) error {
	return runner.ExecSimple("iptables", args)
}

func iptablesCommandOutput(args ...string) (string, error) {
	return runner.ExecSimpleOutput("iptables", args)
}
