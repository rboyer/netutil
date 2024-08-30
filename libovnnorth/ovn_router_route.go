package libovnnorth

import "errors"

type OVN_LogicalRouterRoute struct{}

var (
	LogicalRouterRoute = &OVN_LogicalRouterRoute{}
)

// nextHop can be equal to "discard"
func (o *OVN_LogicalRouterRoute) Add(routerName, prefix, nextHop, outputLRP string) error {
	if routerName == "" {
		return errors.New("no LR name")
	}
	if prefix == "" {
		return errors.New("no prefix")
	}
	if nextHop == "" {
		return errors.New("no next hop")
	}
	return exec("--may-exist", "lr-route-add", routerName, prefix, nextHop)
}

func (o *OVN_LogicalRouterRoute) AddFancy(routerName, prefix, nextHop, outputLRP string) error {
	if routerName == "" {
		return errors.New("no LR name")
	}
	if prefix == "" {
		return errors.New("no prefix")
	}
	if nextHop == "" {
		return errors.New("no next hop")
	}
	if outputLRP == "" {
		return errors.New("no logical router port output")
	}

	return exec(
		"--", "--id=@route", "create", "logical_router_static_route",
		"ip_prefix="+prefix,
		"nexthop="+nextHop,
		"output_port="+outputLRP,
		"--", "set", "logical_router", routerName,
		"static_routes=@route",
	)
}
