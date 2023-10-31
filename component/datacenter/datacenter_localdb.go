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

/*
*
* 本地数据库是Sqlite，用来存储比如Modbus等数据
*
 */
type LocalDb struct {
}

func Init(ldb *LocalDb) error {
	return nil
}

func (ldb *LocalDb) Name() string {
	return "LOCALDB"
}
func (ldb *LocalDb) GetSchemaDetail(goodsId string) SchemaDetail {
	return SchemaDetail{}
}
func (ldb *LocalDb) Query(goodsId, query string) ([]map[string]any, error) {
	return []map[string]any{}, nil
}
