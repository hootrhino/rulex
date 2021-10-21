package httpserver

import (
	"time"
)

type RulexModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
}
type MRule struct {
	RulexModel
	Name        string `gorm:"not null"`
	Description string
	From        string `gorm:"not null"`
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
