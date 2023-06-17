package test

import (
	"fmt"
	"testing"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
)

// go test -timeout 30s -run ^Test_sqlparser github.com/hootrhino/rulex/test -v -count=1

func Test_sqlparser(t *testing.T) {
	sql := `
	select * from (select * from user_items) AS U where U.user_id=1
	`
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		panic(err)
	}
	fmt.Printf("stmt = %+v\n", stmt)
}
