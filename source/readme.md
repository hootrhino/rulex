##

### 开发模板
资源需要实现以下接口:
```go

type XSource interface {
	Test(inEndId string) bool
	Init(inEndId string, cfg map[string]interface{}) error
	Start(CCTX) error
	DataModels() []XDataModel
	Status() SourceState
	Details() *InEnd
	Driver() XExternalDriver
	Topology() []TopologyPoint
	Stop()
}

```
### 加载配置
我们以一个 `COAP Server` 为例来解释,首先定义一个配置结构体:
```go
type coAPConfig struct {
	Port       uint16             `json:"port" validate:"required" title:"端口"`
	DataModels []typex.XDataModel `json:"dataModels" title:"数据模型"`
}
```
然后在Init里面解析外部配置到资源的配置结构体，相当于是加载配置:
```go

func (cc *coAPInEndSource) Init(inEndId string, cfg map[string]interface{}) error {
	cc.PointId = inEndId
	var mainConfig coAPConfig
	if err := utils.BindSourceConfig(cfg, &mainConfig); err != nil {
		return err
	}
	cc.port = mainConfig.Port
	cc.dataModels = mainConfig.DataModels
	return nil
}
```
参数说明
- inEndId：资源元数据的ID
- cfg：配置数据的Map形式结构