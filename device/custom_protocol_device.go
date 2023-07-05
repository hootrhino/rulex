package device

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	serial "github.com/wwhai/tarmserial"
)

// 读出来的字节缓冲默认大小
const __DEFAULT_BUFFER_SIZE = 100

// 传输形式：
// `rawtcp`, `rawudp`, `rs485rawserial`, `rs485rawtcp`
// const rawtcp string = "rawtcp"
// const rawudp string = "rawudp"
// const rs485rawserial string = "rs485rawserial"
// const rs485rawtcp string = "rs485rawtcp"

type _CPDCommonConfig struct {
	Transport *string `json:"transport" validate:"required"` // 传输协议
	RetryTime *int    `json:"retryTime" validate:"required"` // 几次以后重启,0 表示不重启
}

type _CPDProtocol struct {
	//---------------------------------------------------------------------
	// 下面都是校验算法相关配置:
	// -- 例如对[Byte1,Byte2,Byte3,Byte4,Byte5,Byte6,Byte7]用XOR算法比对
	//    从第一个开始，第五个结束[Byte1,Byte2,Byte3,Byte4,Byte5], 比对值位置在第六个[Byte6]
	// 伪代码：XOR(Byte[ChecksumBegin:ChecksumEnd]) == Byte[ChecksumValuePos]
	//---------------------------------------------------------------------
	CheckAlgorithm   string `json:"checkAlgorithm" validate:"required" default:"NONECHECK"` // 校验算法，目前暂时支持: CRC16, XOR
	ChecksumValuePos uint   `json:"checksumValuePos" validate:"required"`                   // 校验值比对位
	ChecksumBegin    uint   `json:"checksumBegin" validate:"required"`                      // 校验算法起始位置
	ChecksumEnd      uint   `json:"checksumEnd" validate:"required"`                        // 校验算法结束位置
}

/*
*
* 自定义协议
*
 */
type _CustomProtocolConfig struct {
	CommonConfig _CPDCommonConfig        `json:"commonConfig" validate:"required"`
	UartConfig   common.CommonUartConfig `json:"uartConfig" validate:"required"`
}
type CustomProtocolDevice struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	serialPort *serial.Port // 现阶段暂时支持串口
	mainConfig _CustomProtocolConfig
	errorCount int // 记录最大容错数，默认5次，出错超过5此就重启
}

func NewCustomProtocolDevice(e typex.RuleX) typex.XDevice {
	mdev := new(CustomProtocolDevice)
	mdev.RuleEngine = e
	mdev.mainConfig = _CustomProtocolConfig{
		CommonConfig: _CPDCommonConfig{},
		UartConfig:   common.CommonUartConfig{},
	}
	return mdev

}

// 初始化
func (mdev *CustomProtocolDevice) Init(devId string, configMap map[string]interface{}) error {
	mdev.PointId = devId
	if err := utils.BindSourceConfig(configMap, &mdev.mainConfig); err != nil {
		return err
	}
	if !utils.SContains([]string{"N", "E", "O"}, mdev.mainConfig.UartConfig.Parity) {
		return errors.New("parity value only one of 'N','O','E'")
	}
	if !utils.SContains([]string{`rawtcp`, `rawudp`, `rs485rawserial`, `rs485rawtcp`},
		*mdev.mainConfig.CommonConfig.Transport) {
		return errors.New("option only one of 'rawtcp','rawudp','rs485rawserial','rs485rawserial'")
	}
	return nil
}

// 启动
func (mdev *CustomProtocolDevice) Start(cctx typex.CCTX) error {
	mdev.Ctx = cctx.Ctx
	mdev.CancelCTX = cctx.CancelCTX
	// 现阶段暂时只支持RS485串口, 以后有需求再支持TCP、UDP
	if *mdev.mainConfig.CommonConfig.Transport == "rs485rawserial" {
		config := serial.Config{
			Name:        mdev.mainConfig.UartConfig.Uart,
			Baud:        mdev.mainConfig.UartConfig.BaudRate,
			Size:        byte(mdev.mainConfig.UartConfig.DataBits),
			Parity:      serial.Parity(mdev.mainConfig.UartConfig.Parity[0]),
			StopBits:    serial.StopBits(mdev.mainConfig.UartConfig.StopBits),
			ReadTimeout: time.Duration(mdev.mainConfig.UartConfig.Timeout) * time.Millisecond,
		}
		serialPort, err := serial.OpenPort(&config)
		if err != nil {
			glogger.GLogger.Error("serialPort start failed:", err)
			return err
		}
		mdev.errorCount = 0
		mdev.serialPort = serialPort
		mdev.status = typex.DEV_UP
		return nil
	}

	return fmt.Errorf("unsupported transport:%s", *mdev.mainConfig.CommonConfig.Transport)
}

/*
*
* 数据读出来，对数据结构有要求, 其中Key必须是个数字或者数字字符串, 例如 1 or "1"
*
 */
func (mdev *CustomProtocolDevice) OnRead(cmd []byte, data []byte) (int, error) {
	return 0, errors.New("unknown read command:" + string(cmd))

}

/*
*
* 写进来的数据格式 参考@Protocol
*
 */

// 把数据写入设备
func (mdev *CustomProtocolDevice) OnWrite(cmd []byte, data []byte) (int, error) {
	return 0, errors.New("unknown write command:" + string(cmd))
}

/*
*
* 外部指令交互, 常用来实现自定义协议等
*
 */
func (mdev *CustomProtocolDevice) OnCtrl(cmd []byte, _ []byte) ([]byte, error) {
	glogger.GLogger.Debug("Time slice SliceRequest:", string(cmd))
	return mdev.ctrl(cmd)
}

// 设备当前状态
func (mdev *CustomProtocolDevice) Status() typex.DeviceState {
	if *mdev.mainConfig.CommonConfig.RetryTime == 0 {
		mdev.status = typex.DEV_UP
	}
	if *mdev.mainConfig.CommonConfig.RetryTime > 0 {
		if mdev.errorCount >= *mdev.mainConfig.CommonConfig.RetryTime {
			mdev.CancelCTX()
			mdev.status = typex.DEV_DOWN
		}
	}
	return mdev.status
}

// 停止设备
func (mdev *CustomProtocolDevice) Stop() {
	mdev.CancelCTX()
	mdev.status = typex.DEV_DOWN
	mdev.serialPort.Close()

}

// 设备属性，是一系列属性描述
func (mdev *CustomProtocolDevice) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (mdev *CustomProtocolDevice) Details() *typex.Device {
	return mdev.RuleEngine.GetDevice(mdev.PointId)
}

// 状态
func (mdev *CustomProtocolDevice) SetState(status typex.DeviceState) {
	mdev.status = status
}

// 驱动
func (mdev *CustomProtocolDevice) Driver() typex.XExternalDriver {
	return nil
}

/*
*
* 设备服务调用
*
 */
func (mdev *CustomProtocolDevice) OnDCACall(_ string, Command string,
	Args interface{}) typex.DCAResult {

	return typex.DCAResult{}
}

// --------------------------------------------------------------------------------------------------
// 内部函数
// --------------------------------------------------------------------------------------------------
func (mdev *CustomProtocolDevice) ctrl(args []byte) ([]byte, error) {
	glogger.GLogger.Debug("Custom Protocol Device Request:", string(args))
	hexs, err1 := hex.DecodeString(string(args))
	if err1 != nil {
		glogger.GLogger.Error(err1)
		return nil, err1
	}
	result := [__DEFAULT_BUFFER_SIZE]byte{}
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(mdev.mainConfig.UartConfig.Timeout)*time.Millisecond)
	count, err2 := utils.SliceRequest(ctx, mdev.serialPort,
		hexs, result[:], false,
		time.Duration(30)*time.Millisecond /*30ms wait*/)
	cancel()
	if err2 != nil {
		glogger.GLogger.Error("Custom Protocol Device Request error: ", err2)
		mdev.errorCount++
		return nil, err2
	}
	dataMap := map[string]string{}
	dataMap["in"] = string(args)
	dataMap["out"] = hex.EncodeToString(result[:count])
	bytes, _ := json.Marshal(dataMap)
	return []byte(bytes), nil
}

// func (mdev *CustomProtocolDevice) checkXOR(b []byte, v int) bool {
// 	return utils.XOR(b) == v
// }
// func (mdev *CustomProtocolDevice) checkCRC(b []byte, v int) bool {

// 	return int(utils.CRC16(b)) == v
// }

// /*
// *
// * Check hex string
// *
//  */
// func (mdev *CustomProtocolDevice) checkHexs(p _CPDProtocol, result []byte) bool {
// 	checkOk := false
// 	if p.CheckAlgorithm == "CRC16" || p.CheckAlgorithm == "crc16" {
// 		glogger.GLogger.Debug("checkCRC:", result[:p.BufferSize],
// 			int(result[:p.BufferSize][p.ChecksumValuePos]))
// 		checkOk = mdev.checkCRC(result[:p.BufferSize],
// 			int(result[:p.BufferSize][p.ChecksumValuePos]))
// 	}
// 	//
// 	if p.CheckAlgorithm == "XOR" || p.CheckAlgorithm == "xor" {
// 		glogger.GLogger.Debug("checkCRC:", result[:p.BufferSize],
// 			int(result[:p.BufferSize][p.ChecksumValuePos]))
// 		checkOk = mdev.checkXOR(result[:p.BufferSize],
// 			int(result[:p.BufferSize][p.ChecksumValuePos]))
// 	}
// 	// NONECHECK: 不校验
// 	if p.CheckAlgorithm == "NONECHECK" {
// 		checkOk = true
// 	}
// 	return checkOk
// }
