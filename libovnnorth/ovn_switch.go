package libovnnorth

import (
	"errors"

	"github.com/rboyer/netutil/internal/util"
)

type OVN_LogicalSwitch struct{}

var (
	LogicalSwitch = &OVN_LogicalSwitch{}
)

func (o *OVN_LogicalSwitch) Add(switchName string) error {
	if switchName == "" {
		return errors.New("no LS name")
	}
	return exec("--may-exist", "ls-add", switchName)
}

func (o *OVN_LogicalSwitch) SetOtherConfig(switchName string, kvs ...string) error {
	if switchName == "" {
		return nil
	}
	opts, err := util.PairsToOptions("other_config", kvs)
	if err != nil {
		return err
	} else if len(opts) == 0 {
		return nil // TODO: erase the options?
	}

	args := append([]string{
		"set", "logical_switch", switchName,
	}, opts...)

	return exec(args...)
}

func (o *OVN_LogicalSwitch) Delete(switchName string) error {
	if switchName == "" {
		return nil
	}
	return exec("--if-exists", "ls-del", switchName)
}
