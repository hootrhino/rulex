package test

import (
	"rulex/plugin/http_server"
	"testing"

	"github.com/ngaut/log"
)

func TestGorm(t *testing.T) {

	// Create
	// sqliteDb.Create(&MInEnd{
	// 	Type:        "MQTT",
	// 	Name:        "MQTT input stream",
	// 	Description: "OK is good",
	// 	Config:      `{"k":"v"}`,
	// })
	m1 := []httpserver.MInEnd{}
	// httpserver.sqliteDb.Find(&inends)
	// t.Logf("%#v", inends)

	if err := httpserver.SqliteDb.Where("Id=?", 1).First(&m1).Error; err != nil {
		t.Error(err)
	} else {
		log.Debug(m1)
	}
}
