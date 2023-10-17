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
	"time"

	"github.com/hootrhino/rulex/component/trailer"
	"github.com/hootrhino/rulex/glogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var __DefaultDataCenter *DataCenter

/*
*
* 留着未来扩充数据中心的功能
*
 */
type DataCenter struct {
}

func InitDataCenter() {
	__DefaultDataCenter = new(DataCenter)
}

/*
*
* 获取表格定义
*
 */
func SchemaList() []SchemaDetail {
	Schemas := []SchemaDetail{}
	trailer.AllGoods().Range(func(key, value any) bool {
		goodsPs := (value.(*trailer.GoodsProcess))
		Schemas = append(Schemas, SchemaDetail{
			Name:        goodsPs.Name,
			LocalPath:   goodsPs.LocalPath,
			NetAddr:     goodsPs.NetAddr,
			CreateTs:    0,
			Size:        0,
			StorePath:   "",
			Description: goodsPs.Description,
		})
		return true
	})
	return Schemas
}

/*
*
* 表结构
*
 */

func SchemaDefineList() ([]SchemaDefine, error) {
	var err error
	ColumnsMap := []SchemaDefine{}
	Columns := []Column{}
	trailer.AllGoods().Range(func(key, value any) bool {
		goodsPs := (value.(*trailer.GoodsProcess))
		grpcConnection, err1 := grpc.Dial(goodsPs.NetAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err1 != nil {
			glogger.GLogger.Error(err1)
			err = err1
			return false
		}
		defer grpcConnection.Close()
		client := trailer.NewTrailerClient(grpcConnection)
		columns, err2 := client.Schema(context.Background(), &trailer.SchemaRequest{})
		if err2 != nil {
			glogger.GLogger.Error(err2)
			err = err2
			return false
		}
		for _, column := range columns.Columns {
			Columns = append(Columns, Column{
				Name:        string(column.Name),
				Type:        string(column.Type),
				Description: string(column.Description),
			})
		}
		Define := SchemaDefine{
			UUID:    goodsPs.Uuid,
			Columns: Columns,
		}
		ColumnsMap = append(ColumnsMap, Define)
		return true
	})
	return ColumnsMap, err
}

/*
*
* 获取仓库详情, 现阶段写死的, 后期会在proto中实现
*
 */
func GetSchemaDetail(goodsId string) SchemaDetail {
	return SchemaDetail{
		Name:        "Test RPC",
		LocalPath:   "/root/app1",
		NetAddr:     "127.0.0.1:4567",
		CreateTs:    uint64(time.Now().Unix()),
		Size:        12.34,
		StorePath:   "/root/data/test.db",
		Description: "An simply demo",
	}
}

/*
*
* 查询，第一个参数是查询请求，针对Sqlite就是SQL，针对mongodb就是JS，根据具体情况而定
*
 */
func Query(goodsId, query string) ([]Column, error) {
	var err error
	Columns := []Column{}
	if goodsPs := trailer.Get(goodsId); goodsPs != nil {
		grpcConnection, err1 := grpc.Dial(goodsPs.NetAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err1 != nil {
			glogger.GLogger.Error(err1)
			err = err1
			return Columns, err
		}
		defer grpcConnection.Close()
		client := trailer.NewTrailerClient(grpcConnection)
		columns, err2 := client.Query(context.Background(), &trailer.DataRowsRequest{
			Query: []byte(query),
		})
		if err2 != nil {
			glogger.GLogger.Error(err2)
			err = err2
			return Columns, err
		}
		for _, column := range columns.Column {
			Columns = append(Columns, Column{
				Name:  string(column.Name),
				Value: string(column.Value),
			})
		}
	}
	return Columns, err
}
