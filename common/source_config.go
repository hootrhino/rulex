package common

//
//
//
type TencentMqttConfig struct {
	ProductId  string `json:"productId" validate:"required" title:"产品名" info:""`
	DeviceName string `json:"deviceName" validate:"required" title:"设备名" info:""`
	//
	Host string `json:"host" validate:"required" title:"服务地址" info:""`
	Port int    `json:"port" validate:"required" title:"服务端口" info:""`
	//
	ClientId string `json:"clientId" validate:"required" title:"客户端ID" info:""`
	Username string `json:"username" validate:"required" title:"连接账户" info:""`
	Password string `json:"password" validate:"required" title:"连接密码" info:""`
}

/*
*
* 自定义UDP协议
*
 */

type RULEXUdpConfig struct {
	Host          string `json:"host" validate:"required" title:"服务地址" info:""`
	Port          int    `json:"port" validate:"required" title:"服务端口" info:""`
	MaxDataLength int    `json:"maxDataLength" validate:"required" title:"最大数据包" info:""`
}

// TextConfig 文本文件相关参数配置
type TextConfig struct {
	Path string `json:"path" validate:"required" title:"文本文件路径" info:""`
	Name string `json:"name" validate:"required" title:"文本文件名称" info:""`
}
