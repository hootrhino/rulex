package target

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/i4de/rulex/common"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"
)

/*
*
* TDengine 的资源输出支持,当前暂时支持HTTP接口的形式，后续逐步会增加UDP、TCP模式
*
 */

type tdEngineTarget struct {
	typex.XStatus
	client     http.Client
	mainConfig common.TDEngineConfig
	status     typex.SourceState
}
type tdrs struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Desc   string `json:"desc"`
}

func NewTdEngineTarget(e typex.RuleX) typex.XTarget {
	td := tdEngineTarget{
		client:     http.Client{},
		mainConfig: common.TDEngineConfig{},
	}
	td.RuleEngine = e
	td.status = typex.SOURCE_DOWN
	return &td

}
func (td *tdEngineTarget) test() bool {
	if err := execQuery(td.client,
		td.mainConfig.Username,
		td.mainConfig.Password,
		"SELECT CLIENT_VERSION();",
		td.url()); err != nil {
		glogger.GLogger.Error(err)
		return false
	}
	return true
}
func (td *tdEngineTarget) url() string {
	return fmt.Sprintf("http://%s:%v/rest/sql/%s",
		td.mainConfig.Fqdn, td.mainConfig.Port, td.mainConfig.DbName)
}

// 测试资源是否可用
func (td *tdEngineTarget) Test(inEndId string) bool {
	return td.test()
}

//
// 注册InEndID到资源
//

func (td *tdEngineTarget) Init(outEndId string, configMap map[string]interface{}) error {
	td.PointId = outEndId

	if err := utils.BindSourceConfig(configMap, &td.mainConfig); err != nil {
		return err
	}
	if td.test() {
		return nil
	}
	return errors.New("tdengine connect error")
}

// 启动资源
func (td *tdEngineTarget) Start(cctx typex.CCTX) error {
	td.Ctx = cctx.Ctx
	td.CancelCTX = cctx.CancelCTX
	//

	if err := execQuery(td.client, td.mainConfig.Username,
		td.mainConfig.Password, td.mainConfig.CreateDbSql, td.url()); err != nil {
		return err
	}
	if err := execQuery(td.client, td.mainConfig.Username,
		td.mainConfig.Password, td.mainConfig.CreateTableSql, td.url()); err != nil {
		return err
	} else {
		td.status = typex.SOURCE_UP
		return nil
	}

}

// 资源是否被启用
func (td *tdEngineTarget) Enabled() bool {
	return true
}

// 数据模型, 用来描述该资源支持的数据, 对应的是云平台的物模型
func (td *tdEngineTarget) DataModels() []typex.XDataModel {
	return td.XDataModels
}

// 重载: 比如可以在重启的时候把某些数据保存起来
func (td *tdEngineTarget) Reload() {

}

// 挂起资源, 用来做暂停资源使用
func (td *tdEngineTarget) Pause() {
}

// 获取资源状态
func (td *tdEngineTarget) Status() typex.SourceState {
	return td.status
}

// 获取资源绑定的的详情
func (td *tdEngineTarget) Details() *typex.OutEnd {
	return td.RuleEngine.GetOutEnd(td.PointId)

}

// 驱动接口, 通常用来和硬件交互
func (td *tdEngineTarget) Driver() typex.XExternalDriver {
	return nil
}

func (td *tdEngineTarget) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

// 停止资源, 用来释放资源
func (td *tdEngineTarget) Stop() {
	td.status = typex.SOURCE_STOP
	td.CancelCTX()
}

func post(client http.Client,
	username string,
	password string,
	sql string,
	url string,
	headers map[string]string) (string, error) {
	body := strings.NewReader(sql)
	request, _ := http.NewRequest("POST", url, body)
	request.Header.Add("Content-Type", "text/plain")
	request.SetBasicAuth(username, password)
	response, err2 := client.Do(request)
	if err2 != nil {
		return "", err2
	}
	if response.StatusCode != 200 {
		bytes0, err3 := io.ReadAll(response.Body)
		if err3 != nil {
			return "", err3
		}
		return "", fmt.Errorf("Error:%v", string(bytes0))
	}
	bytes1, err3 := io.ReadAll(response.Body)
	if err3 != nil {
		return "", err3
	}
	return string(bytes1), nil
}

/*
*
* 执行TdEngine的查询
*
 */
func execQuery(client http.Client, username string, password string, sql string, url string) error {
	var r tdrs
	// {"status":"error","code":534,"desc":"Syntax error in SQL"}
	body, err1 := post(client, username, password, sql, url, map[string]string{})
	if err1 != nil {
		return err1
	}
	err2 := utils.TransformConfig([]byte(body), &r)
	if err2 != nil {
		return err2
	}
	if r.Status == "error" {
		return fmt.Errorf("code;%v, error:%s", r.Code, r.Desc)
	}
	return nil
}

/*
* SQL: INSERT INTO meter VALUES (NOW, %v, %v);
* 数据到达后写入Tdengine, 这里对数据有严格约束，必须是以,分割的字符串
* 比如: 10.22,220.12,123,......
*
 */
func (td *tdEngineTarget) To(data interface{}) (interface{}, error) {
	switch s := data.(type) {
	case string:
		{
			ss := strings.Split(s, ",")
			insertSql := td.mainConfig.InsertSql
			for _, v := range ss {
				insertSql = strings.Replace(insertSql, "%v", strings.TrimSpace(v), 1)
			}
			return execQuery(td.client, td.mainConfig.Username,
				td.mainConfig.Password, insertSql, td.url()), nil
		}
	}
	return nil, nil
}

/*
*
* 配置
*
 */
func (*tdEngineTarget) Configs() *typex.XConfig {
	return &typex.XConfig{}
}
