package source

import (
	"fmt"
	"github.com/i4de/rulex/common"
	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
)

//	从txt文件读取数据，作为数据源
type textInEndSource struct {
	typex.XStatus
	engine     *gin.Engine
	mainConfig common.TextConfig
	status     typex.SourceState
}

func NewTextInEndSource(e typex.RuleX) typex.XSource {
	h := textInEndSource{}
	gin.SetMode(gin.ReleaseMode)
	h.engine = gin.New()
	h.RuleEngine = e
	return &h
}
func (*textInEndSource) Configs() *typex.XConfig {
	return core.GenInConfig(typex.TEXT, "TEXT", common.TextConfig{})
}
func (hh *textInEndSource) Init(inEndId string, configMap map[string]interface{}) error {
	hh.PointId = inEndId
	if err := utils.BindSourceConfig(configMap, &hh.mainConfig); err != nil {
		return err
	}
	return nil
}

//
func (hh *textInEndSource) Start(cctx typex.CCTX) error {
	hh.Ctx = cctx.Ctx
	hh.CancelCTX = cctx.CancelCTX

	path := hh.mainConfig.Path
	file, err := os.Open(path)
	if err != nil {
		glogger.GLogger.Error(err)
	} else {
		glogger.GLogger.Info("TEXT source open on" + fmt.Sprintf("%v/%v", hh.mainConfig.Path, hh.mainConfig.Name))
	}
	//4、读取到相关的文件，在处理完之后，把IO它关闭,释放资源
	defer file.Close()
	//5、开始读取文件的内容
	all, err := ioutil.ReadAll(file)
	if err != nil {
		glogger.GLogger.Error(err)
	} else {
		glogger.GLogger.Info("TEXT source read on" + fmt.Sprintf("%v/%v", hh.mainConfig.Path, hh.mainConfig.Name))
	}
	//6、读取文件的内容后，输出到控制台
	glogger.GLogger.Info(all)

	hh.status = typex.SOURCE_UP

	return nil
}

//
func (mm *textInEndSource) DataModels() []typex.XDataModel {
	return mm.XDataModels
}

//
func (hh *textInEndSource) Stop() {
	hh.CancelCTX()
	hh.status = typex.SOURCE_STOP
}
func (hh *textInEndSource) Reload() {

}
func (hh *textInEndSource) Pause() {

}
func (hh *textInEndSource) Status() typex.SourceState {
	return hh.status
}

func (hh *textInEndSource) Test(inEndId string) bool {
	return true
}

func (hh *textInEndSource) Enabled() bool {
	return hh.Enable
}
func (hh *textInEndSource) Details() *typex.InEnd {
	return hh.RuleEngine.GetInEnd(hh.PointId)
}

func (*textInEndSource) Driver() typex.XExternalDriver {
	return nil
}

//
// 拓扑
//
func (*textInEndSource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

//
// 来自外面的数据
//
func (*textInEndSource) DownStream([]byte) (int, error) {
	return 0, nil
}

//
// 上行数据
//
func (*textInEndSource) UpStream([]byte) (int, error) {
	return 0, nil
}
