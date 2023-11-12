package apis

import (
	"errors"
	"io"
	"strconv"

	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/server"
	"github.com/hootrhino/rulex/plugin/http_server/service"

	"github.com/xuri/excelize/v2"

	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2/json"
)

type DeviceVo struct {
	UUID        string                 `json:"uuid"`
	Gid         string                 `json:"gid"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	State       int                    `json:"state"`
	Config      map[string]interface{} `json:"config"`
	Description string                 `json:"description"`
}

/*
*
* 列表先读数据库，然后读内存，合并状态后输出
*
 */
func DeviceDetail(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	mdev, err := service.GetMDeviceWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400EmptyObj(err))
		return
	}
	DeviceVo := DeviceVo{}
	DeviceVo.UUID = mdev.UUID
	DeviceVo.Name = mdev.Name
	DeviceVo.Type = mdev.Type
	DeviceVo.Description = mdev.Description
	DeviceVo.Config = mdev.GetConfig()
	//
	device := ruleEngine.GetDevice(mdev.UUID)
	if device == nil {
		DeviceVo.State = int(typex.DEV_STOP)
	} else {
		DeviceVo.State = int(device.Device.Status())
	}
	Group := service.GetVisualGroup(mdev.UUID)
	DeviceVo.Gid = Group.UUID
	c.JSON(common.HTTP_OK, common.OkWithData(DeviceVo))
}

/*
*
* 分组查看
*
 */
func ListDeviceByGroup(c *gin.Context, ruleEngine typex.RuleX) {
	Gid, _ := c.GetQuery("uuid")
	devices := []DeviceVo{}
	_, MDevices := service.FindByType(Gid, "DEVICE")
	for _, mdev := range MDevices {
		DeviceVo := DeviceVo{}
		DeviceVo.UUID = mdev.UUID
		DeviceVo.Name = mdev.Name
		DeviceVo.Type = mdev.Type
		DeviceVo.Description = mdev.Description
		DeviceVo.Config = mdev.GetConfig()
		//
		device := ruleEngine.GetDevice(mdev.UUID)
		if device == nil {
			DeviceVo.State = int(typex.DEV_STOP)
		} else {
			DeviceVo.State = int(device.Device.Status())
		}
		Group := service.GetVisualGroup(mdev.UUID)
		DeviceVo.Gid = Group.UUID
		devices = append(devices, DeviceVo)
	}
	c.JSON(common.HTTP_OK, common.OkWithData(devices))
}

// 删除设备
func DeleteDevice(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	Mdev, err := service.GetMDeviceWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 要处理这个空字符串 ""
	if Mdev.BindRules.Len() == 1 && len(Mdev.BindRules[0]) != 0 {
		c.JSON(common.HTTP_OK, common.Error("Can't remove, Already have rule bind:"+Mdev.BindRules.String()))
		return
	}
	// 检查是否有规则被绑定了
	for _, ruleId := range Mdev.BindRules {
		if ruleId != "" {
			_, err0 := service.GetMRuleWithUUID(ruleId)
			if err0 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err0))
				return
			}
		}

	}

	// 检查是否通用Modbus设备.需要同步删除点位表记录
	if Mdev.Type == "GENERIC_MODBUS_POINT_EXCEL" {
		if err := service.DeleteModbusPointAndDevice(uuid); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	} else {
		if err := service.DeleteDevice(uuid); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}

	old := ruleEngine.GetDevice(uuid)
	if old != nil {
		if old.Device.Status() == typex.DEV_UP {
			old.Device.Stop()
		}
	}

	ruleEngine.RemoveDevice(uuid)
	c.JSON(common.HTTP_OK, common.Ok())

}

// 创建设备
func CreateDevice(c *gin.Context, ruleEngine typex.RuleX) {

	form := DeviceVo{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	configJson, err := json.Marshal(form.Config)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	newUUID := utils.DeviceUuid()
	MDevice := model.MDevice{
		UUID:        newUUID,
		Type:        form.Type,
		Name:        form.Name,
		Description: form.Description,
		Config:      string(configJson),
		BindRules:   []string{},
	}
	if err := service.InsertDevice(&MDevice); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 新建大屏的时候必须给一个分组
	if err := service.BindResource(form.Gid, MDevice.UUID); err != nil {
		c.JSON(common.HTTP_OK, common.Error("Group not found"))
		return
	}
	if err := server.LoadNewestDevice(newUUID, ruleEngine); err != nil {
		c.JSON(common.HTTP_OK, common.OkWithMsg(err.Error()))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

// 更新设备
func UpdateDevice(c *gin.Context, ruleEngine typex.RuleX) {

	form := DeviceVo{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	configJson, err := json.Marshal(form.Config)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if form.UUID == "" {
		c.JSON(common.HTTP_OK, common.Error("missing 'uuid' fields"))
		return
	}
	// 更新的时候从数据库往外面拿
	Device, err := service.GetMDeviceWithUUID(form.UUID)
	if err != nil {
		c.JSON(common.HTTP_OK, err)
		return
	}
	MDevice := model.MDevice{
		Type:        form.Type,
		Name:        form.Name,
		Description: form.Description,
		Config:      string(configJson),
	}
	if err := service.UpdateDevice(Device.UUID, &MDevice); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 取消绑定分组,删除原来旧的分组
	Group := service.GetVisualGroup(Device.UUID)
	if err := service.UnBindResource(Group.UUID, Device.UUID); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 重新绑定分组
	if err := service.BindResource(form.Gid, Device.UUID); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err := server.LoadNewestDevice(form.UUID, ruleEngine); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	c.JSON(common.HTTP_OK, common.Ok())
}

// ModbusPoints 获取modbus_excel类型的点位数据
func ModbusPoints(c *gin.Context, ruleEngine typex.RuleX) {
	deviceUuid := c.GetString("deviceUuid")
	list, err := service.AllModbusPointByDeviceUuid(deviceUuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(list))
}

// UpdateModbusPoint 更新modbus_excel类型的点位数据
func UpdateModbusPoint(c *gin.Context, ruleEngine typex.RuleX) {
	type Form struct {
		Id           uint
		DeviceUuid   string `json:"deviceUuid"    gorm:"not null"`
		Tag          string `json:"tag"           gorm:"not null"`
		Function     int    `json:"function"      gorm:"not null"`
		SlaverId     byte   `json:"slaverId"      gorm:"not null"`
		StartAddress uint16 `json:"startAddress"  gorm:"not null"`
		Quality      uint16 `json:"quality"       gorm:"not null"`
	}

	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	err := service.UpdateModbusPoint(model.MModbusPointPosition{
		RulexModel: model.RulexModel{
			ID: form.Id,
		},
		DeviceUuid:   form.DeviceUuid,
		Tag:          form.Tag,
		Function:     form.Function,
		SlaverId:     form.SlaverId,
		StartAddress: form.StartAddress,
		Quality:      form.Quality,
	})

	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	c.JSON(common.HTTP_OK, common.Ok())

}

// ModbusSheetImport 上传Excel文件
func ModbusSheetImport(c *gin.Context, ruleEngine typex.RuleX) {
	// 解析 multipart/form-data 类型的请求体
	err := c.Request.ParseMultipartForm(32 << 20) // 限制上传文件大小为 512MB
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

	deviceUuid := c.Request.Form.Get("deviceUuid")

	// 检查文件类型是否为 Excel
	contentType := header.Header.Get("Content-Type")
	if contentType != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" &&
		contentType != "application/vnd.ms-excel" {
		c.JSON(common.HTTP_OK, common.Error("上传的文件必须是 Excel 格式"))
		return
	}

	// 判断文件大小是否符合要求（1MB）
	if header.Size > 1024*1024 {
		c.JSON(common.HTTP_OK, common.Error("Excel file size cannot be greater than 1MB"))
		return
	}

	list, err := parseModbusPointExcel(file, "Sheet1", deviceUuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	err = service.InsertModbusPointPosition(list)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

func parseModbusPointExcel(r io.Reader,
	sheetName string,
	deviceUuid string) (list []model.MModbusPointPosition, err error) {
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
	// |Tag|Function|SlaverId|StartAddress|Quality|
	if rows[0][0] != "Tag" || rows[0][1] != "Function" || rows[0][2] != "SlaverId" || rows[0][3] != "StartAddress" || rows[0][4] != "Quality" {
		return nil, errors.New("表头不符合要求")
	}

	list = make([]model.MModbusPointPosition, 0)

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		function, _ := strconv.Atoi(row[1])
		slaverId, _ := strconv.ParseInt(row[2], 10, 8)
		address, _ := strconv.ParseUint(row[3], 10, 16)
		quantity, _ := strconv.ParseUint(row[3], 10, 16)
		model := model.MModbusPointPosition{
			DeviceUuid:   deviceUuid,
			Tag:          row[0],
			Function:     function,
			SlaverId:     byte(slaverId),
			StartAddress: uint16(address),
			Quality:      uint16(quantity),
		}
		list = append(list, model)
	}
	return list, nil
}
