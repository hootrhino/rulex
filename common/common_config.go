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
