package common

/*
*
* 西门子PLC的一些配置
*
 */
type S1200Config struct {
	Host        string `json:"host" validate:"required" title:"IP地址"`          // 127.0.0.1
	Port        *int   `json:"port" validate:"required" title:"端口号"`           // 0
	Rack        *int   `json:"rack" validate:"required" title:"架号"`            // 0
	Slot        *int   `json:"slot" validate:"required" title:"槽号"`            // 1
	Model       string `json:"model" validate:"required" title:"型号"`           // s7-200 s7 1500
	Timeout     *int   `json:"timeout" validate:"required" title:"连接超时时间"`     // 5s
	IdleTimeout *int   `json:"idleTimeout" validate:"required" title:"心跳超时时间"` // 5s
	//
	// Weather allow AutoRequest?
	AutoRequest bool `json:"autoRequest" title:"启动轮询"`
	// Request Frequency, default 5 second
	Frequency int64        `json:"frequency" validate:"required" title:"采集频率"`
	Blocks    []S1200Block `json:"blocks" validate:"required" title:"采集配置"` // Db
}
type S1200Block struct {
	Tag     string `json:"tag" title:"数据tag"`  // 数据tag
	Address int    `json:"address" title:"地址"` // 地址
	Start   int    `json:"start" title:"起始地址"` // 起始地址
	Size    int    `json:"size" title:"数据长度"`  // 数据长度
}
type S1200BlockValue struct {
	Tag     string `json:"tag" title:"数据tag"`  // 数据tag
	Address int    `json:"address" title:"地址"` // 地址
	Start   int    `json:"start" title:"起始地址"` // 起始地址
	Size    int    `json:"size" title:"服务地址"`  // 数据长度
	Value   []byte `json:"value" title:"数据长度"` // 值
}
