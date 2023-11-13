package source

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	"github.com/gin-gonic/gin"
)

type httpInEndSource struct {
	typex.XStatus
	engine     *gin.Engine
	mainConfig common.HostConfig
	status     typex.SourceState
}

func NewHttpInEndSource(e typex.RuleX) typex.XSource {
	h := httpInEndSource{}
	gin.SetMode(gin.ReleaseMode)
	h.engine = gin.New()
	h.RuleEngine = e
	return &h
}

func (hh *httpInEndSource) Init(inEndId string, configMap map[string]interface{}) error {
	hh.PointId = inEndId
	if err := utils.BindSourceConfig(configMap, &hh.mainConfig); err != nil {
		return err
	}
	return nil
}

func (hh *httpInEndSource) Start(cctx typex.CCTX) error {
	hh.Ctx = cctx.Ctx
	hh.CancelCTX = cctx.CancelCTX

	hh.engine.POST("/in", func(c *gin.Context) {
		type Form struct {
			Data string `json:"data"`
		}
		var inForm Form
		err := c.BindJSON(&inForm)
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
		} else {
			hh.RuleEngine.WorkInEnd(hh.RuleEngine.GetInEnd(hh.PointId), inForm.Data)
			c.JSON(200, gin.H{
				"message": "success",
				"code":    0,
			})
		}
	})

	go func(ctx context.Context) {
		err := http.ListenAndServe(fmt.Sprintf(":%v", hh.mainConfig.Port), hh.engine)
		if err != nil {
			glogger.GLogger.Error(err)
			return
		}
	}(hh.Ctx)
	hh.status = typex.SOURCE_UP
	glogger.GLogger.Info("HTTP source started on" + " [0.0.0.0]:" + fmt.Sprintf("%v", hh.mainConfig.Port))

	return nil
}

func (mm *httpInEndSource) DataModels() []typex.XDataModel {
	return mm.XDataModels
}

func (hh *httpInEndSource) Stop() {
	hh.status = typex.SOURCE_STOP
	if hh.CancelCTX != nil {
		hh.CancelCTX()
	}
}

func (hh *httpInEndSource) Status() typex.SourceState {
	return hh.status
}

func (hh *httpInEndSource) Test(inEndId string) bool {
	return true
}


func (hh *httpInEndSource) Details() *typex.InEnd {
	return hh.RuleEngine.GetInEnd(hh.PointId)
}

func (*httpInEndSource) Driver() typex.XExternalDriver {
	return nil
}

// 拓扑
func (*httpInEndSource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

// 来自外面的数据
func (*httpInEndSource) DownStream([]byte) (int, error) {
	return 0, nil
}

// 上行数据
func (*httpInEndSource) UpStream([]byte) (int, error) {
	return 0, nil
}
