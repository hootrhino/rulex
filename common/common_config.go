package common

/*
*
* 通用的含有主机:端口的这类配置
*
 */
type HostConfig struct {
	Host string `json:"host" validate:"required" title:"服务地址" info:""`
	Port int    `json:"port" validate:"required" title:"服务端口" info:""`
}

/*
*
* IP 地址
*
 */
type IpConfig struct {
	Timeout int `json:"timeout" validate:"required" title:"连接超时" info:""`
	// Weather allow AutoRequest?
	AutoRequest bool `json:"autoRequest" title:"启动轮询" info:""`
	// Request Frequency, default 5 second
	Frequency int64 `json:"frequency" validate:"required" title:"采集频率" info:""`
	// ["127.0.0.1", "127.0.0.2", "127.0.0.3"]
	Hosts []string `json:"hosts" validate:"required" title:"服务地址" info:""`
}

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
* 通用串口配置
*
 */
type GenericUartConfig struct {
	Tag      string `json:"tag" validate:"required" title:"数据Tag" info:"给数据打标签"`
	Uart     string `json:"uart" validate:"required" title:"串口路径" info:"本地系统的串口路径"`
	BaudRate int    `json:"baudRate" validate:"required" title:"波特率" info:"串口通信波特率"`
	DataBits int    `json:"dataBits" validate:"required" title:"数据位" info:"串口通信数据位"`
	// 结束符, 默认是 '\n'；但是可以自己定义
	Decollator string `json:"decollator" title:"协议分隔符" info:""`
	// Weather allow AutoRequest?
	AutoRequest bool `json:"autoRequest" title:"启动轮询" info:""`
	// Request Frequency, default 5 second
	Frequency int64  `json:"frequency" validate:"required" title:"采集频率" info:""`
	Timeout   int    `json:"timeout" validate:"required" title:"连接超时" info:""`
	Parity    string `json:"parity" validate:"required" title:"奇偶校验" info:"奇偶校验"`
	StopBits  int    `json:"stopBits" validate:"required" title:"停止位" info:"串口通信停止位"`
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
	// Weather allow AutoRequest?
	AutoRequest bool `json:"autoRequest" title:"启动轮询" info:""`
	// Request Frequency, default 5 second
	Frequency int64 `json:"frequency" validate:"required" title:"采集频率" info:""`
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
