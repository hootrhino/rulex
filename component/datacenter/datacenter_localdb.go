// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package datacenter

import (
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/sqlparser"
)

/*
*
* 本地数据库是Sqlite，用来存储比如Modbus等数据
*
 */
type LocalDb struct {
	Sqlite *SqliteDAO
	rulex  typex.RuleX
}

func InitLocalDb(rulex typex.RuleX) DataSource {
	return &LocalDb{
		Sqlite: InitSqliteDAO(rulex),
		rulex:  rulex,
	}
}

func (ldb *LocalDb) Init() error {
	return nil
}

func (ldb *LocalDb) GetSchemaDetail(goodsId string) SchemaDetail {
	return SchemaDetail{
		UUID:        "INTERNAL_DATACENTER",
		SchemaType:  "INTERNAL_DATACENTER",
		Name:        "RULEX内置轻量级数据仓库",
		LocalPath:   ".local",
		NetAddr:     ".local",
		CreateTs:    0,
		Size:        0,
		StorePath:   ".local",
		Description: "本地内部数据中心",
	}
}

/*
*
  - 此处执行SQL
    // 第一行数据
    // Row1 := map[string]any{
    // 	"Key1": 1,
    // 	"Key2": 1,
    // 	"Key3": 1,
    // 	"Key4": 1,
    // 	"Key5": 1,
    // 	"Key6": 1,
    // }
    // // 第二行数据
    // Row2 := map[string]any{
    // 	"Key1": 1,
    // 	"Key2": 1,
    // 	"Key3": 1,
    // 	"Key4": 1,
    // 	"Key5": 1,
    // 	"Key6": 1,
    // }
    // Rows = append(Rows, Row1)
    // Rows = append(Rows, Row2)
*/
const (
	// UnknownType is the zero value for a Type
	UnknownType = 0
	// Select represents a SELECT query
	Select = 1
	// Update represents an UPDATE query
	Update = 2
	// Insert represents an INSERT query
	Insert = 3
	// Delete represents a DELETE query
	Delete = 4
)

func (ldb *LocalDb) Query(goodsId, query string) ([]map[string]any, error) {
	result := []map[string]any{}
	Query, err := sqlparser.Parse(query)
	if err != nil {
		return []map[string]any{}, err
	}
	//select
	if Query.Type == Select {
		err1 := ldb.Sqlite.db.Raw(query).Scan(&result).Error
		if err1 != nil {
			return []map[string]any{}, err1
		}
		return result, nil
	}
	return result, ldb.Sqlite.db.Raw(query).Error
}
