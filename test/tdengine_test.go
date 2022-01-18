package test

import (
	"encoding/json"
	"rulex/core"
	"rulex/typex"
	"testing"
)

type tdEngineConfig struct {
	Fqdn           string `json:"fqdn" validate:"required"`
	Port           int    `json:"port" validate:"required"`
	Username       string `json:"username" validate:"required"`
	Password       string `json:"password" validate:"required"`
	DbName         string `json:"dbName" validate:"required"`
	CreateDbSql    string `json:"createDbSql" validate:"required"`
	CreateTableSql string `json:"createTableSql" validate:"required"`
	InsertSql      string `json:"insertSql" validate:"required"`
}

func Test_gen_td_config(t *testing.T) {
	td := tdEngineConfig{
		Fqdn:           "127.0.0.1",
		Port:           4400,
		Username:       "root",
		Password:       "taosdata",
		DbName:         "test",
		CreateDbSql:    "CREATE DATABASE IF NOT EXISTS device UPDATE 0;",
		CreateTableSql: "CREATE TABLE IF NOT EXISTS meter (ts TIMESTAMP, current FLOAT, valtage FLOAT);",
		InsertSql:      `INSERT INTO meter VALUES (NOW, %v, %v);`,
	}
	b, _ := json.Marshal(td)
	t.Log(string(b))
}
func Test_gen_tdEngineConfig(t *testing.T) {
	c, err := core.RenderOutConfig(typex.TDENGINE_TARGET, "TDENGINE", tdEngineConfig{})
	if err != nil {
		t.Error(err)
	}
	b, _ := json.MarshalIndent(c.Views, "  ", "")
	t.Log(string(b))
}
