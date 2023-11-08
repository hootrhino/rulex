package common

/*
*
* 通用的含有主机:端口的这类配置
*
 */
type HostConfig struct {
	Host    string `json:"host" validate:"required" title:"服务地址"`
	Port    int    `json:"port" validate:"required" title:"服务端口"`
	Timeout int    `json:"timeout,omitempty" title:"连接超时"`
}

/*
*
* IP 地址
*
 */
type IcmpConfig struct {
	Timeout int `json:"timeout" validate:"required" title:"连接超时"`
	// ["127.0.0.1", "127.0.0.2", "127.0.0.3"]
	Hosts []string `json:"hosts" validate:"required" title:"服务地址"`
}

/*
*
* MQTT 连接配置
*
 */
type MqttConfig struct {
	Host     string `json:"host" validate:"required" title:"服务地址"`
	Port     int    `json:"port" validate:"required" title:"服务端口"`
	ClientId string `json:"clientId" validate:"required" title:"客户端ID"`
	Username string `json:"username" validate:"required" title:"连接账户"`
	Password string `json:"password" validate:"required" title:"连接密码"`
	PubTopic string `json:"pubTopic" title:"上报TOPIC" info:"上报TOPIC"` // 上报数据的 Topic
	SubTopic string `json:"subTopic" title:"订阅TOPIC" info:"订阅TOPIC"` // 上报数据的 Topic
}

/*
*
* 4.19重构
*
 */
type CommonUartConfig struct {
	Timeout  int    `json:"timeout" validate:"required"`
	Uart     string `json:"uart" validate:"required"`
	BaudRate int    `json:"baudRate" validate:"required"`
	DataBits int    `json:"dataBits" validate:"required"`
	Parity   string `json:"parity" validate:"required"`
	StopBits int    `json:"stopBits" validate:"required"`
}

/*
*
* SNMP 配置
*
 */
type GenericSnmpConfig struct {
	// Target is an ipv4 address.
	Target string `json:"target" validate:"required" title:"Target" info:"Target"`
	// Port is a port.
	Port uint16 `json:"port" validate:"required" title:"Port" info:"Port"`
	// Transport is the transport protocol to use ("udp" or "tcp"); if unset "udp" will be used.
	Transport string `json:"transport" validate:"required" title:"Transport" info:"Transport"`
	// Community is an SNMP Community string.
	Community string `json:"community" validate:"required" title:"Community" info:"Community"`
}

/*
*
* Sqlite 配置
*
 */
type SqliteConfig struct {
	// 本地数据库名称
	DbName string `json:"dbName" validate:"required"`
	// 数据表名
	CreateTbSql string `json:"createTbSql" validate:"required"`
	// 插入语句, 变量用 ？ 替代，会被替换成实际值
	// Eg: insert into db1.tb1 value(v1, v2, v3)
	InsertSql string `json:"insertSql" validate:"required"`
}
