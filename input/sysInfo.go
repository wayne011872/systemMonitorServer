package input

import (
	"errors"
	myDao "github.com/wayne011872/systemMonitorServer/dao"
)

type SysInfoInput struct {
	*myDao.SysInfo `json:",inline"`
}

func (si *SysInfoInput) Validate() error {
	if si.Ip == "" {
		return errors.New("missing ip")
	}
	return nil
}