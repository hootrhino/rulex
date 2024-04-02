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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hootrhino/rulex/component/iotschema"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

/*
*
* 点位表
*
 */
type kndDataPoint struct {
	DeviceUUID    string `json:"device_uuid"`
	Name          string `json:"name"`
	Alias         string `json:"alias"`
	Function      *int   `json:"function"`
	Group         *int   `json:"group"`
	Status        int    `json:"status"`        // 运行时数据
	LastFetchTime uint64 `json:"lastFetchTime"` // 运行时数据
	Value         string `json:"value"`         // 运行时数据
}

//	{
//		"error": -1, // 错误码
//		"error-message": "未知错误" // 错误的消息
//	}
//
// KND 部分 CNC 系统上运行了 REST API 服务器，用于向第三方开放部分数据接口。
// 服务运行于 HTTP 标准端口（即 80 端口），当前最新版本为 v1.3， 建议使用 v1.3， 请求
// 的基地址为/api/v1.3，
// 兼容 v1.2, 若使用 V1.2,请求基地址为/api/v1.2（以下接口文档中若无特殊说明，均支持
// v1.2 和 v1.3）。
// 例如 CNC 的 IP 地址为 192.168.1.101，那么访问 CNC 的/status 接口应使用地址：
// v1.3 版本 http://192.168.1.101/api/v1.3/status
// v1.2 版本 http://192.168.1.101/api/v1.2/status
// 所有接口只接收 Content-Type 是 application/json 类型的 HTTP 数据，
// 除部分文件相关的接口外，大部分接口返回的也是 application/json 类型的 HTTP 数据。
// 所有数据的编码必须是 UTF-8
// 以下说明中，如未加特殊说明，均表示 HTTP 方法为 GET

type kdnCNCStatus struct {
	// 当前运行状态（运行、停止等）
	// - 0：CNC 处于停止状态
	// - 1：CNC 处于暂停（进给保持）状态
	// - 2：CNC 处于运行状态
	RunStatus int `json:"run-status"`
	// 当前工作模式（录入，自动等）
	// - 0：录入方式
	// - 1：自动方式
	// - 3：编辑方式
	// - 4：单步方式
	// - 5：手动方式
	// - 6：手动编辑（示教）方式
	// - 7：手轮编辑（示教）方式
	// - 8：手轮方式
	// - 9：（机械）回零方式
	// - 10：程序回零方式
	OprMode int  `json:"opr-mode"`
	Ready   bool `json:"ready"` // 是否准备就绪
	// 准备未绪的原因掩码值，可为以下值的 or:
	// - 0x1：急停信号有效
	// - 0x2：伺服准备未绪
	// - 0x4：IO 准备未绪（远程 IO 设备等）
	NotReadyReason int      `json:"not-ready-reason"`
	Alarms         []string `json:"alarms"`         // 当前报警类型的列表，如果当前没有报警，则列表为空
	MachineLock    bool     `json:"machine-lock"`   // 锁定状态
	AuxiliaryLock  bool     `json:"auxiliary-lock"` // 辅助锁状态
	DryRun         bool     `json:"dry-run"`        // 空运行状态
	SingleBlock    bool     `json:"single-block"`   // 单段状态
	OptionalSkip   bool     `json:"optional-skip"`  // 跳段状态
	OptionalStop   bool     `json:"optional-stop"`  // 选择停状态
}
type kdnCNCInfo struct {
	ID              int      `json:"id"`               // 唯一 ID 的 64 位十进制表示
	Type            string   `json:"type"`             // 系统类型
	Manufacturer    string   `json:"manufacturer"`     // 制造商
	ManufactureTime string   `json:"manufacture-time"` // 生产时间
	CncType         string   `json:"cnc-type"`         // 车铣类型
	CncName         string   `json:"cnc-name"`         // 网络参数中的机床名称
	SoftVersion     string   `json:"soft-version"`     // 系统软件版本号
	FpgaVersion     string   `json:"fpga-version"`     // FPGA 版本号
	LadderVersion   string   `json:"ladder-version"`   // 梯图版本号
	NcAxes          []string `json:"nc-axes"`          // 用户 NC 轴列表
	NcRelativeAxes  []string `json:"nc-relative-axes"` // 用户 NC 轴相对坐标地址列表
	Axes            []string `json:"axes"`             // 用户轴列表
	RelativeAxes    []string `json:"relative-axes"`    // 用户轴相对坐标地址列表
}

type kdn_cnc_config struct {
	Host          string         `json:"host" validate:"required"`       // IP:Port
	ApiVersion    int            `json:"apiVersion" validate:"required"` // API 版本,2 | 3
	CNCInfo       kdnCNCInfo     `json:"cncInfo"`
	CNCStatus     kdnCNCStatus   `json:"cncStatus"`
	KndDataPoints []kndDataPoint `json:"kndDataPoints"`
}

/*
*
* 凯帝恩CNC
*
 */
type KDN_CNC struct {
	typex.XStatus
	status     typex.DeviceState
	mainConfig kdn_cnc_config
}

func NewKDN_CNC(e typex.RuleX) typex.XDevice {
	hd := new(KDN_CNC)
	hd.RuleEngine = e
	return hd
}

//  初始化
func (hd *KDN_CNC) Init(devId string, configMap map[string]interface{}) error {
	hd.PointId = devId
	if err := utils.BindSourceConfig(configMap, &hd.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	return nil
}

// 启动
func (hd *KDN_CNC) Start(cctx typex.CCTX) error {
	hd.Ctx = cctx.Ctx
	hd.CancelCTX = cctx.CancelCTX
	kdnCNCInfo, err1 := hd.getCncInfo()
	if err1 != nil {
		glogger.GLogger.Error("Fetch KDN CNC Status Failed", err1,
			" Check your device or configuration")
		hd.status = typex.DEV_DOWN
		return err1
	}
	hd.mainConfig.CNCInfo = kdnCNCInfo
	glogger.GLogger.Info("Fetch KDN CNC Status Success:", kdnCNCInfo.ID)
	hd.status = typex.DEV_UP
	return nil
}

func (hd *KDN_CNC) OnRead(cmd []byte, data []byte) (int, error) {
	return 0, nil
}

// 把数据写入设备
func (hd *KDN_CNC) OnWrite(cmd []byte, b []byte) (int, error) {
	return 0, nil
}

// 设备当前状态
// 主要判别依据是获取CNC的基础信息
func (hd *KDN_CNC) Status() typex.DeviceState {
	_, err1 := hd.getCncInfo()
	if err1 != nil {
		glogger.GLogger.Error("Fetch KDN CNC Status Failed", err1,
			" Check your device or configuration")
		hd.status = typex.DEV_DOWN
	}
	return hd.status
}

// 停止设备
func (hd *KDN_CNC) Stop() {
	hd.status = typex.DEV_DOWN
	hd.CancelCTX()
}

// 设备属性，是一系列属性描述
func (hd *KDN_CNC) Property() []iotschema.IoTSchema {
	return []iotschema.IoTSchema{}
}

// 真实设备
func (hd *KDN_CNC) Details() *typex.Device {
	return hd.RuleEngine.GetDevice(hd.PointId)
}

// 状态
func (hd *KDN_CNC) SetState(status typex.DeviceState) {
	hd.status = status

}

// 驱动
func (hd *KDN_CNC) Driver() typex.XExternalDriver {
	return nil
}

// --------------------------------------------------------------------------------------------------
//
// --------------------------------------------------------------------------------------------------

func (hd *KDN_CNC) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
func (hd *KDN_CNC) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}

/*
*
* HTTP GET
*
 */
func kdnHttpGet(url string) (string, error) {
	var err error
	client := http.DefaultClient
	client.Timeout = 2 * time.Second
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		glogger.GLogger.Warn(err)
		return "", err
	}
	response, err := client.Do(request)
	if err != nil {
		glogger.GLogger.Warn(err)
		return "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		glogger.GLogger.Warn(err)
		return "", err
	}
	return string(body), nil
}

// ==================================================================================================
// API 接口
// ==================================================================================================
/*
*
* 获取CNC的基本信息
*
 */
func (hd *KDN_CNC) getCncInfo() (kdnCNCInfo, error) {
	ApiVersion := "v1.3" // 默认就是1.3
	if hd.mainConfig.ApiVersion == 2 {
		ApiVersion = "v1.2"
	}
	kdnCNCInfo := kdnCNCInfo{}
	status, err1 := kdnHttpGet(fmt.Sprintf("http://%s/api/%v/",
		hd.mainConfig.Host, ApiVersion))
	if err1 != nil {
		return kdnCNCInfo, err1
	}
	if err2 := json.Unmarshal([]byte(status), &kdnCNCInfo); err2 != nil {
		return kdnCNCInfo, err2
	}
	return kdnCNCInfo, nil
}

/*
*
* 获取运行时状态
*
 */
func (hd *KDN_CNC) GetCncStatus() (kdnCNCStatus, error) {
	ApiVersion := "v1.3" // 默认就是1.3
	if hd.mainConfig.ApiVersion == 2 {
		ApiVersion = "v1.2"
	}
	kdnCNCStatus := kdnCNCStatus{}
	status, err1 := kdnHttpGet(fmt.Sprintf("http://%s/api/%v/status",
		hd.mainConfig.Host, ApiVersion))
	if err1 != nil {
		return kdnCNCStatus, err1
	}
	if err2 := json.Unmarshal([]byte(status), &kdnCNCStatus); err2 != nil {
		return kdnCNCStatus, err2
	}
	return kdnCNCStatus, nil
}
