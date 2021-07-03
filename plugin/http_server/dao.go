package plugin

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ngaut/log"
)
// TODO
var localDB *sql.DB

func InitDb() {
	db, err := sql.Open("sqlite3", "./rulex.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	localDB = db
}

//
func createTable(db *sql.DB, sql string) {

}
