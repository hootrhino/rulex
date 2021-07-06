package httpserver

import "gorm.io/gorm"

type MRule struct {
	gorm.Model
	Name        string
	Description string
	From        string
	Actions     string
	Success     string
	Failed      string
}

type MInEnd struct {
	gorm.Model
	Type        string
	Name        string
	Description string
	Config      string
}

type MOutEnd struct {
	gorm.Model
	Type        string
	Name        string
	Description string
	Config      string
}

type MUser struct {
	gorm.Model
	Role        string
	Username    string
	Password    string
	Description string
}
