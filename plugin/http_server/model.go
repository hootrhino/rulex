package httpserver

import (
	"gorm.io/gorm"
)

type MRule struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	From        string `gorm:"not null"`
	Actions     string `gorm:"not null"`
	Success     string `gorm:"not null"`
	Failed      string `gorm:"not null"`
}

type MInEnd struct {
	gorm.Model
	// UUID for origin resource ID
	UUID        string `gorm:"not null"`
	Type        string `gorm:"not null"`
	Name        string `gorm:"not null"`
	Description string
	Config      string
}

type MOutEnd struct {
	gorm.Model
	// UUID for origin resource ID
	UUID        string `gorm:"not null"`
	Type        string `gorm:"not null"`
	Name        string `gorm:"not null"`
	Description string
	Config      string
}

type MUser struct {
	gorm.Model
	Role        string `gorm:"not null"`
	Username    string `gorm:"not null"`
	Password    string `gorm:"not null"`
	Description string
}

type MLock struct {
	Name     string `gorm:"not null"`
	InitLock int    `gorm:"not null"`
}
