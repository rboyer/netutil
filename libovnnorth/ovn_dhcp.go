package libovnnorth

import (
	"fmt"

	"github.com/rboyer/netutil/internal/util"
	"github.com/rboyer/netutil/libovs"
)

type OVN_DHCPOptions struct{}

var (
	DHCPOptions = &OVN_DHCPOptions{}
)

func (o *OVN_DHCPOptions) CreateWithExternalID(subnet string, extID string) (string, error) {
	if subnet == "" {
		return "", fmt.Errorf("DHCP subnet is empty")
	}
	if extID == "" {
		return "", fmt.Errorf("DHCP external id is empty")
	}

	uuid, err := o.GetByExternalID(extID)
	if err != nil {
		return "", err
	} else if uuid != "" {
		return uuid, nil
	}

	if err := exec("dhcp-options-create", subnet, extID); err != nil {
		return "", err
	}

	return o.GetByExternalID(extID)
}

func (o *OVN_DHCPOptions) SetOptions(dhcpUUID string, kvs ...string) error {
	if dhcpUUID == "" {
		return nil
	}

	opts, err := util.PairsToOptions("", kvs)
	if err != nil {
		return err
	} else if len(opts) == 0 {
		return nil // TODO: erase the options?
	}

	args := append([]string{
		"dhcp-options-set-options", dhcpUUID,
	}, opts...)

	return exec(args...)
}

func (o *OVN_DHCPOptions) Delete(dhcpUUID string) error {
	if dhcpUUID == "" {
		return nil
	}
	return exec("dhcp-options-del", dhcpUUID)
}

func (o *OVN_DHCPOptions) GetByExternalID(extID string) (string, error) {
	if extID == "" {
		return "", fmt.Errorf("DHCP external id is empty")
	}
	return execOutput(
		"--bare", "--columns=_uuid", "find", "dhcp_options",
		libovs.ExternalIDs+":"+extID,
	)
}
