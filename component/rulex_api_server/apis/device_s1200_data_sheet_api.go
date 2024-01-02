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
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	siemenscache "github.com/hootrhino/rulex/component/intercache/siemens"
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
	UUID           string `json:"uuid"`
	DeviceUUID     string `json:"device_uuid"`
	SiemensAddress string `json:"siemensAddress"` // 西门子的地址字符串
	Tag            string `json:"tag"`
	Alias          string `json:"alias"`
	DataOrder      string `json:"dataOrder"` // 字节序
	DataType       string `json:"dataType"`
	Frequency      *int64 `json:"frequency"`
	Status         int    `json:"status"`        // 运行时数据
	LastFetchTime  uint64 `json:"lastFetchTime"` // 运行时数据
	Value          string `json:"value"`         // 运行时数据
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
		Slot := siemenscache.GetSlot(deviceUuid)
		Value, ok := Slot[record.UUID]
		Vo := SiemensPointVo{
			UUID:           record.UUID,
			DeviceUUID:     record.DeviceUuid,
			SiemensAddress: record.SiemensAddress,
			Tag:            record.Tag,
			Alias:          record.Alias,
			Frequency:      record.Frequency,
			DataType:       record.DataBlockType,
			DataOrder:      record.DataBlockOrder,
			LastFetchTime:  Value.LastFetchTime, // 运行时
			Value:          Value.Value,         // 运行时
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
func parseRequestSizeByType(s string) (int, error) {
	switch s {
	case "BYTE":
		return 1, nil
	case "SHORT":
		return 2, nil
	case "INT":
		return 4, nil
	case "FLOAT":
		return 4, nil
	default:
		return 0, errors.New("Invalid Block Type")
	}
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
	if len(rows[0]) < 6 {
		return nil, err1
	}
	// Address Tag Alias Type Order Frequency

	if strings.ToLower(rows[0][0]) != "address" ||
		strings.ToLower(rows[0][1]) != "tag" ||
		strings.ToLower(rows[0][2]) != "alias" ||
		strings.ToLower(rows[0][3]) != "type" ||
		strings.ToLower(rows[0][4]) != "order" ||
		strings.ToLower(rows[0][5]) != "frequency" {
		return nil, err1
	}

	list = make([]model.MSiemensDataPoint, 0)
	// Address Tag Alias Type Order Frequency
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		SiemensAddress := row[0]
		Tag := row[1]
		Alias := row[2]
		Type := row[3]
		Order := row[4]
		frequency, _ := strconv.ParseInt(row[5], 10, 8)
		Frequency := int64(frequency)
		Info, errParse1 := utils.ParseSiemensDB(SiemensAddress)
		if errParse1 != nil {
			return nil, errParse1
		}
		_, errParse2 := utils.ParseRequestSize(Info.DataBlockType)
		if errParse2 != nil {
			return nil, errParse2
		}
		model := model.MSiemensDataPoint{
			UUID:           utils.SiemensPointUUID(),
			DeviceUuid:     deviceUuid,
			SiemensAddress: SiemensAddress,
			Tag:            Tag,
			Alias:          Alias,
			DataBlockType:  Type,
			DataBlockOrder: Order,
			Frequency:      &Frequency,
		}
		list = append(list, model)
	}
	return list, nil
}
