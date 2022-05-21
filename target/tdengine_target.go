package target

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"rulex/core"
	"rulex/typex"
	"rulex/utils"
	"strings"

	"github.com/ngaut/log"
)

/*
*
* TDengine 的资源输出支持,当前暂时支持HTTP接口的形式，后续逐步会增加UDP、TCP模式
*
 */

// http://<fqdn>:<port>/rest/sql/[db_name]
// fqnd: 集群中的任一台主机 FQDN 或 IP 地址
// port: 配置文件中 httpPort 配置项，缺省为 6041
// db_name: 可选参数，指定本次所执行的 SQL 语句的默认数据库库名
// curl -u root:taosdata -d 'show databases;' 106.15.225.172:6041/rest/sql
type tdEngineConfig struct {
	Fqdn           string `json:"fqdn" validate:"required"`           // 服务地址
	Port           int    `json:"port" validate:"required"`           // 服务端口
	Username       string `json:"username" validate:"required"`       // 用户
	Password       string `json:"password" validate:"required"`       // 密码
	DbName         string `json:"dbName" validate:"required"`         // 数据库名
	CreateDbSql    string `json:"createDbSql" validate:"required"`    // 建库SQL
	CreateTableSql string `json:"createTableSql" validate:"required"` // 建表SQL
	InsertSql      string `json:"insertSql" validate:"required"`      // 插入SQL
}
type tdEngineTarget struct {
	typex.XStatus
	client         http.Client
	Fqdn           string
	Port           int
	Username       string
	Password       string
	DbName         string
	CreateDbSql    string
	CreateTableSql string
	InsertSql      string
	Url            string
}
type tdrs struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Desc   string `json:"desc"`
}

func NewTdEngineTarget(e typex.RuleX) typex.XTarget {
	td := tdEngineTarget{
		client: http.Client{},
	}
	td.RuleEngine = e
	return &td

}

//
// 测试资源是否可用
//
func (td *tdEngineTarget) Test(inEndId string) bool {
	if err := execQuery(td.client, td.Username, td.Password, "SELECT CLIENT_VERSION();", td.Url); err != nil {
		log.Error(err)
		return false
	}
	return true
}

//
// 注册InEndID到资源
//

func (td *tdEngineTarget) Init(outEndId string, cfg map[string]interface{}) error {
	td.PointId = outEndId
	return nil
}

//
// 启动资源
//
func (td *tdEngineTarget) Start(cctx typex.CCTX) error {
	td.Ctx = cctx.Ctx
	td.CancelCTX = cctx.CancelCTX
	// http://<fqdn>:<port>/rest/sql/[db_name]
	// curl -u root:taosdata -d 'show databases;' 127.0.0.1:6041/rest/sql
	config := td.RuleEngine.GetOutEnd(td.PointId).Config
	var mainConfig tdEngineConfig
	if err := utils.BindSourceConfig(config, &mainConfig); err != nil {
		return err
	}
	td.Fqdn = mainConfig.Fqdn
	td.Port = mainConfig.Port
	td.Username = mainConfig.Username
	td.Password = mainConfig.Password
	td.DbName = mainConfig.DbName
	td.CreateDbSql = mainConfig.CreateDbSql
	td.CreateTableSql = mainConfig.CreateTableSql
	td.InsertSql = mainConfig.InsertSql
	td.Url = fmt.Sprintf("http://%s:%v/rest/sql/%s", td.Fqdn, td.Port, td.DbName)
	if err := execQuery(td.client, td.Username, td.Password, td.CreateDbSql, td.Url); err != nil {
		return err
	}
	return execQuery(td.client, td.Username, td.Password, td.CreateTableSql, td.Url)
}

//
// 资源是否被启用
//
func (td *tdEngineTarget) Enabled() bool {
	return true
}

//
// 数据模型, 用来描述该资源支持的数据, 对应的是云平台的物模型
//
func (td *tdEngineTarget) DataModels() []typex.XDataModel {
	return td.XDataModels
}

//
// 重载: 比如可以在重启的时候把某些数据保存起来
//
func (td *tdEngineTarget) Reload() {

}

//
// 挂起资源, 用来做暂停资源使用
//
func (td *tdEngineTarget) Pause() {
}

//
// 获取资源状态
//
func (td *tdEngineTarget) Status() typex.SourceState {
	if err := execQuery(td.client, td.Username, td.Password, "SELECT CLIENT_VERSION();", td.Url); err != nil {
		log.Error(err)
		return typex.DOWN
	}
	return typex.UP
}

//
// 获取资源绑定的的详情
//
func (td *tdEngineTarget) Details() *typex.OutEnd {
	return td.RuleEngine.GetOutEnd(td.PointId)

}



//
// 驱动接口, 通常用来和硬件交互
//
func (td *tdEngineTarget) Driver() typex.XExternalDriver {
	return nil
}

//
//
//
func (td *tdEngineTarget) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

//
// 停止资源, 用来释放资源
//
func (td *tdEngineTarget) Stop() {
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
		bytes0, err3 := ioutil.ReadAll(response.Body)
		if err3 != nil {
			return "", err3
		}
		return "", fmt.Errorf("Error:%v", string(bytes0))
	}
	bytes1, err3 := ioutil.ReadAll(response.Body)
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
*
* 数据到达后写入Tdengine, 这里对数据有严格约束，必须是以,分割的字符串
* 比如: 10.22,220.12,123,......
*
 */
func (td *tdEngineTarget) To(data interface{}) (interface{}, error) {
	switch s := data.(type) {
	case string:
		{
			ss := strings.Split(s, ",")
			insertSql := td.InsertSql
			for _, v := range ss {
				insertSql = strings.Replace(insertSql, "%v", strings.TrimSpace(v), 1)
			}
			return execQuery(td.client, td.Username, td.Password, insertSql, td.Url), nil
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
	return core.GenOutConfig(typex.TDENGINE_TARGET, "TDENGINE_TARGET", httpConfig{})
}
