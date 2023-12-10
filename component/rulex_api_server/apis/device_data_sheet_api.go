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
	"github.com/hootrhino/rulex/component/interdb"
	common "github.com/hootrhino/rulex/component/rulex_api_server/common"
	"github.com/hootrhino/rulex/component/rulex_api_server/model"
	"github.com/hootrhino/rulex/component/rulex_api_server/service"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"github.com/jinzhu/copier"
	"github.com/xuri/excelize/v2"
)

type ModbusPointVo struct {
	UUID       string `json:"uuid,omitempty"`
	DeviceUUID string `json:"device_uuid"`
	Tag        string `json:"tag"`
	Alias      string `json:"alias"`
	Function   int    `json:"function"`
	SlaverId   byte   `json:"slaverId"`
	Address    uint16 `json:"address"`
	Frequency  int64  `json:"frequency"`
	Quantity   uint16 `json:"quantity"`
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
	err1 := interdb.DB().Model(&model.MModbusDataPoint{}).Count(&count).Error
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
	Result := service.WrapPageResult(*pager, records, count)
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
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 更新点位表
*
 */
func ModbusSheetUpdate(c *gin.Context, ruleEngine typex.RuleX) {
	ModbusDataPoints := []ModbusPointVo{}
	err := c.ShouldBindJSON(&ModbusDataPoints)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	for _, ModbusDataPoint := range ModbusDataPoints {
		if ModbusDataPoint.UUID == "" {
			NewRow := model.MModbusDataPoint{}
			copier.Copy(&NewRow, &ModbusDataPoint)
			NewRow.DeviceUuid = ModbusDataPoint.DeviceUUID
			NewRow.UUID = utils.ModbusPointUUID()
			service.InsertModbusPointPosition(NewRow)
		} else {
			OldRow := model.MModbusDataPoint{}
			copier.Copy(&OldRow, &ModbusDataPoint)
			OldRow.DeviceUuid = ModbusDataPoint.DeviceUUID
			OldRow.UUID = ModbusDataPoint.UUID
			service.UpdateModbusPoint(OldRow)
		}
	}
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
	// tag, alias, function, frequency, slaverId, startAddress, quality
	err1 := errors.New("invalid Sheet Header")
	if len(rows[0]) < 7 {
		return nil, err1
	}
	if rows[0][0] != "tag" ||
		rows[0][1] != "alias" ||
		rows[0][2] != "function" ||
		rows[0][3] != "frequency" ||
		rows[0][4] != "slaverId" ||
		rows[0][5] != "startAddress" ||
		rows[0][6] != "quality" {
		return nil, err1
	}

	list = make([]model.MModbusDataPoint, 0)
	// tag, alias, function, frequency, slaverId, startAddress, quality
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		tag := row[0]
		alias := row[1]
		function, _ := strconv.ParseInt(row[2], 10, 8)
		frequency, _ := strconv.ParseInt(row[3], 10, 8)
		slaverId, _ := strconv.ParseInt(row[4], 10, 8)
		address, _ := strconv.ParseUint(row[5], 10, 16)
		quantity, _ := strconv.ParseUint(row[6], 10, 16)
		model := model.MModbusDataPoint{
			UUID:       utils.ModbusPointUUID(),
			DeviceUuid: deviceUuid,
			Tag:        tag,
			Alias:      alias,
			Function:   int(function),
			SlaverId:   byte(slaverId),
			Address:    uint16(address),
			Frequency:  frequency, //ms
			Quantity:   uint16(quantity),
		}
		list = append(list, model)
	}
	return list, nil
}
