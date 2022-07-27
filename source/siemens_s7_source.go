package source

import (
	"context"
	"encoding/json"
	"time"

	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"

	"github.com/robinson/gos7"
)

var _status typex.SourceState

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
	Host        string `json:"host" validate:"required" title:"IP地址" info:""`          // 127.0.0.1
	Rack        *int   `json:"rack" validate:"required" title:"架号" info:""`            // 0
	Slot        *int   `json:"slot" validate:"required" title:"槽号" info:""`            // 1
	Model       string `json:"model" validate:"required" title:"型号" info:""`           // s7-200 s7 1500
	Timeout     *int   `json:"timeout" validate:"required" title:"连接超时时间" info:""`     // 5s
	IdleTimeout *int   `json:"idleTimeout" validate:"required" title:"心跳超时时间" info:""` // 5s
	Frequency   *int   `json:"frequency" validate:"required" title:"采集频率" info:""`     // 5s
	Dbs         []db   `json:"dbs" validate:"required" title:"采集配置" info:""`           // Db
}
type siemensS7Source struct {
	typex.XStatus
	client      gos7.Client
	Host        string
	Rack        *int
	Slot        *int
	Model       string
	Timeout     *int
	IdleTimeout *int
	Frequency   *int
	Dbs         []db
}

func NewSiemensS7Source(e typex.RuleX) typex.XSource {
	s7 := siemensS7Source{}
	s7.RuleEngine = e
	return &s7
}

//
// 测试资源是否可用
//
func (s7 *siemensS7Source) Test(inEndId string) bool {
	return true
}

//
// 注册InEndID到资源
//

func (s7 *siemensS7Source) Init(inEndId string, cfg map[string]interface{}) error {
	s7.PointId = inEndId
	return nil
}

//
// 启动资源
//
func (s7 *siemensS7Source) Start(cctx typex.CCTX) error {
	s7.Ctx = cctx.Ctx
	s7.CancelCTX = cctx.CancelCTX

	config := s7.RuleEngine.GetInEnd(s7.PointId).Config
	var mainConfig siemensS7config
	if err := utils.BindSourceConfig(config, &mainConfig); err != nil {
		return err
	}
	handler := gos7.NewTCPClientHandler(mainConfig.Host, *mainConfig.Rack, *mainConfig.Slot)
	handler.Timeout = 5 * time.Second
	if err := handler.Connect(); err != nil {
		return err
	}
	handler.Timeout = time.Duration(*mainConfig.Timeout) * time.Second
	handler.IdleTimeout = time.Duration(*mainConfig.IdleTimeout) * time.Second
	s7.client = gos7.NewClient(handler)
	_status = typex.SOURCE_UP
	ticker := time.NewTicker(time.Duration(*mainConfig.Frequency) * time.Second)
	for _, d := range mainConfig.Dbs {
		glogger.GLogger.Infof("Start read: Tag:%v Address:%v Start:%v Size:%v", d.Tag, d.Address, d.Start, d.Size)
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
					_status = typex.SOURCE_DOWN
					glogger.GLogger.Error(err)
				} else {
					// glogger.GLogger.Info("client.AGReadDB dataBuffer:", dataBuffer)
					dbv := dbValue{Value: string(dataBuffer[:d.Size])}
					dbv.Tag = d.Tag
					dbv.Address = d.Address
					dbv.Start = d.Start
					dbv.Size = d.Size
					bytes, _ := json.Marshal(dbv)
					work, err := s7.RuleEngine.WorkInEnd(s7.RuleEngine.GetInEnd(s7.PointId), string(bytes))
					if !work {
						glogger.GLogger.Error(err)
					}
				}
			}

		}(s7.Ctx, d)
	}

	return nil
}

//
// 资源是否被启用
//
func (s7 *siemensS7Source) Enabled() bool {
	return true
}

//
// 数据模型, 用来描述该资源支持的数据, 对应的是云平台的物模型
//
func (s7 *siemensS7Source) DataModels() []typex.XDataModel {
	return s7.XDataModels
}

//
// 获取前端表单定义
//
func (s7 *siemensS7Source) Configs() *typex.XConfig {
	return core.GenInConfig(typex.SIEMENS_S7, "SIEMENS_S7", siemensS7config{})

}

//
// 重载: 比如可以在重启的时候把某些数据保存起来
//
func (s7 *siemensS7Source) Reload() {

}

//
// 挂起资源, 用来做暂停资源使用
//
func (s7 *siemensS7Source) Pause() {

}

//
// 获取资源状态
//
func (s7 *siemensS7Source) Status() typex.SourceState {
	return _status

}

//
// 获取资源绑定的的详情
//
func (s7 *siemensS7Source) Details() *typex.InEnd {
	return s7.RuleEngine.GetInEnd(s7.PointId)
}

//
// 驱动接口, 通常用来和硬件交互
//
func (s7 *siemensS7Source) Driver() typex.XExternalDriver {
	return nil
}

//
//
//
func (s7 *siemensS7Source) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

//
// 停止资源, 用来释放资源
//
func (s7 *siemensS7Source) Stop() {
	if s7.client != nil {
		s7.client = nil
	}
	s7.CancelCTX()
}

//
// 来自外面的数据
//
func (*siemensS7Source) DownStream([]byte) (int, error) {
	return 0, nil
}

//
// 上行数据
//
func (*siemensS7Source) UpStream([]byte) (int, error) {
	return 0, nil
}
