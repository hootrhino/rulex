// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package device

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	golog "log"

	"time"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/component/hwportmanager"
	modbuscache "github.com/hootrhino/rulex/component/intercache/modbus"
	"github.com/hootrhino/rulex/component/interdb"
	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	modbus "github.com/wwhai/gomodbus"
)

// 这是个通用Modbus采集器, 主要用来在通用场景下采集数据，因此需要配合规则引擎来使用
//
// Modbus 采集到的数据如下, LUA 脚本可做解析, 示例脚本可参照 generic_modbus_parse.lua
//
//	{
//	    "d1":{
//	        "tag":"d1",
//	        "function":3,
//	        "slaverId":1,
//	        "address":0,
//	        "quantity":2,
//	        "value":"..."
//	    },
//	    "d2":{
//	        "tag":"d2",
//	        "function":3,
//	        "slaverId":2,
//	        "address":0,
//	        "quantity":2,
//	        "value":"..."
//	    }
//	}
type _GMODCommonConfig struct {
	Mode        string `json:"mode"`
	AutoRequest *bool  `json:"autoRequest"`
}
type _GMODConfig struct {
	CommonConfig _GMODCommonConfig `json:"commonConfig" validate:"required"`
	PortUuid     string            `json:"portUuid"`
	HostConfig   common.HostConfig `json:"hostConfig"`
}

/*
*
* 点位表
*
 */
type ModbusPoint struct {
	UUID      string `json:"uuid,omitempty"` // 当UUID为空时新建
	Tag       string `json:"tag"`
	Alias     string `json:"alias"`
	Function  int    `json:"function"`
	SlaverId  byte   `json:"slaverId"`
	Address   uint16 `json:"address"`
	Frequency int64  `json:"frequency"`
	Quantity  uint16 `json:"quantity"`
	Value     string `json:"value,omitempty"`
}
type generic_modbus_device struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	//
	rtuHandler *modbus.RTUClientHandler
	tcpHandler *modbus.TCPClientHandler
	Client     modbus.Client
	//
	mainConfig   _GMODConfig
	retryTimes   int
	hwPortConfig hwportmanager.UartConfig
	Registers    map[string]*common.RegisterRW
}

/*
*
* 温湿度传感器
*
 */
func NewGenericModbusDevice(e typex.RuleX) typex.XDevice {
	mdev := new(generic_modbus_device)
	mdev.RuleEngine = e
	mdev.mainConfig = _GMODConfig{
		CommonConfig: _GMODCommonConfig{
			AutoRequest: func() *bool {
				b := false
				return &b
			}(),
		},
		PortUuid:   "/dev/ttyS0",
		HostConfig: common.HostConfig{Host: "127.0.0.1", Port: 502, Timeout: 3000},
	}
	mdev.Registers = map[string]*common.RegisterRW{}
	mdev.Busy = false
	mdev.status = typex.DEV_DOWN
	return mdev
}

//  初始化
func (mdev *generic_modbus_device) Init(devId string, configMap map[string]interface{}) error {
	mdev.PointId = devId
	modbuscache.RegisterSlot(mdev.PointId)
	if err := utils.BindSourceConfig(configMap, &mdev.mainConfig); err != nil {
		return err
	}
	if !utils.SContains([]string{"UART", "TCP"}, mdev.mainConfig.CommonConfig.Mode) {
		return errors.New("unsupported mode, only can be one of 'TCP' or 'UART'")
	}
	// 合并数据库里面的点位表
	var list []ModbusPoint
	errDb := interdb.DB().Table("m_modbus_data_points").
		Where("device_uuid=?", devId).Find(&list).Error
	if errDb != nil {
		return errDb
	}
	for _, v := range list {
		// 频率不能太快
		if v.Frequency < 50 {
			return errors.New("'frequency' must grate than 50 millisecond")
		}
		mdev.Registers[v.UUID] = &common.RegisterRW{
			Tag:       v.Tag,
			Alias:     v.Alias,
			Function:  v.Function,
			SlaverId:  v.SlaverId,
			Address:   v.Address,
			Quantity:  v.Quantity,
			Frequency: v.Frequency,
		}
		Status := 0
		LastFetchTime := uint64(time.Now().UnixMilli())
		modbuscache.SetValue(mdev.PointId, v.UUID, modbuscache.RegisterPoint{
			UUID:          v.UUID,
			Status:        Status,
			LastFetchTime: LastFetchTime,
		})
	}

	if mdev.mainConfig.CommonConfig.Mode == "UART" {
		hwPort, err := hwportmanager.GetHwPort(mdev.mainConfig.PortUuid)
		if err != nil {
			return err
		}
		if hwPort.Busy {
			return fmt.Errorf("UART is busying now, Occupied By:%s", hwPort.OccupyBy)
		}
		switch tCfg := hwPort.Config.(type) {
		case hwportmanager.UartConfig:
			{
				mdev.hwPortConfig = tCfg
			}
		default:
			{
				return fmt.Errorf("invalid config:%s", hwPort.Config)
			}
		}
	}
	return nil
}

// 启动
func (mdev *generic_modbus_device) Start(cctx typex.CCTX) error {
	mdev.Ctx = cctx.Ctx
	mdev.CancelCTX = cctx.CancelCTX

	if mdev.mainConfig.CommonConfig.Mode == "UART" {
		hwPort, err := hwportmanager.GetHwPort(mdev.mainConfig.PortUuid)
		if err != nil {
			return err
		}
		if hwPort.Busy {
			return fmt.Errorf("UART is busying now, Occupied By:%s", hwPort.OccupyBy)
		}

		mdev.rtuHandler = modbus.NewRTUClientHandler(hwPort.Name)
		mdev.rtuHandler.BaudRate = mdev.hwPortConfig.BaudRate
		mdev.rtuHandler.DataBits = mdev.hwPortConfig.DataBits
		mdev.rtuHandler.Parity = mdev.hwPortConfig.Parity
		mdev.rtuHandler.StopBits = mdev.hwPortConfig.StopBits
		// timeout 最大不能超过20, 不然无意义
		mdev.rtuHandler.Timeout = time.Duration(mdev.hwPortConfig.Timeout) * time.Millisecond
		if core.GlobalConfig.AppDebugMode {
			mdev.rtuHandler.Logger = golog.New(glogger.GLogger.Writer(),
				"Modbus RTU Mode: "+mdev.PointId+", ", golog.LstdFlags)
		}

		if err := mdev.rtuHandler.Connect(); err != nil {
			return err
		}
		hwportmanager.SetInterfaceBusy(mdev.mainConfig.PortUuid, hwportmanager.HwPortOccupy{
			UUID: mdev.PointId,
			Type: "DEVICE",
			Name: mdev.Details().Name,
		})
		mdev.Client = modbus.NewClient(mdev.rtuHandler)
	}
	if mdev.mainConfig.CommonConfig.Mode == "TCP" {
		mdev.tcpHandler = modbus.NewTCPClientHandler(
			fmt.Sprintf("%s:%v", mdev.mainConfig.HostConfig.Host, mdev.mainConfig.HostConfig.Port),
		)
		if core.GlobalConfig.AppDebugMode {
			mdev.tcpHandler.Logger = golog.New(glogger.GLogger.Writer(),
				"Modbus TCP Mode: "+mdev.PointId+", ", golog.LstdFlags)
		}
		if err := mdev.tcpHandler.Connect(); err != nil {
			return err
		}
		mdev.Client = modbus.NewClient(mdev.tcpHandler)
	}
	//---------------------------------------------------------------------------------
	// Start
	//---------------------------------------------------------------------------------
	if *mdev.mainConfig.CommonConfig.AutoRequest {
		mdev.retryTimes = 0
		go func(ctx context.Context) {
			buffer := make([]byte, common.T_64KB)
			for {
				select {
				case <-ctx.Done():
					{
						return
					}
				default:
					{
					}
				}
				n := 0
				var err error
				if mdev.mainConfig.CommonConfig.Mode == "TCP" {
					n, err = mdev.RTURead(buffer)
				}
				if mdev.mainConfig.CommonConfig.Mode == "UART" {
					n, err = mdev.TCPRead(buffer)
				}
				if err != nil {
					glogger.GLogger.Error(err)
					mdev.retryTimes++
					continue
				}
				// [] {} ""
				if n < 3 {
					continue
				}
				mdev.RuleEngine.WorkDevice(mdev.Details(), string(buffer[:n]))
			}

		}(mdev.Ctx)
	}

	mdev.status = typex.DEV_UP
	return nil
}

// 从设备里面读数据出来
func (mdev *generic_modbus_device) OnRead(cmd []byte, data []byte) (int, error) {
	return 0, nil
}

// 把数据写入设备
func (mdev *generic_modbus_device) OnWrite(cmd []byte, data []byte) (int, error) {
	RegisterW := common.RegisterW{}
	if err := json.Unmarshal(data, &RegisterW); err != nil {
		return 0, err
	}
	dataMap := [1]common.RegisterW{RegisterW}
	for _, r := range dataMap {
		if mdev.mainConfig.CommonConfig.Mode == "TCP" {
			mdev.tcpHandler.SlaveId = r.SlaverId
		}
		if mdev.mainConfig.CommonConfig.Mode == "UART" {
			mdev.rtuHandler.SlaveId = r.SlaverId
		}
		// 5
		if r.Function == common.WRITE_SINGLE_COIL {
			if len(r.Values) > 0 {
				if r.Values[0] == 0 {
					_, err := mdev.Client.WriteSingleCoil(r.Address,
						binary.BigEndian.Uint16([]byte{0x00, 0x00}))
					if err != nil {
						return 0, err
					}
				}
				if r.Values[0] == 1 {
					_, err := mdev.Client.WriteSingleCoil(r.Address,
						binary.BigEndian.Uint16([]byte{0xFF, 0x00}))
					if err != nil {
						return 0, err
					}
				}

			}

		}
		// 15
		if r.Function == common.WRITE_MULTIPLE_COILS {
			_, err := mdev.Client.WriteMultipleCoils(r.Address, r.Quantity, r.Values)
			if err != nil {
				return 0, err
			}
		}
		// 6
		if r.Function == common.WRITE_SINGLE_HOLDING_REGISTER {
			_, err := mdev.Client.WriteSingleRegister(r.Address, binary.BigEndian.Uint16(r.Values))
			if err != nil {
				return 0, err
			}
		}
		// 16
		if r.Function == common.WRITE_MULTIPLE_HOLDING_REGISTERS {

			_, err := mdev.Client.WriteMultipleRegisters(r.Address,
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

// 设备当前状态
func (mdev *generic_modbus_device) Status() typex.DeviceState {
	// 容错5次
	if mdev.retryTimes > 0 {
		return typex.DEV_DOWN
	}
	return mdev.status
}

// 停止设备
func (mdev *generic_modbus_device) Stop() {
	mdev.status = typex.DEV_DOWN
	if mdev.CancelCTX != nil {
		mdev.CancelCTX()
	}
	if mdev.mainConfig.CommonConfig.Mode == "UART" {
		hwportmanager.FreeInterfaceBusy(mdev.mainConfig.PortUuid)
	}
	if mdev.mainConfig.CommonConfig.Mode == "UART" {
		mdev.rtuHandler.Close()
	}
	if mdev.mainConfig.CommonConfig.Mode == "TCP" {
		mdev.tcpHandler.Close()
	}
	modbuscache.UnRegisterSlot(mdev.PointId)

}

// 设备属性，是一系列属性描述
func (mdev *generic_modbus_device) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (mdev *generic_modbus_device) Details() *typex.Device {
	return mdev.RuleEngine.GetDevice(mdev.PointId)
}

// 状态
func (mdev *generic_modbus_device) SetState(status typex.DeviceState) {
	mdev.status = status
}

// 驱动
func (mdev *generic_modbus_device) Driver() typex.XExternalDriver {
	return nil
}
func (mdev *generic_modbus_device) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
func (mdev *generic_modbus_device) OnCtrl([]byte, []byte) ([]byte, error) {
	return []byte{}, nil
}

/*
*
* 串口模式
*
 */
func (mdev *generic_modbus_device) modbusRead(buffer []byte) (int, error) {
	var err error
	var results []byte
	RegisterRWs := []common.RegisterRW{}
	count := len(mdev.Registers)
	if count == 0 {
		return 0, nil
	}
	if mdev.Client == nil {
		return 0, fmt.Errorf("modbus client id not valid")
	}
	for uuid, r := range mdev.Registers {
		if mdev.mainConfig.CommonConfig.Mode == "TCP" {
			mdev.tcpHandler.SlaveId = r.SlaverId
		}
		if mdev.mainConfig.CommonConfig.Mode == "UART" {
			mdev.rtuHandler.SlaveId = r.SlaverId
		}
		if r.Function == common.READ_COIL {
			results, err = mdev.Client.ReadCoils(r.Address, r.Quantity)
			if err != nil {
				count--
				glogger.GLogger.Error(err)
			}
			Value := covertEmptyHex(results)
			Reg := common.RegisterRW{
				Tag:      r.Tag,
				Function: r.Function,
				SlaverId: r.SlaverId,
				Address:  r.Address,
				Quantity: r.Quantity,
				Alias:    r.Alias,
				Value:    Value,
			}
			RegisterRWs = append(RegisterRWs, Reg)
			modbuscache.SetValue(mdev.PointId, uuid, modbuscache.RegisterPoint{
				UUID:          uuid,
				Status:        0,
				Value:         Value,
				LastFetchTime: uint64(time.Now().UnixMilli()),
			})
		}
		if r.Function == common.READ_DISCRETE_INPUT {
			results, err = mdev.Client.ReadDiscreteInputs(r.Address, r.Quantity)
			if err != nil {
				count--
				glogger.GLogger.Error(err)
				continue
			}
			Value := covertEmptyHex(results)
			Reg := common.RegisterRW{
				Tag:      r.Tag,
				Function: r.Function,
				SlaverId: r.SlaverId,
				Address:  r.Address,
				Quantity: r.Quantity,
				Alias:    r.Alias,
				Value:    Value,
			}
			RegisterRWs = append(RegisterRWs, Reg)
			modbuscache.SetValue(mdev.PointId, uuid, modbuscache.RegisterPoint{
				UUID:          uuid,
				Status:        0,
				Value:         Value,
				LastFetchTime: uint64(time.Now().UnixMilli()),
			})
		}
		if r.Function == common.READ_HOLDING_REGISTERS {
			results, err = mdev.Client.ReadHoldingRegisters(r.Address, r.Quantity)
			if err != nil {
				count--
				glogger.GLogger.Error(err)
			}
			Value := covertEmptyHex(results)
			Reg := common.RegisterRW{
				Tag:      r.Tag,
				Function: r.Function,
				SlaverId: r.SlaverId,
				Address:  r.Address,
				Quantity: r.Quantity,
				Alias:    r.Alias,
				Value:    Value,
			}
			RegisterRWs = append(RegisterRWs, Reg)
			modbuscache.SetValue(mdev.PointId, uuid, modbuscache.RegisterPoint{
				UUID:          uuid,
				Status:        0,
				Value:         Value,
				LastFetchTime: uint64(time.Now().UnixMilli()),
			})

		}
		if r.Function == common.READ_INPUT_REGISTERS {
			results, err = mdev.Client.ReadInputRegisters(r.Address, r.Quantity)
			if err != nil {
				count--
				glogger.GLogger.Error(err)
			}
			Value := covertEmptyHex(results)
			Reg := common.RegisterRW{
				Tag:      r.Tag,
				Function: r.Function,
				SlaverId: r.SlaverId,
				Address:  r.Address,
				Quantity: r.Quantity,
				Alias:    r.Alias,
				Value:    Value,
			}
			RegisterRWs = append(RegisterRWs, Reg)
			modbuscache.SetValue(mdev.PointId, uuid, modbuscache.RegisterPoint{
				UUID:          uuid,
				Status:        0,
				Value:         Value,
				LastFetchTime: uint64(time.Now().UnixMilli()),
			})
		}
		time.Sleep(time.Duration(r.Frequency) * time.Millisecond)
	}
	bytes, _ := json.Marshal(RegisterRWs)
	copy(buffer, bytes)
	return len(bytes), nil
}
func (mdev *generic_modbus_device) RTURead(buffer []byte) (int, error) {
	return mdev.modbusRead(buffer)
}
func (mdev *generic_modbus_device) TCPRead(buffer []byte) (int, error) {
	return mdev.modbusRead(buffer)
}
func covertEmptyHex(v []byte) string {
	if len(v) < 1 {
		return ""
	}
	return hex.EncodeToString(v)
}
