package httpserver

import (
	"errors"
	"io"
	"strconv"

	"github.com/xuri/excelize/v2"

	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2/json"
)

/*
*
* 列表先读数据库，然后读内存，合并状态后输出
*
 */
func DeviceDetail(c *gin.Context, hs *HttpApiServer) {
	uuid, _ := c.GetQuery("uuid")
	mdev, err := hs.GetDeviceWithUUID(uuid)
	if err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	device := hs.ruleEngine.GetDevice(mdev.UUID)
	if device == nil {
		// 如果内存里面没有就给安排一个死设备
		tDevice := new(typex.Device)
		tDevice.UUID = mdev.UUID
		tDevice.Name = mdev.Name
		tDevice.Type = typex.DeviceType(mdev.Type)
		tDevice.Description = mdev.Description
		tDevice.BindRules = map[string]typex.Rule{}
		tDevice.Config = mdev.GetConfig()
		tDevice.State = typex.DEV_STOP
		c.JSON(HTTP_OK, OkWithData(tDevice))
		return
	}
	device.State = device.Device.Status()
	c.JSON(HTTP_OK, OkWithData(device))
}
func Devices(c *gin.Context, hs *HttpApiServer) {
	uuid, _ := c.GetQuery("uuid")
	if uuid == "" {
		devices := []*typex.Device{}
		for _, mdev := range hs.AllDevices() {
			device := hs.ruleEngine.GetDevice(mdev.UUID)
			if device == nil {
				tDevice := new(typex.Device)
				tDevice.UUID = mdev.UUID
				tDevice.Name = mdev.Name
				tDevice.Type = typex.DeviceType(mdev.Type)
				tDevice.Description = mdev.Description
				tDevice.BindRules = map[string]typex.Rule{}
				tDevice.Config = map[string]interface{}{}
				tDevice.State = typex.DEV_STOP
				devices = append(devices, tDevice)
			}
			if device != nil {
				device.State = device.Device.Status()
				devices = append(devices, device)
			}
		}
		c.JSON(HTTP_OK, OkWithData(devices))
		return
	}
	mdev, err := hs.GetDeviceWithUUID(uuid)
	if err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	device := hs.ruleEngine.GetDevice(mdev.UUID)
	if device == nil {
		// 如果内存里面没有就给安排一个死设备
		tDevice := new(typex.Device)
		tDevice.UUID = mdev.UUID
		tDevice.Name = mdev.Name
		tDevice.Type = typex.DeviceType(mdev.Type)
		tDevice.Description = mdev.Description
		tDevice.BindRules = map[string]typex.Rule{}
		tDevice.Config = mdev.GetConfig()
		tDevice.State = typex.DEV_STOP
		c.JSON(HTTP_OK, OkWithData(tDevice))
		return
	}
	device.State = device.Device.Status()
	c.JSON(HTTP_OK, OkWithData(device))
}

// 删除设备
func DeleteDevice(c *gin.Context, hs *HttpApiServer) {
	uuid, _ := c.GetQuery("uuid")
	Mdev, err := hs.GetDeviceWithUUID(uuid)
	if err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	// 要处理这个空字符串 ""
	if Mdev.BindRules.Len() == 1 && len(Mdev.BindRules[0]) != 0 {
		c.JSON(HTTP_OK, Error("Can't remove, Already have rule bind:"+Mdev.BindRules.String()))
		return
	}
	// 检查是否有规则被绑定了
	for _, ruleId := range Mdev.BindRules {
		if ruleId != "" {
			_, err0 := hs.GetMRuleWithUUID(ruleId)
			if err0 != nil {
				c.JSON(HTTP_OK, Error400(err0))
				return
			}
		}

	}
	if err := hs.DeleteDevice(uuid); err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	old := hs.ruleEngine.GetDevice(uuid)
	if old != nil {
		if old.Device.Status() == typex.DEV_UP {
			old.Device.Details().State = typex.DEV_STOP
			old.Device.Stop()
		}
	}
	hs.ruleEngine.RemoveDevice(uuid)
	c.JSON(HTTP_OK, Ok())

}

// 创建设备
func CreateDevice(c *gin.Context, hs *HttpApiServer) {
	type Form struct {
		UUID         string                 `json:"uuid"`
		Name         string                 `json:"name"`
		Type         string                 `json:"type"`
		ActionScript string                 `json:"actionScript"`
		Config       map[string]interface{} `json:"config"`
		Description  string                 `json:"description"`
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	configJson, err := json.Marshal(form.Config)
	if err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	newUUID := utils.DeviceUuid()
	if err := hs.InsertDevice(&MDevice{
		UUID:        newUUID,
		Type:        form.Type,
		Name:        form.Name,
		Description: form.Description,
		Config:      string(configJson),
		BindRules:   []string{},
	}); err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	if err := hs.LoadNewestDevice(newUUID); err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	c.JSON(HTTP_OK, Ok())

}

// 更新设备
func UpdateDevice(c *gin.Context, hs *HttpApiServer) {
	type Form struct {
		UUID        string                 `json:"uuid"`
		Name        string                 `json:"name"`
		Type        string                 `json:"type"`
		Config      map[string]interface{} `json:"config"`
		Description string                 `json:"description"`
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	configJson, err := json.Marshal(form.Config)
	if err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	if form.UUID == "" {
		c.JSON(HTTP_OK, Error("missing 'uuid' fields"))
		return
	}
	// 更新的时候从数据库往外面拿
	Device, err := hs.GetDeviceWithUUID(form.UUID)
	if err != nil {
		c.JSON(HTTP_OK, err)
		return
	}

	if err := hs.UpdateDevice(Device.UUID, &MDevice{
		Type:        form.Type,
		Name:        form.Name,
		Description: form.Description,
		Config:      string(configJson),
	}); err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}

	if err := hs.LoadNewestDevice(form.UUID); err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}

	c.JSON(HTTP_OK, Ok())
}

// ModbusSheetImport 上传Excel文件
func ModbusSheetImport(c *gin.Context, hs *HttpApiServer) {
	// 解析 multipart/form-data 类型的请求体
	err := c.Request.ParseMultipartForm(32 << 20) // 限制上传文件大小为 512MB
	if err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	defer file.Close()

	deviceUuid := c.Request.Form.Get("deviceUuid")

	// 检查文件类型是否为 Excel
	contentType := header.Header.Get("Content-Type")
	if contentType != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" &&
		contentType != "application/vnd.ms-excel" {
		c.JSON(HTTP_OK, Error("上传的文件必须是 Excel 格式"))
		return
	}

	list, err := parseModbusPointExcel(file, "Sheet1", deviceUuid)
	if err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}

	err = hs.InsertModbusPointPosition(list)
	if err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	c.JSON(HTTP_OK, Ok())
}

func parseModbusPointExcel(r io.Reader, sheetName string, deviceUuid string) (list []*MModbusPointPosition, err error) {
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

	list = make([]*MModbusPointPosition, 0)

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		function, _ := strconv.Atoi(row[1])
		slaverId, _ := strconv.ParseInt(row[2], 10, 8)
		address, _ := strconv.ParseUint(row[3], 10, 16)
		quantity, _ := strconv.ParseUint(row[3], 10, 16)
		model := &MModbusPointPosition{
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
