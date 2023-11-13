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
	RegisterW := common.RegisterW{}
	if err := json.Unmarshal(data, &RegisterW); err != nil {
		return 0, err
	}
	dataMap := [1]common.RegisterW{RegisterW}
	for _, r := range dataMap {
		d.handler.SlaveId = r.SlaverId
		// 5
		if r.Function == common.WRITE_SINGLE_COIL {
			if len(r.Values) > 0 {
				if r.Values[0] == 0 {
					_, err := d.client.WriteSingleCoil(r.Address,
						binary.BigEndian.Uint16([]byte{0x00, 0x00}))
					if err != nil {
						return 0, err
					}
				}
				if r.Values[0] == 1 {
					_, err := d.client.WriteSingleCoil(r.Address,
						binary.BigEndian.Uint16([]byte{0xFF, 0x00}))
					if err != nil {
						return 0, err
					}
				}

			}

		}
		// 15
		if r.Function == common.WRITE_MULTIPLE_COILS {
			_, err := d.client.WriteMultipleCoils(r.Address, r.Quantity, r.Values)
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

			_, err := d.client.WriteMultipleRegisters(r.Address,
				uint16(len(r.Values))/2, maybePrependZero(r.Values))
			if err != nil {
				return 0, err
			}
		}
	}
	return 0, nil
}
func maybePrependZero(slice []byte) []byte {
	if len(slice)%2 != 0 {
		slice = append([]byte{0}, slice...)
	}
	return slice
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
