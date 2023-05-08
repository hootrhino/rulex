package device

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"github.com/i4de/rulex/common"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"
	"strconv"
	"sync"
	"time"
)

/*
Endpoint：OPC UA服务端的地址，对应gopcua库中的ClientEndpoint
CertFile：PEM格式的证书文件路径，对应gopcua库中的ClientCertFile
KeyFile：PEM格式的私钥文件路径，对应gopcua库中的ClientKeyFile
GenCert：是否生成新的证书，对应gopcua库中的ClientGenCert
Policy：安全策略URL，可以是None、Basic128Rsa15、Basic256、Basic256Sha256中的任意一个，对应gopcua库中的ClientSecurityPolicyURI
Mode：安全模式，可以是None、Sign、SignAndEncrypt中的任意一个，对应gopcua库中的ClientSecurityMode
Auth：认证模式，可以是Anonymous、UserName、Certificate中的任意一个，对应gopcua库中的ClientAuthMode。
*/

type opcuaCommonconfig struct {
	Endpoint  string `json:"endpoint" title:"服务器URL" example:"opc.tcp://NOAH:53530/OPCUA/SimulationServer" info:""`
	Policy    string `json:"policy" title:"消息安全模式" flag:"None、Basic128Rsa15、Basic256、Basic256Sha256" info:""` //可选 四种模式：无、Basic128Rsa15、Basic256、Basic256Sha256
	Mode      string `json:"mode" title:"消息安全模式" flag:"None, Sign, SignAndEncrypt" info:""`                   //可选 三种模式：无、签名、签名加密
	Auth      string `json:"auth" title:"认证方式" title:"one of Anonymous, UserName"  info:""`                   //可选 二种模式：匿名、用户名
	Username  string `json:"username" title:"用户名" info:""`
	Password  string `json:"password" title:"密码" info:""`
	Timeout   int    `json:"timeout" title:"超时" info:""`
	Frequency int64  `json:"frequency" title:"采集频率" validate:"required" info:""`
	RetryTime int    `json:"retryTime" title:"错误次数" validate:"required"` // 几次以后重启,0 表示不重启
}
type OpcuaNode struct {
	Tag         string `json:"tag" validate:"required" title:"数据Tag" info:""`
	Description string `json:"description" validate:"required"`
	NodeID      string `json:"nodeId" validate:"required" title:"NodeID" example:"ns=1;s=Test"`
	DataType    string `json:"dataType" title:"数据类型" tag:"String" info:""`
	Value       string `json:"value" title:"值" info:"从OPCUA获取的值"` //不需要配置
}
type opcua_CustomProtocolConfig struct {
	OpcuaCommonconfig opcuaCommonconfig `json:"commonConfig" validate:"required"`
	Opcnodes          []OpcuaNode       `json:"opcuanodes" validate:"required" title:"采集节点" info:""`
}
type genericOpcuaDevice struct {
	typex.XStatus
	status       typex.DeviceState
	RuleEngine   typex.RuleX
	driver       typex.XExternalDriver
	client       *opcua.Client
	mainConfig   opcua_CustomProtocolConfig
	subscription *opcua.Subscription
	locker       sync.Locker
	errorCount   int // 记录最大容错数，默认5次，出错超过5此就重启
}
type PolicyFlag string

const (
	POLICY_NONE           PolicyFlag = "None"
	POLICY_BASIC128RSA15  PolicyFlag = "Basic128Rsa15"
	POLICY_BASIC256       PolicyFlag = "Basic256"
	POLICY_BASIC256SHA256 PolicyFlag = "Basic256Sha256"
)

// Auth 认证模式  枚举
type AuthType string

const (
	AUTH_ANONYMOUS AuthType = "Anonymous"
	AUTH_USERNAME  AuthType = "UserName"
)

type SecurityMode string

const (
	MODE_NONE             SecurityMode = "None"
	MODE_SIGN             SecurityMode = "Sign"
	MODE_SIGN_AND_ENCRYPT SecurityMode = "SignAndEncrypt"
)

func NewGenericOpcuaDevice(e typex.RuleX) typex.XDevice {
	opc := new(genericOpcuaDevice)
	opc.RuleEngine = e
	opc.locker = &sync.Mutex{}
	//opc.mainConfig = opcua_CustomProtocolConfig{}
	opc.mainConfig = opcua_CustomProtocolConfig{
		OpcuaCommonconfig: opcuaCommonconfig{},
		Opcnodes:          []OpcuaNode{},
	}
	opc.Busy = false
	opc.status = typex.DEV_DOWN
	return opc
}

// 初始化配置文件
func (sd *genericOpcuaDevice) Init(devId string, configMap map[string]interface{}) error {
	sd.PointId = devId
	if err := utils.BindSourceConfig(configMap, &sd.mainConfig); err != nil {
		return err
	}
	return nil
}

func (opcdev *genericOpcuaDevice) Start(cctx typex.CCTX) error {
	opcdev.Ctx = cctx.Ctx
	opcdev.CancelCTX = cctx.CancelCTX
	// 新建OPC UA 客户端
	endpoints, err := opcua.GetEndpoints(cctx.Ctx, opcdev.mainConfig.OpcuaCommonconfig.Endpoint)
	if err != nil {
		glogger.GLogger.Error("create opcua client failed:", err)
		return err
	}
	ep := opcua.SelectEndpoint(endpoints, opcdev.mainConfig.OpcuaCommonconfig.Policy, ua.MessageSecurityModeFromString(opcdev.mainConfig.OpcuaCommonconfig.Mode))
	if ep == nil {
		glogger.GLogger.Error("Setting opcua client failed:", err)
		return err
	}
	//初始化配置
	opts := []opcua.Option{
		opcua.SecurityPolicy(opcdev.mainConfig.OpcuaCommonconfig.Policy),
		opcua.SecurityModeString(opcdev.mainConfig.OpcuaCommonconfig.Mode),
	}
	//判断登录方式
	switch AuthType(opcdev.mainConfig.OpcuaCommonconfig.Auth) {
	case AUTH_USERNAME:
		opts = append(opts, opcua.AuthUsername(opcdev.mainConfig.OpcuaCommonconfig.Username, opcdev.mainConfig.OpcuaCommonconfig.Password))
		opts = append(opts, opcua.SecurityFromEndpoint(ep, ua.UserTokenTypeUserName))
		break
	default:
		opts = append(opts, opcua.AuthAnonymous())
		opts = append(opts, opcua.SecurityFromEndpoint(ep, ua.UserTokenTypeAnonymous))
		break
	}
	//连接opcua -判断连接是否正常
	opcdev.client = opcua.NewClient(ep.EndpointURL, opcua.SecurityMode(ua.MessageSecurityModeNone))
	if err := opcdev.client.Connect(cctx.Ctx); err != nil {
		glogger.GLogger.Error("Connect opcua client failed:", err)
		return err
	}
	opcua.RequestTimeout(time.Duration(opcdev.mainConfig.OpcuaCommonconfig.Timeout) * time.Millisecond)
	// 起一个线程去判断是否要轮询
	go func(ctx context.Context, Driver typex.XExternalDriver) {
		ticker := time.NewTicker(time.Duration(opcdev.mainConfig.OpcuaCommonconfig.Frequency) * time.Millisecond)
		buffer := make([]byte, common.T_64KB)
		for {
			<-ticker.C
			select {
			case <-ctx.Done():
				{
					ticker.Stop()
					return
				}
			default:
				{
				}
			}
			//----------------------------------------------------------------------------------
			if opcdev.Busy {
				continue
			}

			opcdev.Busy = true
			opcdev.locker.Lock()
			n, err := opcdev.readNodes([]byte{}, buffer)
			opcdev.locker.Unlock()
			if err != nil {
				glogger.GLogger.Error(err)
			} else {
				//周期轮询node数据
				opcdev.RuleEngine.WorkDevice(opcdev.Details(), string(buffer[:n]))
			}
			opcdev.Busy = false
		}
	}(cctx.Ctx, opcdev.driver)

	opcdev.status = typex.DEV_UP
	return nil
}
func (opcdev *genericOpcuaDevice) OnRead(cmd []byte, data []byte) (int, error) {

	n, err := opcdev.readNodes(cmd, data)
	if err != nil {
		glogger.GLogger.Error(err)
		opcdev.status = typex.DEV_DOWN
	}
	return n, err
}
func (opcdev *genericOpcuaDevice) readNodes(cmd []byte, data []byte) (int, error) {
	dataMap := map[string]OpcuaNode{}
	//遍历所有的寄存器
	for _, r := range opcdev.mainConfig.Opcnodes {
		// 设置一个间隔时间防止低级CPU黏包等
		time.Sleep(time.Duration(100) * time.Millisecond)
		id, err := ua.ParseNodeID(r.NodeID)
		if err != nil {
			glogger.GLogger.Error("invalid node id: %v", err)
		}
		req := &ua.ReadRequest{
			MaxAge: 2000,
			NodesToRead: []*ua.ReadValueID{
				{NodeID: id},
			},
			TimestampsToReturn: ua.TimestampsToReturnBoth,
		}
		ctx := context.Background()
		resp, err := opcdev.client.ReadWithContext(ctx, req)
		if err != nil {
			opcdev.errorCount++
			glogger.GLogger.Error("Read failed: %s", err)
		}
		if resp.Results[0].Status != ua.StatusOK {
			opcdev.errorCount++
			glogger.GLogger.Error("Status not OK: %v", resp.Results[0].Status)

		}
		value := OpcuaNode{
			Tag:         r.Tag,
			NodeID:      r.NodeID,
			Description: r.Description,
			DataType:    r.DataType,
			Value:       "",
		}
		value.Value, err = interfaceToString(resp.Results[0].Value.Value())
		dataMap[r.Tag] = value
		if err != nil {
			opcdev.errorCount++
			glogger.GLogger.Error("OPCUA value not match type: %s", err)
		}

	}
	bytes, _ := json.Marshal(dataMap)

	copy(data, bytes)
	return len(bytes), nil

}
func interfaceToString(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case int8:
		return strconv.FormatInt(int64(v), 10), nil
	case int16:
		return strconv.FormatInt(int64(v), 10), nil
	case int32:
		return strconv.FormatInt(int64(v), 10), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint64:
		return strconv.FormatUint(v, 10), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	default:
		return "", fmt.Errorf("unsupported type: %T", value)
	}
}

func (sd *genericOpcuaDevice) OnWrite(cmd []byte, data []byte) (int, error) {
	n, err := sd.driver.Write(cmd, data)
	//opc 协议写数据 待完善
	if err != nil {
		glogger.GLogger.Error(err)
		sd.status = typex.DEV_DOWN
	}
	return n, err
}

// 设备当前状态
func (sd *genericOpcuaDevice) Status() typex.DeviceState {
	if sd.mainConfig.OpcuaCommonconfig.RetryTime == 0 {
		sd.status = typex.DEV_UP
	}
	if sd.mainConfig.OpcuaCommonconfig.RetryTime > 0 {
		if sd.errorCount >= sd.mainConfig.OpcuaCommonconfig.RetryTime {
			sd.status = typex.DEV_DOWN
		}
	}
	return sd.status
}

// 停止设备
func (sd *genericOpcuaDevice) Stop() {
	sd.status = typex.DEV_STOP
	sd.CancelCTX()
	if sd.driver != nil {
		sd.driver.Stop()
		//释放连接
		sd.client.CloseWithContext(sd.Ctx)
		sd.driver = nil
	}
}

// 设备属性，是一系列属性描述
func (sd *genericOpcuaDevice) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (sd *genericOpcuaDevice) Details() *typex.Device {
	return sd.RuleEngine.GetDevice(sd.PointId)
}

// 状态
func (sd *genericOpcuaDevice) SetState(status typex.DeviceState) {
	sd.status = status

}

// 驱动
func (sd *genericOpcuaDevice) Driver() typex.XExternalDriver {
	return sd.driver
}

func (sd *genericOpcuaDevice) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
