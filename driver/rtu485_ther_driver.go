// 485 温湿度传感器驱动案例
// 这是个很简单的485温湿度传感器驱动, 目的是为了演示厂商如何实现自己的设备底层驱动
// 本驱动完成于2022年4月28日, 温湿度传感器资料请移步文档
package driver

import "rulex/typex"

type rtu485_THer_Driver struct {
}

func NewRtu485_THer_Driver() typex.XExternalDriver {

	return &rtu485_THer_Driver{}
}
func (rtu485 *rtu485_THer_Driver) Test() error {
	return nil
}

func (rtu485 *rtu485_THer_Driver) Init() error {
	return nil
}

func (rtu485 *rtu485_THer_Driver) Work() error {
	return nil
}

func (rtu485 *rtu485_THer_Driver) State() typex.DriverState {
	return typex.RUNNING
}

func (rtu485 *rtu485_THer_Driver) SetState(_ typex.DriverState) {

}

//---------------------------------------------------
// 读写接口是给LUA标准库用的, 驱动只管实现读写逻辑即可
//---------------------------------------------------
func (rtu485 *rtu485_THer_Driver) Read(_ []byte) (int, error) {
	return 0, nil
}

func (rtu485 *rtu485_THer_Driver) Write(_ []byte) (int, error) {
	return 0, nil

}

//---------------------------------------------------
func (rtu485 *rtu485_THer_Driver) DriverDetail() *typex.DriverDetail {
	return nil
}

func (rtu485 *rtu485_THer_Driver) Stop() error {
	return nil
}
