package libovnnorth

import (
	"errors"
	"strconv"

	"github.com/rboyer/netutil/internal/util"
)

type OVN_LogicalRouterPort struct{}

var (
	LogicalRouterPort = &OVN_LogicalRouterPort{}
)

func (o *OVN_LogicalRouterPort) Add(routerName, routerPortName, mac, cidr string) error {
	if routerName == "" {
		return errors.New("no LR name")
	}
	if routerPortName == "" {
		return errors.New("no LRP name")
	}
	if mac == "" {
		return errors.New("no LRP MAC")
	}
	if cidr == "" {
		return errors.New("no LRP CIDR")
	}

	return exec(
		"--may-exist", "lrp-add", routerName, routerPortName, mac, cidr,
	)
}

func (o *OVN_LogicalRouterPort) SetGatewayChassis(routerPortName string, chassisUUID string, weight uint) error {
	if routerPortName == "" {
		return nil
	}
	if chassisUUID == "" {
		return nil
	}
	weightStr := strconv.FormatUint(uint64(weight), 10)
	return exec("lrp-set-gateway-chassis", routerPortName, chassisUUID, weightStr)
}

func (o *OVN_LogicalRouterPort) SetOptions(routerPortName string, kvs ...string) error {
	if routerPortName == "" {
		return nil
	}
	opts, err := util.PairsToOptions("options", kvs)
	if err != nil {
		return err
	} else if len(opts) == 0 {
		return nil // TODO: erase the options?
	}

	args := append([]string{
		"set", "Logical_Router_Port", routerPortName,
	}, opts...)

	return exec(args...)
}
