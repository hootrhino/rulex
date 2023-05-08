package driver

import (
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	modbus "github.com/wwhai/gomodbus"
)

/*
*
* Modbus RTU
*
 */
type modBusTCPDriver struct {
	state      typex.DriverState
	handler    *modbus.TCPClientHandler
	client     modbus.Client
	RuleEngine typex.RuleX
	Registers  []common.RegisterRW
	device     *typex.Device
}

func NewModBusTCPDriver(
	d *typex.Device,
	e typex.RuleX,
	Registers []common.RegisterRW,
	handler *modbus.TCPClientHandler,
	client modbus.Client) typex.XExternalDriver {
	return &modBusTCPDriver{
		state:      typex.DRIVER_UP,
		device:     d,
		RuleEngine: e,
		client:     client,
		handler:    handler,
		Registers:  Registers,
	}

}
func (d *modBusTCPDriver) Test() error {
	return nil
}

func (d *modBusTCPDriver) Init(map[string]string) error {
	return nil
}

func (d *modBusTCPDriver) Work() error {
	return nil
}

func (d *modBusTCPDriver) State() typex.DriverState {
	return d.state
}

func (d *modBusTCPDriver) Read(cmd []byte, data []byte) (int, error) {
	dataMap := map[string]common.RegisterRW{}
	for _, r := range d.Registers {
		d.handler.SlaveId = r.SlaverId
		if r.Function == common.READ_COIL {
			results, err := d.client.ReadCoils(r.Address, r.Quantity)
			if err != nil {
				glogger.GLogger.Error(err)
			}
			value := common.RegisterRW{
				Tag:      r.Tag,
				Function: r.Function,
				SlaverId: r.SlaverId,
				Address:  r.Address,
				Quantity: r.Quantity,
				Value:    covertEmptyHex(results),
			}
			dataMap[r.Tag] = value
		}
		if r.Function == common.READ_DISCRETE_INPUT {
			results, err := d.client.ReadDiscreteInputs(r.Address, r.Quantity)
			if err != nil {
				glogger.GLogger.Error(err)
			}
			value := common.RegisterRW{
				Tag:      r.Tag,
				Function: r.Function,
				SlaverId: r.SlaverId,
				Address:  r.Address,
				Quantity: r.Quantity,
				Value:    covertEmptyHex(results),
			}
			dataMap[r.Tag] = value

		}
		if r.Function == common.READ_HOLDING_REGISTERS {
			results, err := d.client.ReadHoldingRegisters(r.Address, r.Quantity)
			if err != nil {
				glogger.GLogger.Error(err)
			}
			value := common.RegisterRW{
				Tag:      r.Tag,
				Function: r.Function,
				SlaverId: r.SlaverId,
				Address:  r.Address,
				Quantity: r.Quantity,
				Value:    covertEmptyHex(results),
			}
			dataMap[r.Tag] = value
		}
		if r.Function == common.READ_INPUT_REGISTERS {
			results, err := d.client.ReadInputRegisters(r.Address, r.Quantity)
			if err != nil {
				glogger.GLogger.Error(err)
			}
			value := common.RegisterRW{
				Tag:      r.Tag,
				Function: r.Function,
				SlaverId: r.SlaverId,
				Address:  r.Address,
				Quantity: r.Quantity,
				Value:    covertEmptyHex(results),
			}
			dataMap[r.Tag] = value
		}
		// 设置一个间隔时间防止低级CPU黏包等
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
	bytes, _ := json.Marshal(dataMap)
	copy(data, bytes)
	return len(bytes), nil

}

func (d *modBusTCPDriver) Write(cmd []byte, data []byte) (int, error) {
	dataMap := []common.RegisterRW{}
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return 0, err
	}
	for _, r := range d.Registers {
		if r.Function == common.WRITE_SINGLE_COIL {
			_, err := d.client.WriteSingleCoil(r.Address, binary.BigEndian.Uint16([]byte(r.Value)[0:2]))
			if err != nil {
				return 0, err
			}
		}
		if r.Function == common.WRITE_MULTIPLE_COILS {
			_, err := d.client.WriteMultipleCoils(r.Address, r.Quantity, []byte(r.Value))
			if err != nil {
				return 0, err
			}
		}
		if r.Function == common.WRITE_SINGLE_HOLDING_REGISTER {
			_, err := d.client.WriteSingleRegister(r.Address, binary.BigEndian.Uint16([]byte(r.Value)[0:2]))
			if err != nil {
				return 0, err
			}
		}
		if r.Function == common.WRITE_MULTIPLE_HOLDING_REGISTERS {
			_, err := d.client.WriteMultipleRegisters(r.Address, r.Quantity, []byte(r.Value))
			if err != nil {
				return 0, err
			}
		}
	}
	return 0, nil
}

func (d *modBusTCPDriver) DriverDetail() typex.DriverDetail {
	return typex.DriverDetail{
		Name:        "Generic ModBus TCP Driver",
		Type:        "TCP",
		Description: "Generic ModBus TCP Driver",
	}
}

func (d *modBusTCPDriver) Stop() error {
	d.handler.Close()
	d = nil
	return nil
}
