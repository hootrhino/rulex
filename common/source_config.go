package common

type GenericIoTHUBMqttConfig struct {
	Mode       string `json:"mode" validate:"required" title:"模式:GW(网关)|DC(直连)"`
	ProductId  string `json:"productId" validate:"required" title:"产品名"`
	DeviceName string `json:"deviceName" validate:"required" title:"设备名"`
	//
	Host string `json:"host" validate:"required" title:"服务地址"`
	Port int    `json:"port" validate:"required" title:"服务端口"`
	//
	ClientId string `json:"clientId" validate:"required" title:"客户端ID"`
	Username string `json:"username" validate:"required" title:"连接账户"`
	Password string `json:"password" validate:"required" title:"连接密码"`
}

/*
*
* 自定义UDP协议
*
 */

type RULEXUdpConfig struct {
	Host          string `json:"host" validate:"required" title:"服务地址"`
	Port          int    `json:"port" validate:"required" title:"服务端口"`
	MaxDataLength int    `json:"maxDataLength" title:"最大数据包"`
}
