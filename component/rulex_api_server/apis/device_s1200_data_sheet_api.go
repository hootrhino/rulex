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

type SiemensPointVo struct {
	UUID          string `json:"uuid,omitempty"`
	DeviceUUID    string `json:"device_uuid"`
	Tag           string `json:"tag"`
	Alias         string `json:"alias"`
	Type          string `json:"type"`
	Frequency     *int64 `json:"frequency"`
	Address       *int   `json:"address"`
	Start         *int   `json:"start"`
	Size          *int   `json:"size"`
	Status        int    `json:"status"`        // 运行时数据
	LastFetchTime uint64 `json:"lastFetchTime"` // 运行时数据
	Value         string `json:"value"`         // 运行时数据
}

/*
*
* 特殊设备需要和外界交互，这里主要就是一些设备的点位表导入导出等支持
*  http://127.0.0.1:2580/api/v1/Siemens_data_sheet/export
 */

// SiemensPoints 获取Siemens_excel类型的点位数据
func SiemensPointsExport(c *gin.Context, ruleEngine typex.RuleX) {
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
// SELECT * FROM `m_Siemens_data_points` WHERE
// `m_Siemens_data_points`.`device_uuid` = "DEVICEDQNLO8"
// ORDER BY
// created_at DESC LIMIT 2 OFFSET 10
func SiemensSheetPageList(c *gin.Context, ruleEngine typex.RuleX) {
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	deviceUuid, _ := c.GetQuery("device_uuid")
	db := interdb.DB()
	tx := db.Scopes(service.Paginate(*pager))
	var count int64
	err1 := interdb.DB().Model(&model.MSiemensDataPoint{}).Count(&count).Error
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	var records []model.MSiemensDataPoint
	result := tx.Order("created_at DESC").Find(&records,
		&model.MSiemensDataPoint{DeviceUuid: deviceUuid})
	if result.Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(result.Error))
		return
	}
	recordsVo := []SiemensPointVo{}

	for _, record := range records {
		recordsVo = append(recordsVo, SiemensPointVo{
			UUID:          record.UUID,
			DeviceUUID:    record.DeviceUuid,
			Tag:           record.Tag,
			Type:          record.Type,
			Alias:         record.Alias,
			Address:       record.Address,
			Frequency:     record.Frequency,
			Start:         record.Start,
			Size:          record.Size,
			Status:        1,                              // 运行时
			LastFetchTime: uint64(time.Now().UnixMilli()), // 运行时
			Value:         "00000000",                     // 运行时
		})
	}
	Result := service.WrapPageResult(*pager, recordsVo, count)
	c.JSON(common.HTTP_OK, common.OkWithData(Result))
}

/*
*
* 删除单行
*
 */
func SiemensSheetDeleteAll(c *gin.Context, ruleEngine typex.RuleX) {
	type Form struct {
		UUIDs      []string `json:"uuids"`
		DeviceUUID string   `json:"device_uuid"`
	}
	form := Form{}
	if Error := c.ShouldBindJSON(&form); Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(Error))
		return
	}
	err := service.DeleteAllSiemensPointByDevice(form.DeviceUUID)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	ruleEngine.RestartDevice(form.DeviceUUID)
	c.JSON(common.HTTP_OK, common.Ok())

}
func SiemensSheetDelete(c *gin.Context, ruleEngine typex.RuleX) {
	type Form struct {
		UUIDs      []string `json:"uuids"`
		DeviceUUID string   `json:"device_uuid"`
	}
	form := Form{}
	if Error := c.ShouldBindJSON(&form); Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(Error))
		return
	}
	err := service.DeleteSiemensPointByDevice(form.UUIDs, form.DeviceUUID)
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
func SiemensSheetUpdate(c *gin.Context, ruleEngine typex.RuleX) {
	type Form struct {
		DeviceUUID        string           `json:"device_uuid"`
		SiemensDataPoints []SiemensPointVo `json:"siemens_data_points"`
	}
	form := Form{}
	// SiemensDataPoints := []SiemensPointVo{}
	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	for _, SiemensDataPoint := range form.SiemensDataPoints {
		if SiemensDataPoint.UUID == "" {
			NewRow := model.MSiemensDataPoint{}
			copier.Copy(&NewRow, &SiemensDataPoint)
			NewRow.DeviceUuid = SiemensDataPoint.DeviceUUID
			NewRow.UUID = utils.SiemensPointUUID()
			err0 := service.InsertSiemensPointPosition(NewRow)
			if err0 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err0))
				return
			}
		} else {
			OldRow := model.MSiemensDataPoint{}
			copier.Copy(&OldRow, &SiemensDataPoint)
			OldRow.DeviceUuid = SiemensDataPoint.DeviceUUID
			OldRow.UUID = SiemensDataPoint.UUID
			err0 := service.UpdateSiemensPoint(OldRow)
			if err0 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err0))
				return
			}
		}
	}
	ruleEngine.RestartDevice(form.DeviceUUID)
	c.JSON(common.HTTP_OK, common.Ok())

}

// SiemensSheetImport 上传Excel文件
func SiemensSheetImport(c *gin.Context, ruleEngine typex.RuleX) {
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
	if Device.Type != typex.SIEMENS_PLC.String() {
		c.JSON(common.HTTP_OK,
			common.Error("Invalid Device Type, Only Support Import Siemens Device"))
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
	list, err := parseSiemensPointExcel(file, "Sheet1", deviceUuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err = service.InsertSiemensPointPositions(list); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	ruleEngine.RestartDevice(deviceUuid)
	c.JSON(common.HTTP_OK, common.Ok())
}

func parseSiemensPointExcel(
	r io.Reader,
	sheetName string,
	deviceUuid string) (list []model.MSiemensDataPoint, err error) {
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
	//
	err1 := errors.New("invalid Sheet Header")
	if len(rows[0]) < 7 {
		return nil, err1
	}
	//tag alias type frequency address size
	if rows[0][0] != "tag" ||
		rows[0][1] != "alias" ||
		rows[0][2] != "type" ||
		rows[0][3] != "frequency" ||
		rows[0][4] != "address" ||
		rows[0][5] != "start" ||
		rows[0][6] != "size" {
		return nil, err1
	}

	list = make([]model.MSiemensDataPoint, 0)
	//tag alias type frequency address start size
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		tag := row[0]
		alias := row[1]
		DbType := row[2]
		frequency, _ := strconv.ParseInt(row[3], 10, 8)
		address, _ := strconv.ParseUint(row[4], 10, 16)
		start, _ := strconv.ParseUint(row[5], 10, 16)
		size, _ := strconv.ParseUint(row[6], 10, 16)
		Address := int(address)
		Frequency := int64(frequency)
		Start := int(start)
		Size := int(size)
		model := model.MSiemensDataPoint{
			UUID:       utils.SiemensPointUUID(),
			DeviceUuid: deviceUuid,
			Tag:        tag,
			Alias:      alias,
			Type:       DbType,
			Frequency:  &Frequency,
			Address:    &Address,
			Start:      &Start,
			Size:       &Size,
		}
		list = append(list, model)
	}
	return list, nil
}
