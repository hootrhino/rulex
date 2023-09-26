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
type modBusRtuDriver struct {
	state      typex.DriverState
	handler    *modbus.RTUClientHandler
	client     modbus.Client
	RuleEngine typex.RuleX
	Registers  []common.RegisterRW
	device     *typex.Device
	frequency  int64
}

func NewModBusRtuDriver(
	d *typex.Device,
	e typex.RuleX,
	Registers []common.RegisterRW,
	handler *modbus.RTUClientHandler,
	client modbus.Client, frequency int64) typex.XExternalDriver {
	return &modBusRtuDriver{
		state:      typex.DRIVER_UP,
		device:     d,
		RuleEngine: e,
		client:     client,
		handler:    handler,
		Registers:  Registers,
		frequency:  frequency,
	}

}
func (d *modBusRtuDriver) Test() error {
	return nil
}

func (d *modBusRtuDriver) Init(map[string]string) error {
	return nil
}

func (d *modBusRtuDriver) Work() error {
	return nil
}

func (d *modBusRtuDriver) State() typex.DriverState {
	return d.state
}

func (d *modBusRtuDriver) Read(cmd []byte, data []byte) (int, error) {
	var err error
	var results []byte
	dataMap := map[string]common.RegisterRW{}
	count := len(d.Registers)
	for _, r := range d.Registers {
		d.handler.SlaveId = r.SlaverId
		if r.Function == common.READ_COIL {
			results, err = d.client.ReadCoils(r.Address, r.Quantity)
			if err != nil {
				count--
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
			results, err = d.client.ReadDiscreteInputs(r.Address, r.Quantity)
			if err != nil {
				count--
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
			results, err = d.client.ReadHoldingRegisters(r.Address, r.Quantity)
			if err != nil {
				count--
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
			results, err = d.client.ReadInputRegisters(r.Address, r.Quantity)
			if err != nil {
				count--
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
		time.Sleep(time.Duration(d.frequency) * time.Millisecond)
	}
	bytes, _ := json.Marshal(dataMap)
	copy(data, bytes)
	// 只要有部分成功，哪怕有一个设备出故障也认为是正常的，上层可以根据Value来判断
	ll := len(d.Registers)
	if ll > 0 && count > 0 {
		return len(bytes), nil
	}
	return len(bytes), err

}

/*
*
* 支持Modbus写入
*
 */
func (d *modBusRtuDriver) Write(_ []byte, data []byte) (int, error) {
	dataMap := []common.RegisterW{}
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return 0, err
	}
	for _, r := range dataMap {
		// 5
		if r.Function == common.WRITE_SINGLE_COIL {
			_, err := d.client.WriteSingleCoil(r.Address, binary.BigEndian.Uint16(r.Values))
			if err != nil {
				return 0, err
			}
		}
		// 15
		if r.Function == common.WRITE_MULTIPLE_COILS {
			_, err := d.client.WriteMultipleCoils(r.Address, uint16(len(r.Values)), r.Values)
			if err != nil {
				return 0, err
			}
		}
		// 6
		if r.Function == common.WRITE_SINGLE_HOLDING_REGISTER {
			_, err := d.client.WriteSingleRegister(r.Address, binary.BigEndian.Uint16(r.Values))
			if err != nil {
				return 0, err
			}
		}
		// 16
		if r.Function == common.WRITE_MULTIPLE_HOLDING_REGISTERS {
			_, err := d.client.WriteMultipleRegisters(r.Address, uint16(len(r.Values)), r.Values)
			if err != nil {
				return 0, err
			}
		}
	}
	return 0, nil
}

func (d *modBusRtuDriver) DriverDetail() typex.DriverDetail {
	return typex.DriverDetail{
		Name:        "Generic ModBus RTU Driver",
		Type:        "UART",
		Description: "Generic ModBus RTU Driver",
	}
}

func (d *modBusRtuDriver) Stop() error {
	return d.handler.Close()
}
