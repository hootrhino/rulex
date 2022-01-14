package resource

import (
	"context"
	"encoding/json"
	"rulex/typex"
	"rulex/utils"
	"time"

	"github.com/ngaut/log"
	"github.com/robinson/gos7"
)

/*
*
* stateAddress: 配置一个地址给RULEX采集状态用，一般来说，只要有返回值就表示成功
*
 */
type stateAddress struct {
	Address int `json:"address"` // 地址
	Start   int `json:"start"`   // 起始地址
	Size    int `json:"size"`    // 数据长度
}
type db struct {
	Tag     string `json:"tag"`     // 数据tag
	Address int    `json:"address"` // 地址
	Start   int    `json:"start"`   // 起始地址
	Size    int    `json:"size"`    // 数据长度
}
type dbValue struct {
	db
	Value string `json:"value"`
}
type siemensS7config struct {
	Host         string       `json:"host" validate:"required" title:"IP地址" info:""`          // 127.0.0.1
	Rack         *int         `json:"rack" validate:"required" title:"架号" info:""`            // 0
	Slot         *int         `json:"slot" validate:"required" title:"槽号" info:""`            // 1
	Timeout      *int         `json:"timeout" validate:"required" title:"连接超时时间" info:""`     // 5s
	IdleTimeout  *int         `json:"idleTimeout" validate:"required" title:"心跳超时时间" info:""` // 5s
	Frequency    *int         `json:"frequency" validate:"required" title:"采集频率" info:""`     // 5s
	StateAddress stateAddress `json:"stateAddress" validate:"required" title:"状态地址" info:""`  // stateAddress
	Dbs          []db         `json:"dbs" validate:"required" title:"采集配置" info:""`           // Db
}
type siemensS7Resource struct {
	typex.XStatus
	client       gos7.Client
	stateAddress stateAddress
}

func NewSiemensS7Resource(e typex.RuleX) typex.XResource {
	s7 := siemensS7Resource{}
	s7.RuleEngine = e
	return &s7
}

//
// 测试资源是否可用
//
func (s7 *siemensS7Resource) Test(inEndId string) bool {
	return true
}

//
// 注册InEndID到资源
//
func (s7 *siemensS7Resource) Register(inEndId string) error {
	s7.PointId = inEndId
	return nil
}

//
// 启动资源
//
func (s7 *siemensS7Resource) Start() error {
	config := s7.RuleEngine.GetInEnd(s7.PointId).Config
	var mainConfig siemensS7config
	if err := utils.BindResourceConfig(config, &mainConfig); err != nil {
		return err
	}
	handler := gos7.NewTCPClientHandler(mainConfig.Host, *mainConfig.Rack, *mainConfig.Slot)
	if err := handler.Connect(); err != nil {
		return err
	}
	handler.Timeout = time.Duration(*mainConfig.Timeout) * time.Second
	handler.IdleTimeout = time.Duration(*mainConfig.IdleTimeout) * time.Second
	s7.stateAddress = mainConfig.StateAddress
	s7.client = gos7.NewClient(handler)
	ticker := time.NewTicker(time.Duration(*mainConfig.Frequency) * time.Second)
	for _, d := range mainConfig.Dbs {
		log.Infof("Start read: Tag:%v Address:%v Start:%v Size:%v", d.Tag, d.Address, d.Start, d.Size)
		go func(ctx context.Context, d db) {

			dataBuffer := make([]byte, 512)
			for {

				<-ticker.C
				select {
				case <-ctx.Done():
					{
						return
					}
				default:
					{

					}
				}
				if s7.client == nil {
					return
				}
				err := s7.client.AGReadDB(d.Address, d.Start, d.Size, dataBuffer)
				if err != nil {
					log.Error(err)
				} else {
					// log.Info("client.AGReadDB dataBuffer:", dataBuffer)
					dbv := dbValue{Value: string(dataBuffer[:d.Size])}
					dbv.Tag = d.Tag
					dbv.Address = d.Address
					dbv.Start = d.Start
					dbv.Size = d.Size
					bytes, _ := json.Marshal(dbv)
					s7.RuleEngine.Work(s7.RuleEngine.GetInEnd(s7.PointId), string(bytes))
				}
			}

		}(context.Background(), d)
	}

	return nil
}

//
// 资源是否被启用
//
func (s7 *siemensS7Resource) Enabled() bool {
	return true
}

//
// 数据模型, 用来描述该资源支持的数据, 对应的是云平台的物模型
//
func (s7 *siemensS7Resource) DataModels() []typex.XDataModel {
	return []typex.XDataModel{}
}

//
// 获取前端表单定义
//
func (s7 *siemensS7Resource) Configs() typex.XConfig {
	return typex.XConfig{}
}

//
// 重载: 比如可以在重启的时候把某些数据保存起来
//
func (s7 *siemensS7Resource) Reload() {

}

//
// 挂起资源, 用来做暂停资源使用
//
func (s7 *siemensS7Resource) Pause() {

}

//
// 获取资源状态
//
func (s7 *siemensS7Resource) Status() typex.ResourceState {
	if s7.client != nil {
		b := make([]byte, 1)
		err := s7.client.AGReadDB(s7.stateAddress.Address, s7.stateAddress.Start, s7.stateAddress.Size, b)
		if err != nil {
			log.Error(err)
			return typex.DOWN
		} else {
			return typex.UP
		}
	} else {
		return typex.DOWN
	}

}

//
// 获取资源绑定的的详情
//
func (s7 *siemensS7Resource) Details() *typex.InEnd {
	return s7.RuleEngine.GetInEnd(s7.PointId)
}

//
// 不经过规则引擎处理的直达数据接口
//
func (s7 *siemensS7Resource) OnStreamApproached(data string) error {
	return nil
}

//
// 驱动接口, 通常用来和硬件交互
//
func (s7 *siemensS7Resource) Driver() typex.XExternalDriver {
	return nil
}

//
//
//
func (s7 *siemensS7Resource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

//
// 停止资源, 用来释放资源
//
func (s7 *siemensS7Resource) Stop() {
	if s7.client != nil {
		s7.client = nil
	}
	context.Background().Done()
}
