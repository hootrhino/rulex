package dao

import "gorm.io/gorm"

/*
*
* DAO 接口
*
 */
type DAO interface {
	Init(string) error
	RegisterModel()
	Name() string
	DB() *gorm.DB
	Stop()
}
