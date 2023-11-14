package common

/*
*
* Mongodb 配置
*
 */
type MongoConfig struct {
	MongoUrl   string `json:"mongoUrl" validate:"required" title:"URL"`
	Database   string `json:"database" validate:"required" title:"数据库"`
	Collection string `json:"collection" validate:"required" title:"集合"`
}

// http://<fqdn>:<port>/rest/sql/[db_name]
// fqnd: 集群中的任一台主机 FQDN 或 IP 地址
// port: 配置文件中 httpPort 配置项，缺省为 6041
// db_name: 可选参数，指定本次所执行的 SQL 语句的默认数据库库名
// curl -u root:taosdata -d 'show databases;' 106.15.225.172:6041/rest/sql
type TDEngineConfig struct {
	Fqdn           string `json:"fqdn" validate:"required" title:"地址"`              // 服务地址
	Port           int    `json:"port" validate:"required" title:"端口"`              // 服务端口
	Username       string `json:"username" validate:"required" title:"用户"`          // 用户
	Password       string `json:"password" validate:"required" title:"密码"`          // 密码
	DbName         string `json:"dbName" validate:"required" title:"数据库名"`          // 数据库名
	CreateDbSql    string `json:"createDbSql" validate:"required" title:"建库SQL"`    // 建库SQL
	CreateTableSql string `json:"createTableSql" validate:"required" title:"建表SQL"` // 建表SQL
}

/*
*
* HTTP
*
 */
type HTTPConfig struct {
	Url     string            `json:"url" validate:"required" title:"URL"`
	Headers map[string]string `json:"headers" validate:"required" title:"HTTP Headers"`
}

/*
*
*
*
 */
type GrpcConfig struct {
	Host string `json:"host" validate:"required" title:"地址"`
	Port int    `json:"port" validate:"required" title:"端口"`
	Type string `json:"type" title:"类型"`
}

/*
*
*
*
 */

type NatsConfig struct {
	Username string `json:"username" validate:"required" title:"用户"`
	Password string `json:"password" validate:"required" title:"密码"`
	Host     string `json:"host" validate:"required" title:"地址"`
	Port     int    `json:"port" validate:"required" title:"端口"`
	Topic    string `json:"topic" validate:"required" title:"转发Topic"`
}
