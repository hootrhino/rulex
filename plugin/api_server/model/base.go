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

// Value 写入数据库之前，对数据做类型转换
func (f stringList) Value() (driver.Value, error) {
	b, err := json.Marshal(f)
	return string(b), err
}

// Scan 将数据库中取出的数据，赋值给目标类型
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
