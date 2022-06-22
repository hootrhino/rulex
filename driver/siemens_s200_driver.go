package driver

import (
	"github.com/i4de/rulex/typex"

	"github.com/robinson/gos7"
)

/*
*
* 西门子S200驱动
*
 */
type siemens_s200_driver struct {
	state      typex.DriverState
	s7client   gos7.Client
	device     *typex.Device
	RuleEngine typex.RuleX
}

func NewS200Driver(d *typex.Device, e typex.RuleX, s7client gos7.Client) typex.XExternalDriver {
	return &siemens_s200_driver{
		state:      typex.DRIVER_STOP,
		device:     d,
		RuleEngine: e,
		s7client:   s7client,
	}
}

func (s200 *siemens_s200_driver) Test() error {
	return nil
}

func (s200 *siemens_s200_driver) Init(_ map[string]string) error {
	return nil
}

func (s200 *siemens_s200_driver) Work() error {
	return nil
}

func (s200 *siemens_s200_driver) State() typex.DriverState {
	return typex.DRIVER_RUNNING
}

func (s200 *siemens_s200_driver) Read(_ []byte) (int, error) {
	return 0, nil
}

func (s200 *siemens_s200_driver) Write(_ []byte) (int, error) {
	return 0, nil
}

func (s200 *siemens_s200_driver) DriverDetail() typex.DriverDetail {
	return typex.DriverDetail{
		Name:        "SIEMENS_S200",
		Type:        "TCP",
		Description: "SIEMENS S200 系列 PLC 驱动",
	}
}

func (s200 *siemens_s200_driver) Stop() error {
	return nil
}
