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

package target

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteTarget struct {
	typex.XStatus
	mainConfig common.SqliteConfig
	status     typex.SourceState
	db         *gorm.DB
}

func NewSqliteTarget(e typex.RuleX) typex.XTarget {
	ht := new(SqliteTarget)
	ht.RuleEngine = e
	ht.mainConfig = common.SqliteConfig{}
	ht.status = typex.SOURCE_DOWN
	return ht
}

func (sqt *SqliteTarget) Init(outEndId string, configMap map[string]interface{}) error {
	sqt.PointId = outEndId
	//
	if err := utils.BindSourceConfig(configMap, &sqt.mainConfig); err != nil {
		return err
	}
	return nil

}
func (sqt *SqliteTarget) Start(cctx typex.CCTX) error {
	sqt.Ctx = cctx.Ctx
	sqt.CancelCTX = cctx.CancelCTX
	//
	db, err := gorm.Open(sqlite.Open(sqt.mainConfig.DbName), &gorm.Config{})
	if err != nil {
		return err
	}
	// 检查是否存在该表
	if err1 := db.Raw(sqt.mainConfig.CreateTbSql).Error; err1 != nil {
		return err1
	}

	sqt.db = db
	//
	sqt.status = typex.SOURCE_UP
	glogger.GLogger.Info("SqliteTarget started")
	return nil
}

func (sqt *SqliteTarget) Status() typex.SourceState {
	return typex.SOURCE_UP

}

/*
*
* 数据转存SQLITE
*
 */
func (sqt *SqliteTarget) To(data interface{}) (interface{}, error) {

	if sqt.db != nil {
		// data 必须是个列表: [1, 2, 3]
		// insert into db1.tb1 value(?, ?, ?), [1, 2, 3]
		if reflect.TypeOf(data).Kind() == reflect.Slice {
			sql := sqt.mainConfig.InsertSql
			for _, v := range data.([]interface{}) {
				sql = strings.Replace(sqt.mainConfig.InsertSql, "?",
					fmt.Sprintf("%v", v), 1)
			}
			if err1 := sqt.db.Raw(sql).Error; err1 != nil {
				return nil, err1
			}
		}
	}
	return nil, fmt.Errorf("sqlite target database error")
}

func (sqt *SqliteTarget) Stop() {
	sqt.status = typex.SOURCE_STOP
	sqt.CancelCTX()
	sqt.db = nil
}
func (sqt *SqliteTarget) Details() *typex.OutEnd {
	return sqt.RuleEngine.GetOutEnd(sqt.PointId)
}
