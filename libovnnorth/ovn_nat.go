package libovnnorth

type OVN_NAT struct{}

var (
	NAT = &OVN_NAT{}
)

func (o *OVN_NAT) AddSNAT(routerName, externalIP, logicalCIDR string) error {
	return exec(
		"--may-exist", "lr-nat-add", routerName,
		"snat",
		externalIP,
		logicalCIDR,
		// options:stateless=false
	)
}
