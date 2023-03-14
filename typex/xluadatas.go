package typex

// LUA 表增强
type LMap map[string]LObject
type LList []LObject

func (lMap LMap) ToString() string {
	return ""
}
func (lMap LList) ToString() string {
	return ""
}

type LObject struct {
	Type  int
	Value interface{}
}

func (obj LObject) ToI32() int32 {
	return 0
}
func (obj LObject) ToI64() int64 {
	return 0
}
func (obj LObject) ToF32() float32 {
	return 0
}
func (obj LObject) ToF64() float64 {
	return 0
}
func (obj LObject) ToString() string {
	return ""
}
