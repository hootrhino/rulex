package httpserver

import (
	"database/sql/driver"
	"time"

	"gopkg.in/square/go-jose.v2/json"
)

type RulexModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
}
type stringList []string

func (f stringList) Value() (driver.Value, error) {
	b, err := json.Marshal(f)
	return string(b), err
}

func (f *stringList) Scan(data interface{}) error {
	return json.Unmarshal([]byte(data.(string)), f)
}

type MRule struct {
	RulexModel
	UUID        string     `gorm:"not null"`
	Name        string     `gorm:"not null"`
	Type        string     // 脚本类型，目前支持"lua"和"expr"两种
	FromSource  stringList `gorm:"not null type:string[]"`
	FromDevice  stringList `gorm:"not null type:string[]"`
	Expression  string     `gorm:"not null"` // Expr脚本
	Actions     string     `gorm:"not null"`
	Success     string     `gorm:"not null"`
	Failed      string     `gorm:"not null"`
	Description string
}

type MInEnd struct {
	RulexModel
	// UUID for origin source ID
	UUID        string `gorm:"not null"`
	Type        string `gorm:"not null"`
	Name        string `gorm:"not null"`
	Description string
	Config      string
	XDataModels string
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
	UUID string `gorm:"not null"`
	Name string `gorm:"not null"`
	Type string `gorm:"not null"`
	//   这个字段本来用于给设备单独新建脚本的，但是目前已经有了规则，所以先留空
	// 或许以后会用到
	ActionScript string
	Config       string
	Description  string
}

//
// 外挂
//

type MGoods struct {
	RulexModel
	UUID        string     `gorm:"not null"`
	Addr        string     `gorm:"not null"`
	Description string     `gorm:"not null"`
	Args        stringList `gorm:"not null"`
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
	Filepath    string `gorm:"not null"` // 文件路径, 是相对于main的apps目录
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
	Version     string `gorm:"not null"` // 版本号
	Filepath    string `gorm:"not null"` // 文件路径, 是相对于main的apps目录
	Description string `gorm:"not null"`
}
