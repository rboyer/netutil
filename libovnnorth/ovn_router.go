package libovnnorth

import "errors"

type OVN_LogicalRouter struct{}

var (
	LogicalRouter = &OVN_LogicalRouter{}
)

func (o *OVN_LogicalRouter) Add(routerName string) error {
	if routerName == "" {
		return errors.New("no LR name")
	}
	return exec("--may-exist", "lr-add", routerName)
}

func (o *OVN_LogicalRouter) Delete(routerName string) error {
	if routerName == "" {
		return nil
	}
	return exec("--if-exists", "lr-del", routerName)
}
