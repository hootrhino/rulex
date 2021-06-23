package test

import (
	"testing"

	"github.com/marianogappa/sqlparser"
)

func TestParse(t *testing.T) {
	query, err := sqlparser.Parse("SELECT a, b, c FROM 'data'  WHERE k = '1' AND f > '2'")
	if err != nil {
		t.Log(err)
	}
	t.Log(query.Type)
	t.Log(query.Fields)
	t.Log(query.TableName)
	t.Logf("%#v", query.Conditions)
}
