package common

/*
*
* 西门子PLC的一些配置
*
 */
type S1200Config struct {
	Host        string       `json:"host" validate:"required" title:"IP地址:端口号"`      // 127.0.0.1
	Model       string       `json:"model" validate:"required" title:"型号"`           // s7-200 s7 1500
	Rack        *int         `json:"rack" validate:"required" title:"架号"`            // 0
	Slot        *int         `json:"slot" validate:"required" title:"槽号"`            // 1
	Timeout     *int         `json:"timeout" validate:"required" title:"连接超时时间"`     // 5s
	IdleTimeout *int         `json:"idleTimeout" validate:"required" title:"心跳超时时间"` // 5s
	AutoRequest *bool        `json:"autoRequest" title:"启动轮询"`
	Blocks      []S1200Block `json:"blocks" validate:"required" title:"采集配置"` // Db
}
type S1200Block struct {
	Tag       string `json:"tag" validate:"required" title:"数据tag"` // 数据tag
	Type      string `json:"type" validate:"required" title:"地址"`   // MB | DB |FB
	Frequency int64  `json:"frequency" validate:"required" title:"采集频率"`
	Address   int    `json:"address" validate:"required" title:"地址"`            // 地址
	Start     int    `json:"start" validate:"required" title:"起始地址"`            // 起始地址
	Size      int    `json:"size" validate:"required" title:"服务地址"`             // 数据长度
	Value     string `json:"value,omitempty" validate:"required" title:"数据Hex"` // 值
}
