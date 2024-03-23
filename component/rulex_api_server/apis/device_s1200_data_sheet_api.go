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
	"errors"
	"fmt"
	"github.com/hootrhino/rulex/glogger"
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
	UUID           string   `json:"uuid"`
	DeviceUUID     string   `json:"device_uuid"`
	SiemensAddress string   `json:"siemensAddress"` // 西门子的地址字符串
	Tag            string   `json:"tag"`
	Alias          string   `json:"alias"`
	DataOrder      string   `json:"dataOrder"` // 字节序
	DataType       string   `json:"dataType"`
	Frequency      *int64   `json:"frequency"`
	Weight         *float64 `json:"weight"`        // 权重
	Status         int      `json:"status"`        // 运行时数据
	LastFetchTime  uint64   `json:"lastFetchTime"` // 运行时数据
	Value          string   `json:"value"`         // 运行时数据
}

/*
*
* 特殊设备需要和外界交互，这里主要就是一些设备的点位表导入导出等支持
*  http://127.0.0.1:2580/api/v1/Siemens_data_sheet/export
 */

// SiemensPoints 获取Siemens_excel类型的点位数据
func SiemensPointsExport(c *gin.Context, ruleEngine typex.RuleX) {
	deviceUuid, _ := c.GetQuery("device_uuid")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment;filename=%v.xlsx",
		time.Now().UnixMilli()))
	var records []model.MSiemensDataPoint
	result := interdb.DB().Order("created_at DESC").Find(&records,
		&model.MSiemensDataPoint{DeviceUuid: deviceUuid})
	if result.Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(result.Error))
		return
	}
	Headers := []string{
		"address", "tag", "alias", "type", "order", "weight", "frequency",
	}
	xlsx := excelize.NewFile()
	defer func() {
		if err := xlsx.Close(); err != nil {
			glogger.GLogger.Errorf("close excel file, err=%v", err)
		}
	}()
	cell, _ := excelize.CoordinatesToCellName(1, 1)
	xlsx.SetSheetRow("Sheet1", cell, &Headers)

	if len(records) > 1 {
		for idx, record := range records[0:] {
			Row := []string{
				record.SiemensAddress,
				record.Tag,
				record.Alias,
				record.DataBlockType,
				record.DataBlockOrder,
				fmt.Sprintf("%f", *record.Weight),
				fmt.Sprintf("%d", *record.Frequency),
			}
			cell, _ = excelize.CoordinatesToCellName(1, idx+2)
			xlsx.SetSheetRow("Sheet1", cell, &Row)
		}
	}

	xlsx.WriteTo(c.Writer)

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
	err1 := interdb.DB().Model(&model.MSiemensDataPoint{}).
		Where("device_uuid=?", deviceUuid).Count(&count).Error
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
			Weight:         record.Weight,
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
* 检查点位合法性
*
 */
func checkSiemensDataPoints(M SiemensPointVo) error {
	if M.Tag == "" {
		return fmt.Errorf("Missing required param 'tag'")
	}
	if len(M.Tag) > 256 {
		return fmt.Errorf("Tag length must range of 1-256")
	}
	if M.Alias == "" {
		return fmt.Errorf("Missing required param 'alias'")
	}
	if len(M.Alias) > 256 {
		return fmt.Errorf("Alias length must range of 1-256")
	}
	if M.SiemensAddress == "" {
		return fmt.Errorf("Missing required param 'address'")
	}

	if M.Frequency == nil {
		return fmt.Errorf("Missing required param 'frequency'")
	}
	if *M.Frequency < 50 {
		return fmt.Errorf("Frequency must greater than 50ms")
	}
	if *M.Frequency > 100000 {
		return fmt.Errorf("Frequency must little than 100s")
	}

	switch M.DataType {
	case "I", "Q", "BYTE":
		if M.DataOrder != "A" {
			return fmt.Errorf("invalid '%s' order '%s'", M.DataType, M.DataOrder)
		}
	case "SHORT", "USHORT", "INT16", "UINT16":
		if !utils.SContains([]string{"AB", "BA"}, M.DataOrder) {
			return fmt.Errorf("'Invalid '%s' order '%s'", M.DataType, M.DataOrder)
		}
	case "RAW", "INT", "INT32", "UINT", "UINT32", "FLOAT", "UFLOAT":
		if !utils.SContains([]string{"ABCD", "DCBA", "CDAB"}, M.DataOrder) {
			return fmt.Errorf("invalid '%s' order '%s'", M.DataType, M.DataOrder)
		}
	default:
		return fmt.Errorf("invalid '%s' order '%s'", M.DataType, M.DataOrder)
	}
	if M.Weight == nil {
		return fmt.Errorf("invalid Weight value:%d", M.Weight)
	}
	if !utils.IsValidColumnName(M.Tag) {
		return fmt.Errorf("'Invalid Tag Name:%d", M.Tag)
	}
	return nil
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
		if err := checkSiemensDataPoints(SiemensDataPoint); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		if SiemensDataPoint.UUID == "" {
			NewRow := model.MSiemensDataPoint{}
			copier.Copy(&NewRow, &SiemensDataPoint)
			NewRow.DeviceUuid = SiemensDataPoint.DeviceUUID
			NewRow.UUID = utils.SiemensPointUUID()
			NewRow.DataBlockType = SiemensDataPoint.DataType
			NewRow.DataBlockOrder = SiemensDataPoint.DataOrder
			NewRow.Weight = SiemensDataPoint.Weight
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
			OldRow.DataBlockType = SiemensDataPoint.DataType
			OldRow.DataBlockOrder = SiemensDataPoint.DataOrder
			OldRow.Weight = utils.HandleZeroValue(SiemensDataPoint.Weight)

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
	// Address Tag Alias Type Order Frequency

	if strings.ToLower(rows[0][0]) != "address" ||
		strings.ToLower(rows[0][1]) != "tag" ||
		strings.ToLower(rows[0][2]) != "alias" ||
		strings.ToLower(rows[0][3]) != "type" ||
		strings.ToLower(rows[0][4]) != "order" ||
		strings.ToLower(rows[0][5]) != "weight" ||
		strings.ToLower(rows[0][6]) != "frequency" {
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
		Weight, _ := strconv.ParseFloat(row[5], 32)
		if Weight == 0 {
			Weight = 1 // 防止解析异常的时候系数0
		}
		frequency, _ := strconv.ParseUint(row[6], 10, 8)
		Frequency := int64(frequency)
		_, errParse1 := utils.ParseSiemensDB(SiemensAddress)
		if errParse1 != nil {
			return nil, errParse1
		}
		model := model.MSiemensDataPoint{
			UUID:           utils.SiemensPointUUID(),
			DeviceUuid:     deviceUuid,
			SiemensAddress: SiemensAddress,
			Tag:            Tag,
			Alias:          Alias,
			DataBlockType:  Type,
			DataBlockOrder: utils.GetDefaultDataOrder(Type, Order),
			Frequency:      &Frequency,
			Weight:         &Weight,
		}
		list = append(list, model)
	}
	return list, nil
}
