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
type fromSource []string
type fromDevice []string

func (f fromSource) Value() (driver.Value, error) {
	b, err := json.Marshal(f)
	return string(b), err
}

func (f *fromSource) Scan(data interface{}) error {
	return json.Unmarshal([]byte(data.(string)), f)
}

func (f fromDevice) Value() (driver.Value, error) {
	b, err := json.Marshal(f)
	return string(b), err
}

func (f *fromDevice) Scan(data interface{}) error {
	return json.Unmarshal([]byte(data.(string)), f)
}

type MRule struct {
	RulexModel
	UUID        string `gorm:"not null"`
	Name        string `gorm:"not null"`
	Description string
	FromSource  fromSource `gorm:"not null type:string[]"`
	FromDevice  fromDevice `gorm:"not null type:string[]"`
	Actions     string     `gorm:"not null"`
	Success     string     `gorm:"not null"`
	Failed      string     `gorm:"not null"`
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
	UUID         string `gorm:"not null"`
	Name         string `gorm:"not null"`
	Type         string `gorm:"not null"`
	ActionScript string
	Config       string
	Description  string
}
