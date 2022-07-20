package driver

/*
*
* 西门子PLC的一些配置
*
 */

type S1200Block struct {
	Tag     string `json:"tag"`     // 数据tag
	Address int    `json:"address"` // 地址
	Start   int    `json:"start"`   // 起始地址
	Size    int    `json:"size"`    // 数据长度
}
type S1200BlockValue struct {
	Tag     string `json:"tag"`     // 数据tag
	Address int    `json:"address"` // 地址
	Start   int    `json:"start"`   // 起始地址
	Size    int    `json:"size"`    // 数据长度
	Value   []byte `json:"value"`
}
