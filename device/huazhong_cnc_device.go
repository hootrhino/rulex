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
	"time"

	hnc8cache "github.com/hootrhino/rulex/component/intercache/hnccnc"
	"github.com/hootrhino/rulex/component/interdb"
	"github.com/hootrhino/rulex/component/iotschema"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

/*
*
* 华中数控点位表; 点表保存在全局Sqlite
*
 */
type hncDataPoint struct {
	UUID        string `json:"uuid,omitempty"`                   // 当UUID为空时新建
	Name        string `json:"name" validate:"required"`         // 点位名称
	Alias       string `json:"alias" validate:"required"`        // 点位名称
	ApiFunction string `json:"api_function" validate:"required"` // API路径
	Group       int    `json:"group" validate:"required"`        // 分组采集
	Address     string `json:"address" validate:"required"`      // 地址
}

// HNC8 NC-Link 接口参口手册 版本 v1.0 2019-02-13

type hnc8_cnc_config struct {
	SerialNumber  string                  `json:"serialNumber" validate:"required"` // CNC 序列号
	Host          string                  `json:"host" validate:"required"`         // IP:Port
	ApiVersion    int                     `json:"apiVersion" validate:"required"`   // API 版本,2 | 3
	AutoRequest   *bool                   `json:"autoRequest" validate:"required"`  // 是否自动请求
	Frequency     int64                   `json:"frequency" validate:"required"`    // 请求频率
	HncDataPoints map[string]hncDataPoint `json:"hncDataPoint" validate:"required"` // 点位表
}

/*
*
* 凯帝恩CNC
*
 */
type HNC8_CNC struct {
	typex.XStatus
	mainConfig hnc8_cnc_config
	status     typex.DeviceState
}

func NewHNC8_CNC(e typex.RuleX) typex.XDevice {
	hnc8Cnc := new(HNC8_CNC)
	hnc8Cnc.RuleEngine = e
	AutoRequest := true
	hnc8Cnc.mainConfig = hnc8_cnc_config{
		SerialNumber:  "HNC-1",
		Host:          "127.0.0.1:8080",
		ApiVersion:    3,
		AutoRequest:   &AutoRequest,
		Frequency:     1000,
		HncDataPoints: map[string]hncDataPoint{},
	}
	return hnc8Cnc
}

//  初始化
func (hnc8Cnc *HNC8_CNC) Init(devId string, configMap map[string]interface{}) error {
	hnc8Cnc.PointId = devId
	hnc8cache.RegisterSlot(hnc8Cnc.PointId)
	if err := utils.BindSourceConfig(configMap, &hnc8Cnc.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	var hncDataPointList []hncDataPoint
	errDb := interdb.DB().Table("m_hnc8_data_points").
		Where("device_uuid=?", devId).Find(&hncDataPointList).Error
	if errDb != nil {
		return errDb
	}
	// 加载点位
	for _, MHncDataPoint := range hncDataPointList {
		hnc8Cnc.mainConfig.HncDataPoints[MHncDataPoint.UUID] = hncDataPoint{
			UUID:        MHncDataPoint.UUID,
			Name:        MHncDataPoint.Name,
			Alias:       MHncDataPoint.Alias,
			ApiFunction: MHncDataPoint.ApiFunction,
			Group:       MHncDataPoint.Group,
			Address:     MHncDataPoint.Address,
		}
		LastFetchTime := uint64(time.Now().UnixMilli())
		hnc8cache.SetValue(hnc8Cnc.PointId,
			MHncDataPoint.UUID, hnc8cache.Hnc8RegisterPoint{
				UUID:          MHncDataPoint.UUID,
				Status:        0,
				LastFetchTime: LastFetchTime,
				Value:         "",
			})
	}
	return nil
}

// 启动
func (hnc8Cnc *HNC8_CNC) Start(cctx typex.CCTX) error {
	hnc8Cnc.Ctx = cctx.Ctx
	hnc8Cnc.CancelCTX = cctx.CancelCTX
	go func(Ctx context.Context) {
		for {
			select {
			case <-Ctx.Done():
				{
					return
				}
			default:
				{
				}
			}
			glogger.GLogger.Debug("调试模式, 采集数据中。。。。")
			time.Sleep(1000 * time.Millisecond)
		}

	}(hnc8Cnc.Ctx)
	hnc8Cnc.status = typex.DEV_UP
	return nil
}

func (hnc8Cnc *HNC8_CNC) OnRead(cmd []byte, data []byte) (int, error) {

	return 0, nil
}

// 把数据写入设备
func (hnc8Cnc *HNC8_CNC) OnWrite(cmd []byte, b []byte) (int, error) {
	return 0, nil
}

// 设备当前状态
func (hnc8Cnc *HNC8_CNC) Status() typex.DeviceState {
	return typex.DEV_UP
}

// 停止设备
func (hnc8Cnc *HNC8_CNC) Stop() {
	hnc8Cnc.status = typex.DEV_DOWN
	if hnc8Cnc.CancelCTX != nil {
		hnc8Cnc.CancelCTX()
	}
	hnc8cache.UnRegisterSlot(hnc8Cnc.PointId)
}

// 设备属性，是一系列属性描述
func (hnc8Cnc *HNC8_CNC) Property() []iotschema.IoTSchema {
	return []iotschema.IoTSchema{}
}

// 真实设备
func (hnc8Cnc *HNC8_CNC) Details() *typex.Device {
	return hnc8Cnc.RuleEngine.GetDevice(hnc8Cnc.PointId)
}

// 状态
func (hnc8Cnc *HNC8_CNC) SetState(status typex.DeviceState) {
	hnc8Cnc.status = status

}

// 驱动
func (hnc8Cnc *HNC8_CNC) Driver() typex.XExternalDriver {
	return nil
}

// --------------------------------------------------------------------------------------------------
//
// --------------------------------------------------------------------------------------------------

func (hnc8Cnc *HNC8_CNC) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
func (hnc8Cnc *HNC8_CNC) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}
