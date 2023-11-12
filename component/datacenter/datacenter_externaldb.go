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
	"context"

	"github.com/hootrhino/rulex/component/trailer"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
*
* 外部扩展数据库
*
 */
type ExternalDb struct {
	rulex typex.RuleX
}

func InitExternalDb(rulex typex.RuleX) DataSource {
	Edb := new(ExternalDb)
	Edb.rulex = rulex
	return Edb
}
func (db *ExternalDb) Init() error {
	return nil
}
func (db *ExternalDb) Name() string {
	return "EXTERNAL_DATACENTER"
}
func (db *ExternalDb) GetSchemaDetail(goodsId string) SchemaDetail {
	return SchemaDetail{
		UUID:        "EXTERNAL_DATACENTER",
		SchemaType:  "EXTERNAL_DATACENTER",
		Name:        "外部数据中心",
		LocalPath:   ".EXTERNAL_DATACENTER",
		NetAddr:     ".EXTERNAL_DATACENTER",
		CreateTs:    0,
		Size:        0,
		StorePath:   ".EXTERNAL_DATACENTER",
		Description: "外部数据中心",
	}
}

/*
*
* 去调用RPC的Query
*
 */
func (db *ExternalDb) Query(goodsId, query string) ([]map[string]any, error) {
	var err error
	Rows := []map[string]any{}
	if goodsPs := trailer.Get(goodsId); goodsPs != nil {
		grpcConnection, err1 := grpc.Dial(goodsPs.Info.NetAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err1 != nil {
			glogger.GLogger.Error(err1)
			err = err1
			return Rows, err
		}
		defer grpcConnection.Close()
		client := trailer.NewTrailerClient(grpcConnection)
		columns, err2 := client.Query(context.Background(), &trailer.DataRowsRequest{
			Query: []byte(query),
		})
		if err2 != nil {
			glogger.GLogger.Error(err2)
			err = err2
			return Rows, err
		}
		for _, row := range columns.Row {
			Row := map[string]any{}
			for _, column := range row.Column {
				Row[string(column.GetName())] = covertGoTypeToJsType(column)
			}
			Rows = append(Rows, Row)
		}
	}
	return Rows, err
}
