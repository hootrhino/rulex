package model

import (
	"database/sql/driver"
	"gopkg.in/square/go-jose.v2/json"
	"time"
)

type BaseModel struct {
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

func (f stringList) String() string {
	b, _ := json.Marshal(f)
	return string(b)
}
func (f stringList) Len() int {
	return len(f)
}
