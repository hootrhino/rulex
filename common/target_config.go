package common

/*
*
* Mongodb 配置
*
 */
type MongoConfig struct {
	MongoUrl   string `json:"mongoUrl" validate:"required" title:"URL" info:""`
	Database   string `json:"database" validate:"required" title:"数据库" info:""`
	Collection string `json:"collection" validate:"required" title:"集合" info:""`
}

// http://<fqdn>:<port>/rest/sql/[db_name]
// fqnd: 集群中的任一台主机 FQDN 或 IP 地址
// port: 配置文件中 httpPort 配置项，缺省为 6041
// db_name: 可选参数，指定本次所执行的 SQL 语句的默认数据库库名
// curl -u root:taosdata -d 'show databases;' 106.15.225.172:6041/rest/sql
type TDEngineConfig struct {
	Fqdn           string `json:"fqdn" validate:"required" title:"地址" info:""`              // 服务地址
	Port           int    `json:"port" validate:"required" title:"端口" info:""`              // 服务端口
	Username       string `json:"username" validate:"required" title:"用户" info:""`          // 用户
	Password       string `json:"password" validate:"required" title:"密码" info:""`          // 密码
	DbName         string `json:"dbName" validate:"required" title:"数据库名" info:""`          // 数据库名
	CreateDbSql    string `json:"createDbSql" validate:"required" title:"建库SQL" info:""`    // 建库SQL
	CreateTableSql string `json:"createTableSql" validate:"required" title:"建表SQL" info:""` // 建表SQL
	InsertSql      string `json:"insertSql" validate:"required" title:"写入SQL" info:""`      // 插入SQL
}

/*
*
* HTTP
*
 */
type HTTPConfig struct {
	Url     string            `json:"url" title:"URL" info:""`
	Headers map[string]string `json:"headers" title:"HTTP Headers" info:""`
}

/*
*
*
*
 */
type GrpcConfig struct {
	Host string `json:"host" validate:"required" title:"地址" info:""`
	Port int    `json:"port" validate:"required" title:"端口" info:""`
	Type string `json:"type" title:"类型" info:""`
}

/*
*
*
*
 */

type NatsConfig struct {
	User     string `json:"user" validate:"required" title:"用户" info:""`
	Password string `json:"password" validate:"required" title:"密码" info:""`
	Host     string `json:"host" validate:"required" title:"地址" info:""`
	Port     int    `json:"port" validate:"required" title:"端口" info:""`
	Topic    string `json:"topic" validate:"required" title:"转发Topic" info:""`
}
