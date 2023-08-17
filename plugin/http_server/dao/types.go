package dao

import "gorm.io/gorm"

/*
*
* DAO 接口
*
 */
type DAO interface {
	Init(string) error
	InitializeData()
	RegisterModel(dst ...interface{})
	Name() string
	DB() *gorm.DB
	Stop()
}
