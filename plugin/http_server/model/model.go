package model

import (
	"database/sql/driver"
	"time"

	"gopkg.in/square/go-jose.v2/json"
)

type RulexModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
}
type StringList []string

/*
*
* 给GORM用的
*
 */
func (f StringList) Value() (driver.Value, error) {
	b, err := json.Marshal(f)
	return string(b), err
}

/*
*
* 给GORM用的
*
 */
func (f *StringList) Scan(data interface{}) error {
	return json.Unmarshal([]byte(data.(string)), f)
}

func (f StringList) String() string {
	b, _ := json.Marshal(f)
	return string(b)
}
func (f StringList) Len() int {
	return len(f)
}

type MRule struct {
	RulexModel
	UUID        string     `gorm:"not null"`
	Name        string     `gorm:"not null"`
	Type        string     // 脚本类型，目前支持"lua"
	FromSource  StringList `gorm:"not null type:string[]"`
	FromDevice  StringList `gorm:"not null type:string[]"`
	Actions     string     `gorm:"not null"`
	Success     string     `gorm:"not null"`
	Failed      string     `gorm:"not null"`
	Description string
}

type MInEnd struct {
	RulexModel
	// UUID for origin source ID
	UUID      string     `gorm:"not null"`
	Type      string     `gorm:"not null"`
	Name      string     `gorm:"not null"`
	BindRules StringList `json:"bindRules"` // 与之关联的规则表["A","B","C"]

	Description string
	Config      string
	XDataModels string
}

func (md MInEnd) GetConfig() map[string]interface{} {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(md.Config), &result)
	if err != nil {
		return map[string]interface{}{}
	}
	return result
}

type MOutEnd struct {
	RulexModel
	// UUID for origin source ID
	UUID        string `gorm:"not null"`
	Type        string `gorm:"not null"`
	Name        string `gorm:"not null"`
	Description string
	Config      string
}

func (md MOutEnd) GetConfig() map[string]interface{} {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(md.Config), &result)
	if err != nil {
		return map[string]interface{}{}
	}
	return result
}

type MUser struct {
	RulexModel
	Role        string `gorm:"not null"`
	Username    string `gorm:"not null"`
	Password    string `gorm:"not null"`
	Description string
}

// 设备元数据
type MDevice struct {
	RulexModel
	UUID        string `gorm:"not null"`
	Name        string `gorm:"not null"`
	Type        string `gorm:"not null"`
	Config      string
	BindRules   StringList `json:"bindRules"` // 与之关联的规则表["A","B","C"]
	Description string
}

func (md MDevice) GetConfig() map[string]interface{} {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(md.Config), &result)
	if err != nil {
		return map[string]interface{}{}
	}
	return result
}

//
// 外挂
//

type MGoods struct {
	RulexModel
	UUID        string `gorm:"not null"`
	LocalPath   string `gorm:"not null"`
	GoodsType   string `gorm:"not null"` // LOCAL, EXTERNAL
	ExecuteType string `gorm:"not null"` // exe,elf,js,py....
	AutoStart   *bool  `gorm:"not null"`
	NetAddr     string `gorm:"not null"`
	Args        string `gorm:"not null"`
	Description string `gorm:"not null"`
}

/*
*
* LUA应用
*
 */
type MApp struct {
	RulexModel
	UUID        string `gorm:"not null"` // 名称
	Name        string `gorm:"not null"` // 名称
	Version     string `gorm:"not null"` // 版本号
	AutoStart   *bool  `gorm:"not null"` // 允许启动
	LuaSource   string `gorm:"not null"` // LuaSource
	Description string `gorm:"not null"` // 文件路径, 是相对于main的apps目录
}

/*
*
* 用户上传的AI数据[v0.5.0]
*
 */
type MAiBase struct {
	RulexModel
	UUID        string `gorm:"not null"` // 名称
	Name        string `gorm:"not null"` // 名称
	Type        string `gorm:"not null"` // 类型
	IsBuildIn   bool   `gorm:"not null"` // 是否内建
	Version     string `gorm:"not null"` // 版本号
	Filepath    string `gorm:"not null"` // 文件路径, 是相对于main的apps目录
	Description string `gorm:"not null"`
}

// MModbusPointPosition modbus数据点位
type MModbusPointPosition struct {
	RulexModel
	DeviceUuid   string `json:"deviceUuid"    gorm:"not null"`
	Tag          string `json:"tag"           gorm:"not null"`
	Function     int    `json:"function"      gorm:"not null"`
	SlaverId     byte   `json:"slaverId"      gorm:"not null"`
	StartAddress uint16 `json:"startAddress"  gorm:"not null"`
	Quality      uint16 `json:"quality"       gorm:"not null"`
}

//--------------------------------------------------------------------------------------------------
// 0.6.0
//--------------------------------------------------------------------------------------------------
/*
*
* 大屏
*
 */
type MVisual struct {
	RulexModel
	UUID      string `gorm:"not null"` // 名称
	Name      string `gorm:"not null"` // 名称
	Type      string `gorm:"not null"` // 类型
	Status    bool   `gorm:"not null"` // 状态, EDITING, PUBLISH
	Content   string `gorm:"not null"` // 大屏的内容
	Thumbnail string `gorm:"not null"` // 缩略图
}

/*
*
* 通用分组
*
 */
type MGenericGroup struct {
	RulexModel
	UUID   string `gorm:"not null"`
	Name   string `gorm:"not null"` // 名称
	Type   string `gorm:"not null"` // 组的类型, DEVICE: 设备分组
	Parent string `gorm:"not null"` // 上级, 如果是0表示根节点
}

/*
*
* 分组表的绑定关系表
*
 */
type MGenericGroupRelation struct {
	RulexModel
	UUID string `gorm:"not null"`
	Gid  string `gorm:"not null"` // 分组ID
	Rid  string `gorm:"not null"` // 被绑定方
}

/*
*
* 应用商店里面的各类协议脚本
*
 */
type MProtocolApp struct {
	RulexModel
	UUID    string `gorm:"not null"` // 名称
	Name    string `gorm:"not null"` // 名称
	Type    string `gorm:"not null"` // 类型: IN OUT DEVICE APP
	Content string `gorm:"not null"` // 协议包的内容
}

/*
*
* 系统配置参数, 直接以String保存，完了以后再加工成Dto结构体
*
 */
type MNetworkConfig struct {
	RulexModel
	Type        string     `gorm:"not null"` // 类型: ubuntu16, ubuntu18
	Interface   string     `gorm:"not null"` // eth1 eth0
	Address     string     `gorm:"not null"`
	Netmask     string     `gorm:"not null"`
	Gateway     string     `gorm:"not null"`
	DNS         StringList `gorm:"not null"`
	DHCPEnabled *bool      `gorm:"not null"`
}

/*
*
* 无线网络配置
*
 */
type MWifiConfig struct {
	RulexModel
	Interface string `gorm:"not null"`
	SSID      string `gorm:"not null"`
	Password  string `gorm:"not null"`
	Security  string `gorm:"not null"` // wpa2-psk wpa3-psk
}

/**
 * 定时任务
 */
type MCronTask struct {
	RulexModel
	UUID      string    `gorm:"not null; default:''" json:"uuid"`
	Name      string    `gorm:"not null;" json:"name"`
	CronExpr  string    `gorm:"not null" json:"cronExpr"` // linux cron standard
	Enable    string    `json:"enable"`                   // 0-disable 1-enable
	TaskType  int       `json:"taskType"`                 // 1-shell 目前
	Command   string    `json:"command"`                  // 目前不使用，默认都是linux shell
	Args      *string   `json:"args"`                     // "-param1 -param2 -param3"
	IsRoot    string    `json:"isRoot"`                   // 0-false 1-true
	WorkDir   string    `json:"workDir"`                  // 目前不使用，默认工作路径和网关工作路径保持一致
	Env       string    `json:"env"`                      // ["A=e1", "B=e2", "C=e3"]
	Script    string    `json:"script"`                   // 脚本内容
	UpdatedAt time.Time `json:"updatedAt"`
}

/**
 * 任务结果
 */
type MCronResult struct {
	RulexModel
	TaskUuid  string    `gorm:"not null; default:''" json:"taskUuid,omitempty"`
	Status    string    `json:"status"`             // 1-running 2-end
	ExitCode  string    `json:"exitCode,omitempty"` // 0-success other-failed
	LogPath   string    `json:"logPath,omitempty"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

type PageRequest struct {
	Current     int `json:"current,omitempty"`
	Size        int `json:"size,omitempty"`
	SearchCount int `json:"searchCount,omitempty"`
}

type PageResult struct {
	Current int `json:"current"`
	Size    int `json:"size"`
	Total   int `json:"total"`
	Records any `json:"records"`
}
