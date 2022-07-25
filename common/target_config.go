package common

/*
*
* MQTT 连接配置
*
 */
type MqttConfig struct {
	Host     string `json:"host" validate:"required" title:"服务地址" info:""`
	Port     int    `json:"port" validate:"required" title:"服务端口" info:""`
	ClientId string `json:"clientId" validate:"required" title:"客户端ID" info:""`
	Username string `json:"username" validate:"required" title:"连接账户" info:""`
	Password string `json:"password" validate:"required" title:"连接密码" info:""`
	PubTopic string `json:"pubTopic" title:"上报TOPIC" info:"上报TOPIC"` // 上报数据的 Topic
	SubTopic string `json:"subTopic" title:"订阅TOPIC" info:"订阅TOPIC"` // 上报数据的 Topic
}

/*
*
* Mongodb 配置
*
 */
type MongoConfig struct {
	MongoUrl   string `json:"mongoUrl" validate:"required"`
	Database   string `json:"database" validate:"required"`
	Collection string `json:"collection" validate:"required"`
}

// http://<fqdn>:<port>/rest/sql/[db_name]
// fqnd: 集群中的任一台主机 FQDN 或 IP 地址
// port: 配置文件中 httpPort 配置项，缺省为 6041
// db_name: 可选参数，指定本次所执行的 SQL 语句的默认数据库库名
// curl -u root:taosdata -d 'show databases;' 106.15.225.172:6041/rest/sql
type TDEngineConfig struct {
	Fqdn           string `json:"fqdn" validate:"required"`           // 服务地址
	Port           int    `json:"port" validate:"required"`           // 服务端口
	Username       string `json:"username" validate:"required"`       // 用户
	Password       string `json:"password" validate:"required"`       // 密码
	DbName         string `json:"dbName" validate:"required"`         // 数据库名
	CreateDbSql    string `json:"createDbSql" validate:"required"`    // 建库SQL
	CreateTableSql string `json:"createTableSql" validate:"required"` // 建表SQL
	InsertSql      string `json:"insertSql" validate:"required"`      // 插入SQL
	Url            string `json:"url"`
}

/*
*
* HTTP
*
 */
type HTTPConfig struct {
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

/*
*
*
*
 */
type GrpcConfig struct {
	Host string `json:"host" validate:"required"`
	Port int    `json:"port" validate:"required"`
	Type string `json:"type" validate:"required"`
}

/*
*
*
*
 */

type NatsConfig struct {
	User     string `json:"user" validate:"required"`
	Password string `json:"password" validate:"required"`
	Host     string `json:"host" validate:"required"`
	Port     string `json:"port" validate:"required"`
	Topic    string `json:"topic" validate:"required"`
}
