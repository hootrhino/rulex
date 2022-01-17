package test

import (
	"encoding/json"
	"testing"
)

func Test_gen_td_config(t *testing.T) {
	type tdEngineConfig struct {
		Fqdn           string `json:"fqdn" validate:"required"`
		Port           int    `json:"port" validate:"required"`
		Username       string `json:"username" validate:"required"`
		Password       string `json:"password" validate:"required"`
		DbName         string `json:"dbName" validate:"required"`
		CerateTableSql string `json:"cerateTableSql" validate:"required"`
		InsertSql      string `json:"insertSql" validate:"required"`
	}
	td := tdEngineConfig{
		Fqdn:           "127.0.0.1",
		Port:           4400,
		Username:       "root",
		Password:       "taosdata",
		DbName:         "test",
		CerateTableSql: "CREATE TABLE IF NOT EXISTS test (ts TIMESTAMP, data JSON);",
		InsertSql:      `INSERT INTO d21001 USING meters TAGS ('Beijing.Chaoyang', 2) VALUES ('2021-07-13 14:06:32.272', 10.2, 219, 0.32);`,
	}
	b, _ := json.Marshal(td)
	t.Log(string(b))
}
