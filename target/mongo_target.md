<!--
 Copyright (C) 2023 wwhai

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Affero General Public License as
 published by the Free Software Foundation, either version 3 of the
 License, or (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Affero General Public License for more details.

 You should have received a copy of the GNU Affero General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
-->

# Mongodb 客户端
## 简介
本组件实现将数据写入Mongodb。
## 配置
```go
type MongoConfig struct {
	MongoUrl   string `json:"mongoUrl" validate:"required" title:"URL"`
	Database   string `json:"database" validate:"required" title:"数据库"`
	Collection string `json:"collection" validate:"required" title:"集合"`
}
```
字段解释
- mongoUrl: Mongodb的URL，例如：mongodb://localhost:27017
- database: 数据库
- collection: 数据集合

## 示例

快速启动一个Docker Mongo容器：

```sh
docker run -d -v "./data/db" -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=123456 -p 27017:27017 --name mongodb-container mongo
```
编写脚本

```lua
function(args)
    local err = data:DToMongo('uuid', args)
	return true, args
end
```
然后查看Mongodb的集合看是否有数据写入。