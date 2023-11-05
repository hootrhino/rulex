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

package service

import (
	"encoding/json"

	"github.com/hootrhino/rulex/component/interdb"
	"github.com/hootrhino/rulex/plugin/http_server/model"
)

type UartConfigDto struct {
	Timeout  int
	Uart     string
	BaudRate int
	DataBits int
	Parity   string
	StopBits int
}
type HwInterfaceDto struct {
	UUID     string
	Name     string        // 接口名称
	Type     string        // 接口类型, UART(串口),USB(USB),FD(通用文件句柄)
	Alias    string        // 别名
	Busy     bool          // 是否被占
	OccupyBy string        // 被谁占用了
	Config   UartConfigDto // 配置, 串口配置、或者网卡、USB等
}

func (u UartConfigDto) JsonString() string {
	if bytes, err := json.Marshal(u); err != nil {
		return ""
	} else {
		return string(bytes)
	}
}

/*
*
* 所有的接口列表配置
*
 */
func AllHwInterface() ([]model.MHwInterface, error) {
	Ifaces := []model.MHwInterface{}
	return Ifaces, interdb.DB().
		Model(&model.MHwInterface{}).Find(&Ifaces).Error
}

/*
*
* 配置WIFI HwInterface
*
 */
func UpdateHwInterfaceConfig(MHwInterface model.MHwInterface) error {
	Model := model.MHwInterface{}
	return interdb.DB().
		Model(Model).
		Where("uuid=? AND name=?",
			MHwInterface.UUID, MHwInterface.Name).
		Updates(MHwInterface).Error
}

/*
*
* 获取HwInterface的配置信息
*
 */
func GetHwInterfaceConfig(uuid string) (model.MHwInterface, error) {
	MHwInterface := model.MHwInterface{}
	err := interdb.DB().
		Where("uuid=?", uuid).
		Find(&MHwInterface).Error
	return MHwInterface, err
}

/*
*
* 初始化网卡配置参数
*
 */
func InitHwIfaceConfig() error {

	S1 := model.MHwInterface{
		UUID:        "/dev/ttyS1",
		Name:        "/dev/ttyS1",
		Type:        "UART",
		Alias:       "RS4851(A1B1)",
		Description: "RS4851号串口",
	}
	uartCfg1 := UartConfigDto{
		Timeout:  3000,
		Uart:     "/dev/ttyS1",
		BaudRate: 9600,
		DataBits: 8,
		Parity:   "N",
		StopBits: 1,
	}
	S1.Config = uartCfg1.JsonString()
	S2 := model.MHwInterface{
		UUID:        "/dev/ttyS2",
		Name:        "/dev/ttyS2",
		Type:        "UART",
		Alias:       "RS4851(A2B2)",
		Description: "RS4852号串口",
	}
	uartCfg2 := UartConfigDto{
		Timeout:  3000,
		Uart:     "/dev/ttyS2",
		BaudRate: 9600,
		DataBits: 8,
		Parity:   "N",
		StopBits: 1,
	}
	S2.Config = uartCfg2.JsonString()
	//
	err1 := interdb.DB().
		Table("m_hw_interfaces").Where("uuid", "/dev/ttyS1").
		FirstOrCreate(&S1).Error
	if err1 != nil {
		return err1
	}
	err2 := interdb.DB().
		Table("m_hw_interfaces").Where("uuid", "/dev/ttyS2").
		FirstOrCreate(&S2).Error
	if err2 != nil {
		return err2
	}
	return nil
}
