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
	"time"

	"github.com/gin-gonic/gin"
	hnccnc "github.com/hootrhino/rulex/component/intercache/hnccnc"
	"github.com/hootrhino/rulex/component/interdb"

	common "github.com/hootrhino/rulex/component/rulex_api_server/common"
	"github.com/hootrhino/rulex/component/rulex_api_server/model"
	"github.com/hootrhino/rulex/component/rulex_api_server/service"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"github.com/xuri/excelize/v2"
)

type Hnc8PointVo struct {
	UUID          string `json:"uuid,omitempty"`
	DeviceUUID    string `json:"device_uuid"`
	Name          string `json:"name"`
	Alias         string `json:"alias"`
	ApiFunction   string `json:"apiFunction"`
	Group         *int   `json:"group"`
	Address       string `json:"address"`
	Status        int    `json:"status"`        // 运行时数据
	LastFetchTime uint64 `json:"lastFetchTime"` // 运行时数据
	Value         string `json:"value"`         // 运行时数据
}

/*
*
* 特殊设备需要和外界交互，这里主要就是一些设备的点位表导入导出等支持
*  http://127.0.0.1:2580/api/v1/Hnc8_data_sheet/export
 */

// Hnc8Points 获取Hnc8_excel类型的点位数据
func Hnc8PointsExport(c *gin.Context, ruleEngine typex.RuleX) {
	deviceUuid, _ := c.GetQuery("device_uuid")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment;filename=%v.xlsx",
		time.Now().UnixMilli()))
	var records []model.MHnc8DataPoint
	result := interdb.DB().Order("created_at DESC").Find(&records,
		&model.MHnc8DataPoint{DeviceUuid: deviceUuid})
	if result.Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(result.Error))
		return
	}
	Headers := []string{
		"name", "alias", "function", "group", "address",
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
				record.Name, record.Alias, record.ApiFunction, record.Address,
			}
			cell, _ = excelize.CoordinatesToCellName(1, idx+2)
			xlsx.SetSheetRow("Sheet1", cell, &Row)
		}
	}
	xlsx.WriteTo(c.Writer)
}

// 分页获取
// SELECT * FROM `m_Hnc8_data_points` WHERE
// `m_Hnc8_data_points`.`device_uuid` = "DEVICEDQNLO8"
// ORDER BY
// created_at DESC LIMIT 2 OFFSET 10
func Hnc8SheetPageList(c *gin.Context, ruleEngine typex.RuleX) {
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	deviceUuid, _ := c.GetQuery("device_uuid")
	db := interdb.DB()
	tx := db.Scopes(service.Paginate(*pager))
	var count int64
	err1 := interdb.DB().Model(&model.MHnc8DataPoint{}).
		Where("device_uuid=?", deviceUuid).Count(&count).Error
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	var records []model.MHnc8DataPoint
	result := tx.Order("created_at DESC").Find(&records,
		&model.MHnc8DataPoint{DeviceUuid: deviceUuid})
	if result.Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(result.Error))
		return
	}
	recordsVo := []Hnc8PointVo{}

	for _, record := range records {
		Slot := hnccnc.GetSlot(deviceUuid)
		Value, ok := Slot[record.UUID]
		Vo := Hnc8PointVo{
			UUID:        record.UUID,
			DeviceUUID:  record.DeviceUuid,
			Name:        record.Name,
			Alias:       record.Alias,
			ApiFunction: record.ApiFunction,
			Group:       &record.Group,
			Address:     record.Address,
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
func Hnc8SheetDeleteAll(c *gin.Context, ruleEngine typex.RuleX) {
	type Form struct {
		UUIDs      []string `json:"uuids"`
		DeviceUUID string   `json:"device_uuid"`
	}
	form := Form{}
	if Error := c.ShouldBindJSON(&form); Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(Error))
		return
	}
	err := service.DeleteAllHnc8PointByDevice(form.DeviceUUID)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	ruleEngine.RestartDevice(form.DeviceUUID)
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
*删除
*
 */
func Hnc8SheetDelete(c *gin.Context, ruleEngine typex.RuleX) {
	type Form struct {
		UUIDs      []string `json:"uuids"`
		DeviceUUID string   `json:"device_uuid"`
	}
	form := Form{}
	if Error := c.ShouldBindJSON(&form); Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(Error))
		return
	}
	err := service.DeleteHnc8PointByDevice(form.UUIDs, form.DeviceUUID)
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
func checkHnc8DataPoints(M Hnc8PointVo) error {
	if M.Name == "" {
		return fmt.Errorf("'Missing required param 'name'")
	}
	if len(M.Name) > 256 {
		return fmt.Errorf("'Tag length must range of 1-256")
	}
	if M.Alias == "" {
		return fmt.Errorf("'Missing required param 'alias'")
	}
	if len(M.Alias) > 256 {
		return fmt.Errorf("'Alias length must range of 1-256")
	}
	return nil
}

/*
*
* 更新点位表
*
 */
func Hnc8SheetUpdate(c *gin.Context, ruleEngine typex.RuleX) {
	type Form struct {
		DeviceUUID     string        `json:"device_uuid"`
		Hnc8DataPoints []Hnc8PointVo `json:"Hnc8_data_points"`
	}
	// Hnc8DataPoints := []Hnc8PointVo{}
	form := Form{}
	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	for _, Hnc8DataPoint := range form.Hnc8DataPoints {
		if err := checkHnc8DataPoints(Hnc8DataPoint); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		if Hnc8DataPoint.UUID == "" {
			NewRow := model.MHnc8DataPoint{
				UUID:        utils.HNC8PointUUID(),
				DeviceUuid:  Hnc8DataPoint.DeviceUUID,
				Name:        Hnc8DataPoint.Name,
				Alias:       Hnc8DataPoint.Alias,
				ApiFunction: Hnc8DataPoint.ApiFunction,
				Group:       *Hnc8DataPoint.Group,
				Address:     Hnc8DataPoint.Address,
			}
			err0 := service.InsertHnc8PointPosition(NewRow)
			if err0 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err0))
				return
			}
		} else {
			OldRow := model.MHnc8DataPoint{
				UUID:        Hnc8DataPoint.UUID,
				DeviceUuid:  Hnc8DataPoint.DeviceUUID,
				Name:        Hnc8DataPoint.Name,
				Alias:       Hnc8DataPoint.Alias,
				ApiFunction: Hnc8DataPoint.ApiFunction,
				Group:       *Hnc8DataPoint.Group,
				Address:     Hnc8DataPoint.Address,
			}
			err0 := service.UpdateHnc8Point(OldRow)
			if err0 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err0))
				return
			}
		}
	}
	ruleEngine.RestartDevice(form.DeviceUUID)
	c.JSON(common.HTTP_OK, common.Ok())

}

// Hnc8SheetImport 上传Excel文件
func Hnc8SheetImport(c *gin.Context, ruleEngine typex.RuleX) {
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
	if Device.Type != typex.HNC8.String() {
		c.JSON(common.HTTP_OK,
			common.Error("Invalid Device Type, Only Support Import Hnc8 Device"))
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
	list, err := parseHnc8PointExcel(file, "Sheet1", deviceUuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err = service.InsertHnc8PointPositions(list); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	ruleEngine.RestartDevice(deviceUuid)
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 解析表格
*
 */

func parseHnc8PointExcel(r io.Reader, sheetName string,
	deviceUuid string) (list []model.MHnc8DataPoint, err error) {
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
	// name, alias, function, group, address
	err1 := errors.New("'Invalid Sheet Header, must follow fixed format: 【name, alias, function, group, address】")
	if len(rows[0]) < 5 {
		return nil, err1
	}
	// 严格检查表结构
	if rows[0][0] != "name" ||
		rows[0][1] != "alias" ||
		rows[0][2] != "function" ||
		rows[0][3] != "group" ||
		rows[0][4] != "address" {
		return nil, err1
	}

	list = make([]model.MHnc8DataPoint, 0)
	// name, alias, function, group, address
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		name := row[0]
		alias := row[1]
		function := row[2]
		group, _ := strconv.ParseUint(row[3], 10, 8)
		Group := int(group)
		address := row[4]
		if err := checkHnc8DataPoints(Hnc8PointVo{
			Name:        name,
			Alias:       alias,
			ApiFunction: function,
			Group:       &Group,
			Address:     address,
		}); err != nil {
			return nil, err
		}
		//
		model := model.MHnc8DataPoint{
			UUID:        utils.HNC8PointUUID(),
			DeviceUuid:  deviceUuid,
			Name:        name,
			Alias:       alias,
			ApiFunction: function,
			Group:       Group,
			Address:     address,
		}
		list = append(list, model)
	}
	return list, nil
}
