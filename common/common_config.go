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
	Tag       string `json:"tag" validate:"required" title:"数据Tag" info:"给数据打标签"`
	Uart      string `json:"uart" validate:"required" title:"串口路径" info:"本地系统的串口路径"`
	BaudRate  int    `json:"baudRate" validate:"required" title:"波特率" info:"串口通信波特率"`
	DataBits  int    `json:"dataBits" validate:"required" title:"数据位" info:"串口通信数据位"`
	Frequency int64  `json:"frequency" validate:"required" title:"采集频率" info:""`
	Timeout   int    `json:"timeout" validate:"required" title:"连接超时" info:""`
	Parity    string `json:"parity" validate:"required" title:"奇偶校验" info:"奇偶校验"`
	StopBits  int    `json:"stopBits" validate:"required" title:"停止位" info:"串口通信停止位"`
}
