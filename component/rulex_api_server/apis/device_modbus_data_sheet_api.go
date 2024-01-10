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

package apis

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	modbuscache "github.com/hootrhino/rulex/component/intercache/modbus"
	"github.com/hootrhino/rulex/component/interdb"
	common "github.com/hootrhino/rulex/component/rulex_api_server/common"
	"github.com/hootrhino/rulex/component/rulex_api_server/model"
	"github.com/hootrhino/rulex/component/rulex_api_server/service"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"github.com/xuri/excelize/v2"
)

type ModbusPointVo struct {
	UUID          string   `json:"uuid,omitempty"`
	DeviceUUID    string   `json:"device_uuid"`
	Tag           string   `json:"tag"`
	Alias         string   `json:"alias"`
	Function      *int     `json:"function"`
	SlaverId      *byte    `json:"slaverId"`
	Address       *uint16  `json:"address"`
	Frequency     *int64   `json:"frequency"`
	Quantity      *uint16  `json:"quantity"`
	Type          string   `json:"type"`          // 数据类型
	Order         string   `json:"order"`         // 字节序
	Weight        *float64 `json:"weight"`        // 权重
	Status        int      `json:"status"`        // 运行时数据
	LastFetchTime uint64   `json:"lastFetchTime"` // 运行时数据
	Value         string   `json:"value"`         // 运行时数据

}

/*
*
* 特殊设备需要和外界交互，这里主要就是一些设备的点位表导入导出等支持
*  http://127.0.0.1:2580/api/v1/modbus_data_sheet/export
 */

// ModbusPoints 获取modbus_excel类型的点位数据
func ModbusPointsExport(c *gin.Context, ruleEngine typex.RuleX) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment;filename=%v.csv", time.Now().UnixMilli()))
	csvWriter := csv.NewWriter(c.Writer)
	csvWriter.WriteAll([][]string{
		{"h1", "h2", "h3"},
		{"11", "12", "13"},
		{"21", "22", "23"},
		{"31", "32", "33"},
	})
	csvWriter.Flush()
}

// 分页获取
// SELECT * FROM `m_modbus_data_points` WHERE
// `m_modbus_data_points`.`device_uuid` = "DEVICEDQNLO8"
// ORDER BY
// created_at DESC LIMIT 2 OFFSET 10
func ModbusSheetPageList(c *gin.Context, ruleEngine typex.RuleX) {
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	deviceUuid, _ := c.GetQuery("device_uuid")
	db := interdb.DB()
	tx := db.Scopes(service.Paginate(*pager))
	var count int64
	err1 := interdb.DB().Model(&model.MModbusDataPoint{}).
		Where("device_uuid=?", deviceUuid).Count(&count).Error
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	var records []model.MModbusDataPoint
	result := tx.Order("created_at DESC").Find(&records,
		&model.MModbusDataPoint{DeviceUuid: deviceUuid})
	if result.Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(result.Error))
		return
	}
	recordsVo := []ModbusPointVo{}

	for _, record := range records {
		Slot := modbuscache.GetSlot(deviceUuid)
		Value, ok := Slot[record.UUID]
		Vo := ModbusPointVo{
			UUID:          record.UUID,
			DeviceUUID:    record.DeviceUuid,
			Tag:           record.Tag,
			Alias:         record.Alias,
			Function:      record.Function,
			SlaverId:      record.SlaverId,
			Address:       record.Address,
			Frequency:     record.Frequency,
			Quantity:      record.Quantity,
			Type:          record.Type,
			Order:         record.Order,
			Weight:        record.Weight,
			LastFetchTime: Value.LastFetchTime, // 运行时
			Value:         Value.Value,         // 运行时
		}
		if ok {
			Vo.Status = func() int {
				if Value.Value == "" {
					return 0
				}
				return 1
			}() // 运行时
			Vo.LastFetchTime = Value.LastFetchTime // 运行时
			Vo.Value = Value.Value                 // 运行时
			recordsVo = append(recordsVo, Vo)
		} else {
			recordsVo = append(recordsVo, Vo)
		}
	}
	Result := service.WrapPageResult(*pager, recordsVo, count)
	c.JSON(common.HTTP_OK, common.OkWithData(Result))
}

/*
*
* 删除单行
*
 */
func ModbusSheetDeleteAll(c *gin.Context, ruleEngine typex.RuleX) {
	type Form struct {
		UUIDs      []string `json:"uuids"`
		DeviceUUID string   `json:"device_uuid"`
	}
	form := Form{}
	if Error := c.ShouldBindJSON(&form); Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(Error))
		return
	}
	err := service.DeleteAllModbusPointByDevice(form.DeviceUUID)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	ruleEngine.RestartDevice(form.DeviceUUID)
	c.JSON(common.HTTP_OK, common.Ok())

}
func ModbusSheetDelete(c *gin.Context, ruleEngine typex.RuleX) {
	type Form struct {
		UUIDs      []string `json:"uuids"`
		DeviceUUID string   `json:"device_uuid"`
	}
	form := Form{}
	if Error := c.ShouldBindJSON(&form); Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(Error))
		return
	}
	err := service.DeleteModbusPointByDevice(form.UUIDs, form.DeviceUUID)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	ruleEngine.RestartDevice(form.DeviceUUID)
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 更新点位表
*
 */
func ModbusSheetUpdate(c *gin.Context, ruleEngine typex.RuleX) {
	type Form struct {
		DeviceUUID       string          `json:"device_uuid"`
		ModbusDataPoints []ModbusPointVo `json:"modbus_data_points"`
	}
	// ModbusDataPoints := []ModbusPointVo{}
	form := Form{}
	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	for _, ModbusDataPoint := range form.ModbusDataPoints {
		if ModbusDataPoint.UUID == "" {
			NewRow := model.MModbusDataPoint{
				UUID:       utils.ModbusPointUUID(),
				DeviceUuid: ModbusDataPoint.DeviceUUID,
				Tag:        ModbusDataPoint.Tag,
				Alias:      ModbusDataPoint.Alias,
				Function:   ModbusDataPoint.Function,
				SlaverId:   ModbusDataPoint.SlaverId,
				Address:    ModbusDataPoint.Address,
				Frequency:  ModbusDataPoint.Frequency,
				Quantity:   ModbusDataPoint.Quantity,
				Type:       ModbusDataPoint.Type,
				Order:      ModbusDataPoint.Order,
				Weight:     ModbusDataPoint.Weight,
			}
			err0 := service.InsertModbusPointPosition(NewRow)
			if err0 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err0))
				return
			}
		} else {
			OldRow := model.MModbusDataPoint{
				UUID:       ModbusDataPoint.UUID,
				DeviceUuid: ModbusDataPoint.DeviceUUID,
				Tag:        ModbusDataPoint.Tag,
				Alias:      ModbusDataPoint.Alias,
				Function:   ModbusDataPoint.Function,
				SlaverId:   ModbusDataPoint.SlaverId,
				Address:    ModbusDataPoint.Address,
				Frequency:  ModbusDataPoint.Frequency,
				Quantity:   ModbusDataPoint.Quantity,
				Type:       ModbusDataPoint.Type,
				Order:      ModbusDataPoint.Order,
				Weight:     ModbusDataPoint.Weight,
			}
			err0 := service.UpdateModbusPoint(OldRow)
			if err0 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err0))
				return
			}
		}
	}
	ruleEngine.RestartDevice(form.DeviceUUID)
	c.JSON(common.HTTP_OK, common.Ok())

}

// ModbusSheetImport 上传Excel文件
func ModbusSheetImport(c *gin.Context, ruleEngine typex.RuleX) {
	// 解析 multipart/form-data 类型的请求体
	err := c.Request.ParseMultipartForm(1024 * 1024 * 10)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	defer file.Close()
	deviceUuid := c.Request.Form.Get("device_uuid")
	type DeviceDto struct {
		UUID string
		Name string
		Type string
	}
	Device := DeviceDto{}
	errDb := interdb.DB().Table("m_devices").
		Where("uuid=?", deviceUuid).Find(&Device).Error
	if errDb != nil {
		c.JSON(common.HTTP_OK, common.Error400(errDb))
		return
	}
	if Device.Type != typex.GENERIC_MODBUS.String() {
		c.JSON(common.HTTP_OK,
			common.Error("Invalid Device Type, Only Support Import Modbus Device"))
		return
	}
	contentType := header.Header.Get("Content-Type")
	if contentType != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" &&
		contentType != "application/vnd.ms-excel" {
		c.JSON(common.HTTP_OK, common.Error("File Must be Excel Sheet"))
		return
	}
	// 判断文件大小是否符合要求（10MB）
	if header.Size > 1024*1024*10 {
		c.JSON(common.HTTP_OK, common.Error("Excel file size cannot be greater than 10MB"))
		return
	}
	list, err := parseModbusPointExcel(file, "Sheet1", deviceUuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err = service.InsertModbusPointPositions(list); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	ruleEngine.RestartDevice(deviceUuid)
	c.JSON(common.HTTP_OK, common.Ok())
}

func parseModbusPointExcel(
	r io.Reader,
	sheetName string,
	deviceUuid string) (list []model.MModbusDataPoint, err error) {
	excelFile, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	defer func() {
		excelFile.Close()
	}()
	// 读取表格
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	// 判断首行标头
	// tag, alias, function, frequency, slaverId, address, quality
	err1 := errors.New("invalid Sheet Header")
	if len(rows[0]) < 10 {
		return nil, err1
	}
	if rows[0][0] != "tag" ||
		rows[0][1] != "alias" ||
		rows[0][2] != "function" ||
		rows[0][3] != "frequency" ||
		rows[0][4] != "slaverId" ||
		rows[0][5] != "address" ||
		rows[0][6] != "quality" ||
		rows[0][7] != "type" ||
		rows[0][8] != "order" ||
		rows[0][9] != "weight" {
		return nil, err1
	}

	list = make([]model.MModbusDataPoint, 0)
	// tag, alias, function, frequency, slaverId, address, quality
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		tag := row[0]
		alias := row[1]
		function, _ := strconv.ParseUint(row[2], 10, 8)
		frequency, _ := strconv.ParseUint(row[3], 10, 8)
		slaverId, _ := strconv.ParseUint(row[4], 10, 8)
		address, _ := strconv.ParseUint(row[5], 10, 16)
		quantity, _ := strconv.ParseUint(row[6], 10, 16)
		Type := row[7]
		Order := row[8]
		Weight, _ := strconv.ParseFloat(row[9], 32)
		if Weight == 0 {
			Weight = 1 // 防止解析异常的时候系数0
		}
		Function := int(function)
		SlaverId := byte(slaverId)
		Address := uint16(address)
		Frequency := int64(frequency)
		Quantity := uint16(quantity)
		model := model.MModbusDataPoint{
			UUID:       utils.ModbusPointUUID(),
			DeviceUuid: deviceUuid,
			Tag:        tag,
			Alias:      alias,
			Function:   &Function,
			SlaverId:   &SlaverId,
			Address:    &Address,
			Frequency:  &Frequency, //ms
			Quantity:   &Quantity,
			Type:       Type,
			Order:      Order,
			Weight:     &Weight,
		}
		list = append(list, model)
	}
	return list, nil
}
