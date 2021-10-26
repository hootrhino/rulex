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
type from []string

func (f from) Value() (driver.Value, error) {
	b, err := json.Marshal(f)
	return string(b), err
}

func (f *from) Scan(data interface{}) error {
	return json.Unmarshal([]byte(data.(string)), f)
}

type MRule struct {
	RulexModel
	Name        string `gorm:"not null"`
	Description string
	From        from   `gorm:"not null type:string[]"`
	Actions     string `gorm:"not null"`
	Success     string `gorm:"not null"`
	Failed      string `gorm:"not null"`
}

type MInEnd struct {
	RulexModel
	// UUID for origin resource ID
	UUID        string `gorm:"not null"`
	Type        string `gorm:"not null"`
	Name        string `gorm:"not null"`
	Description string
	Config      string
}

type MOutEnd struct {
	RulexModel
	// UUID for origin resource ID
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

type MLock struct {
	Name     string `gorm:"not null"`
	InitLock int    `gorm:"not null"`
}
