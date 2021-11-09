package driver

import "rulex/typex"

type DemoDriver struct {
}

func (d *DemoDriver) Test() error {
	panic("not implemented") // TODO: Implement
}

func (d *DemoDriver) Init() error {
	panic("not implemented") // TODO: Implement
}

func (d *DemoDriver) Work() error {
	panic("not implemented") // TODO: Implement
}

func (d *DemoDriver) State() typex.DriverState {
	panic("not implemented") // TODO: Implement
}

func (d *DemoDriver) SetState(_ typex.DriverState) {
	panic("not implemented") // TODO: Implement
}

//---------------------------------------------------
// 读写接口是给LUA标准库用的, 驱动只管实现读写逻辑即可
//---------------------------------------------------
func (d *DemoDriver) Read(_ []byte) (int, error) {
	panic("not implemented") // TODO: Implement
}

func (d *DemoDriver) Write(_ []byte) (int, error) {
	panic("not implemented") // TODO: Implement
}

//---------------------------------------------------
func (d *DemoDriver) DriverDetail() *typex.DriverDetail {
	panic("not implemented") // TODO: Implement
}

func (d *DemoDriver) Stop() error {
	panic("not implemented") // TODO: Implement
}
