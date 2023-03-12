package device

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"
	serial "github.com/tarm/serial"
)

// 传输形式：
// `rawtcp`, `rawudp`, `rs485rawserial`, `rs485rawtcp`
// const rawtcp string = "rawtcp"
// const rawudp string = "rawudp"
// const rs485rawserial string = "rs485rawserial"
// const rs485rawtcp string = "rs485rawtcp"

type _CommonConfig struct {
	Transport string `json:"transport" validate:"required"` // 传输协议
	WaitTime  int    `json:"waitTime" validate:"required"`  // 单个轮询间隔
}
type _UartConfig struct {
	Timeout  int    `json:"timeout" validate:"required"`
	Uart     string `json:"uart" validate:"required"`
	BaudRate int    `json:"baudRate" validate:"required"`
	DataBits int    `json:"dataBits" validate:"required"`
	Parity   string `json:"parity" validate:"required"`
	StopBits int    `json:"stopBits" validate:"required"`
}
type _ProtocolArg struct {
	In  string `json:"in" validate:"required"` // 十六进制字符串
	Out string `json:"out"`                    // 十六进制字符串
}
type _Protocol struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	RW          int    `json:"rw" validate:"required"`         // 1:RO 2:WO 3:RW
	BufferSize  int    `json:"bufferSize" validate:"required"` // 缓冲区大小
	Timeout     int    `json:"timeout" validate:"required"`    // 指令的等待时间, 在 Timeout 范围读 BufferSize 个字节, 否则就直接失败
	//---------------------------------------------------------------------
	// 下面都是校验算法相关配置:
	// -- 例如对[Byte1,Byte2,Byte3,Byte4,Byte5,Byte6,Byte7]用XOR算法比对
	//    从第一个开始，第五个结束[Byte1,Byte2,Byte3,Byte4,Byte5], 比对值位置在第六个[Byte6]
	// 伪代码：XOR(Byte[ChecksumBegin:ChecksumEnd]) == Byte[ChecksumValuePos]
	//---------------------------------------------------------------------
	CheckAlgorithm   string `json:"checkAlgorithm" validate:"required"`   // 校验算法，目前暂时支持: CRC16, XOR
	ChecksumValuePos uint   `json:"checksumValuePos" validate:"required"` // 校验值比对位
	ChecksumBegin    uint   `json:"checksumBegin" validate:"required"`    // 校验算法起始位置
	ChecksumEnd      uint   `json:"checksumEnd" validate:"required"`      // 校验算法结束位置
	AutoRequest      bool   `json:"autoRequest" validate:"required"`      // 是否开启轮询
	AutoRequestGap   uint   `json:"autoRequestGap" validate:"required"`   // 轮询间隔
	//---------------------------------------------------------------------
	ProtocolArg _ProtocolArg `json:"protocol" validate:"required"` // 参数
}

/*
*
* 自定义协议
*
 */
type _CustomProtocolConfig struct {
	CommonConfig _CommonConfig        `json:"commonConfig" validate:"required"`
	UartConfig   _UartConfig          `json:"uartConfig" validate:"required"`
	DeviceConfig map[string]_Protocol `json:"deviceConfig" validate:"required"`
}
type CustomProtocolDevice struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	serialPort *serial.Port // 现阶段暂时支持串口
	// tcpConn    *net.TCPConn // rawtcp 以后支持
	// udpConn    *net.UDPConn // rawudp 以后支持
	mainConfig _CustomProtocolConfig
	locker     sync.Locker
	errorCount int // 记录最大容错数，默认5次，出错超过5此就重启
}

func NewCustomProtocolDevice(e typex.RuleX) typex.XDevice {
	mdev := new(CustomProtocolDevice)
	mdev.RuleEngine = e
	mdev.locker = &sync.Mutex{}
	mdev.mainConfig = _CustomProtocolConfig{
		CommonConfig: _CommonConfig{},
		UartConfig:   _UartConfig{},
		DeviceConfig: map[string]_Protocol{},
	}
	mdev.status = typex.DEV_DOWN
	mdev.errorCount = 0
	return mdev

}

// 初始化
func (mdev *CustomProtocolDevice) Init(devId string, configMap map[string]interface{}) error {
	mdev.PointId = devId
	if err := utils.BindSourceConfig(configMap, &mdev.mainConfig); err != nil {
		return err
	}
	if !contains([]string{"N", "E", "O"}, mdev.mainConfig.UartConfig.Parity) {
		return errors.New("parity value only one of 'N','O','E'")
	}
	if !contains([]string{`rawtcp`, `rawudp`, `rs485rawserial`, `rs485rawtcp`},
		mdev.mainConfig.CommonConfig.Transport) {
		return errors.New("option only one of 'rawtcp','rawudp','rs485rawserial','rs485rawserial'")
	}
	// parse hex format
	for _, v := range mdev.mainConfig.DeviceConfig {
		// 检查指令是否符合十六进制
		if _, err := hex.DecodeString(v.ProtocolArg.In); err != nil {
			errMsg := fmt.Sprintf("invalid hex format:%s", v.ProtocolArg.In)
			glogger.GLogger.Error(errMsg)
			return fmt.Errorf(errMsg)
		}
		if v.ProtocolArg.Out != "" {
			if _, err := hex.DecodeString(v.ProtocolArg.Out); err != nil {
				errMsg := fmt.Sprintf("invalid hex format:%s", v.ProtocolArg.Out)
				glogger.GLogger.Error(errMsg)

				return fmt.Errorf(errMsg)
			}
		}
		// 目前暂时就先支持这几个算法
		if !contains([]string{"XOR", "xor", "CRC16", "crc16",
			"CRC32", "crc32", "NONECHECK"}, v.CheckAlgorithm) {
			return errors.New("unsupported check algorithm")
		}
		//------------------------------------------------------------------------------------------
		// 校验参数检查
		//------------------------------------------------------------------------------------------
		// 1. 检查区间是否越界
		if v.ChecksumBegin+v.ChecksumEnd > uint(v.BufferSize) {
			errMsg := fmt.Sprintf("check size [%d] out of buffer range:%v",
				v.ChecksumEnd, v.BufferSize)
			glogger.GLogger.Error(errMsg)
			return fmt.Errorf(errMsg)
		}
		// 2. 校验位是否超出缓冲长度
		if v.ChecksumValuePos > uint(v.BufferSize) {
			errMsg := fmt.Sprintf("checksum position [%d] out of buffer range:%v",
				v.ChecksumEnd, v.BufferSize)
			glogger.GLogger.Error(errMsg)
			return fmt.Errorf(errMsg)
		}

	}
	return nil
}

// 启动
func (mdev *CustomProtocolDevice) Start(cctx typex.CCTX) error {
	mdev.Ctx = cctx.Ctx
	mdev.CancelCTX = cctx.CancelCTX
	// 现阶段暂时只支持RS485串口, 以后有需求再支持TCP、UDP
	if mdev.mainConfig.CommonConfig.Transport == "rs485rawserial" {
		config := serial.Config{
			Name:     mdev.mainConfig.UartConfig.Uart,
			Baud:     mdev.mainConfig.UartConfig.BaudRate,
			Size:     byte(mdev.mainConfig.UartConfig.DataBits),
			Parity:   serial.Parity(mdev.mainConfig.UartConfig.Parity[0]),
			StopBits: serial.StopBits(mdev.mainConfig.UartConfig.StopBits),
		}
		serialPort, err := serial.OpenPort(&config)
		if err != nil {
			glogger.GLogger.Error("serialPort start failed:", err)
			return err
		}
		mdev.serialPort = serialPort
		// 起一个线程去判断是否要轮询
		go func(ctx context.Context, pp map[string]_Protocol) {
			result := [100]byte{} // 全局buf, 默认是100字节, 应该能覆盖绝大多数报文了
			for {
				select {
				case <-ctx.Done():
					return
				default:
					{
					}
				}
				//----------------------------------------------------------------------------------
				for _, p := range pp {
					if !p.AutoRequest {
						continue
					}
					hexs, err0 := hex.DecodeString(p.ProtocolArg.In)
					if err0 != nil {
						glogger.GLogger.Error(err0)
						mdev.errorCount++
						continue
					}
					if core.GlobalConfig.AppDebugMode {
						log.Println("[AppDebugMode] Write data:", hexs)
					}
					if mdev.serialPort == nil {
						mdev.status = typex.DEV_DOWN
						return
					}
					if _, err1 := mdev.serialPort.Write(hexs); err1 != nil {
						glogger.GLogger.Error("mdev.serialPort.Write error: ", err1)
						mdev.errorCount++
						continue
					}
					// 协议等待响应时间毫秒
					time.Sleep(time.Duration(p.AutoRequestGap) * time.Millisecond)
					if _, err2 := io.ReadAtLeast(mdev.serialPort, result[:p.BufferSize],
						p.BufferSize); err2 != nil {
						glogger.GLogger.Error("serialPort.ReadAtLeast error: ", err2)
						continue
					}
					if core.GlobalConfig.AppDebugMode {
						log.Println("[AppDebugMode] Write data:", p.ProtocolArg.In)
						log.Println("[AppDebugMode] Read data:", result[:p.BufferSize])
					}
					dataMap := map[string]string{}
					checkOk := false
					if p.CheckAlgorithm == "CRC16" || p.CheckAlgorithm == "crc16" {
						glogger.GLogger.Debug("checkCRC:", result[:p.BufferSize],
							int(result[:p.BufferSize][p.ChecksumValuePos]))
						checkOk = mdev.checkCRC(result[:p.BufferSize],
							int(result[:p.BufferSize][p.ChecksumValuePos]))

					}
					if p.CheckAlgorithm == "XOR" || p.CheckAlgorithm == "xor" {
						glogger.GLogger.Debug("checkXOR:", result[:p.BufferSize],
							int(result[:p.BufferSize][p.ChecksumValuePos]))
						checkOk = mdev.checkXOR(result[:p.BufferSize],
							int(result[:p.BufferSize][p.ChecksumValuePos]))
					}
					// NOCHECK: 不校验
					if p.CheckAlgorithm == "NOCHECK" {
						checkOk = true
					}
					if checkOk {
						// 返回给lua参数是十六进制大写字符串
						dataMap["name"] = p.Name
						dataMap["in"] = p.ProtocolArg.In
						dataMap["out"] = hex.EncodeToString(result[:p.BufferSize])
						bytes, _ := json.Marshal(dataMap)
						// 返回是十六进制大写字符串
						mdev.RuleEngine.WorkDevice(mdev.Details(), string(bytes))
					}
				}
				time.Sleep(time.Duration(mdev.mainConfig.CommonConfig.WaitTime) * time.Millisecond)
			}
		}(mdev.Ctx, mdev.mainConfig.DeviceConfig)
		mdev.status = typex.DEV_UP
		return nil
	}

	return fmt.Errorf("unsupported transport:%s", mdev.mainConfig.CommonConfig.Transport)
}

/*
*
* 数据读出来，对数据结构有要求, 其中Key必须是个数字或者数字字符串, 例如 1 or "1"
*
 */
func (mdev *CustomProtocolDevice) OnRead(cmd int, data []byte) (int, error) {
	// 拿到命令的索引
	p, exists := mdev.mainConfig.DeviceConfig[fmt.Sprintf("%d", cmd)]
	if exists {
		mdev.locker.Lock()
		hexs, err0 := hex.DecodeString(p.ProtocolArg.In)
		if err0 != nil {
			glogger.GLogger.Error(err0)
			mdev.errorCount++
			return 0, err0
		}
		if _, err1 := mdev.serialPort.Write(hexs); err1 != nil {
			glogger.GLogger.Error("serialPort.Write error: ", err1)
			mdev.errorCount++
			return 0, err1
		}
		mdev.locker.Unlock()

		// 协议等待响应时间毫秒
		time.Sleep(time.Duration(p.AutoRequestGap) * time.Millisecond)
		result := [100]byte{} // 全局buf, 默认是100字节, 应该能覆盖绝大多数报文了
		mdev.locker.Lock()

		if _, err2 := io.ReadAtLeast(mdev.serialPort, result[:p.BufferSize],
			p.BufferSize); err2 != nil {
			glogger.GLogger.Error("serialPort.ReadAtLeast error: ", err2)
			return 0, err2
		}
		mdev.locker.Unlock()

		if core.GlobalConfig.AppDebugMode {
			log.Println("[AppDebugMode] Write data:", p.ProtocolArg.In)
			log.Println("[AppDebugMode] Read data:", result[:p.BufferSize])
		}
		// 返回值
		dataMap := map[string]string{}
		checkOk := false
		if p.CheckAlgorithm == "CRC16" || p.CheckAlgorithm == "crc16" {
			glogger.GLogger.Debug("checkCRC:", result[:p.BufferSize],
				int(result[:p.BufferSize][p.ChecksumValuePos]))
			checkOk = mdev.checkCRC(result[:p.BufferSize],
				int(result[:p.BufferSize][p.ChecksumValuePos]))
		}
		//
		if p.CheckAlgorithm == "XOR" || p.CheckAlgorithm == "xor" {
			glogger.GLogger.Debug("checkCRC:", result[:p.BufferSize],
				int(result[:p.BufferSize][p.ChecksumValuePos]))
			checkOk = mdev.checkCRC(result[:p.BufferSize],
				int(result[:p.BufferSize][p.ChecksumValuePos]))
		}
		// NONECHECK: 不校验
		if p.CheckAlgorithm == "NONECHECK" {
			checkOk = true
		}
		if checkOk {
			// 返回给lua参数是十六进制大写字符串
			dataMap["name"] = p.Name
			dataMap["in"] = p.ProtocolArg.In
			dataMap["out"] = hex.EncodeToString(result[:p.BufferSize])
			bytes, _ := json.Marshal(dataMap)
			// 返回是十六进制大写字符串
			copy(data, bytes)
			return len(bytes), nil
		}
	}
	return 0, errors.New("unknown read command")

}

/*
*
* 写进来的数据格式 参考@Protocol
*
 */
type writeProtocol struct {
	Name             string `json:"name" validate:"required"`
	BufferSize       int    `json:"bufferSize" validate:"required"`
	TimeGap          int    `json:"timeGap" validate:"required"`
	CheckAlgorithm   string `json:"checkAlgorithm" validate:"required"`
	ChecksumValuePos uint   `json:"checksumValuePos" validate:"required"`
	ChecksumBegin    uint   `json:"checksumBegin" validate:"required"`
	ChecksumEnd      uint   `json:"checksumEnd" validate:"required"`
	In               string `json:"in" validate:"required"`
	Out              string `json:"out"`
}

// 把数据写入设备
// 根据第二个参数来找配置进去的自定义协议, 必须进来一个可识别的指令
// 其中cmd常为0,为无意义参数
func (mdev *CustomProtocolDevice) OnWrite(_ int, data []byte) (int, error) {

	return 0, errors.New("unknown write command")
}

// 设备当前状态
func (mdev *CustomProtocolDevice) Status() typex.DeviceState {
	if mdev.errorCount >= 5 {
		mdev.status = typex.DEV_DOWN
	}
	return mdev.status
}

// 停止设备
func (mdev *CustomProtocolDevice) Stop() {
	mdev.status = typex.DEV_STOP
	mdev.CancelCTX()
	if mdev.serialPort != nil {
		mdev.serialPort.Close()
		mdev.serialPort = nil
	}
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
* 设备服务调用，一般第三个参数为请求 body 传空
*
 */
func (mdev *CustomProtocolDevice) OnDCACall(_ string, Command string,
	Args interface{}) typex.DCAResult {
	T := reflect.TypeOf(Args)
	dcaResult := typex.DCAResult{Error: nil, Data: ""}
	if T.Name() != "[]interface{}" {
		dcaResult.Error = fmt.Errorf("error type:%s", T.Name())
		return dcaResult

	}
	wp := writeProtocol{CheckAlgorithm: "NONECHECK", TimeGap: 60}
	if err := json.Unmarshal([]byte((Args.([]string))[0]), &wp); err != nil {
		dcaResult.Error = err
		return dcaResult
	}
	mdev.locker.Lock()
	if _, err := mdev.serialPort.Write([]byte(wp.In)); err != nil {
		glogger.GLogger.Error("serialPort.Write error: ", err)
		mdev.errorCount++
		dcaResult.Error = err
		return dcaResult
	}
	mdev.locker.Unlock()
	time.Sleep(time.Duration(wp.TimeGap) * time.Millisecond)
	result := [100]byte{}
	//
	mdev.locker.Lock()
	if _, err := io.ReadAtLeast(mdev.serialPort, result[:wp.BufferSize],
		wp.BufferSize); err != nil {
		glogger.GLogger.Error("serialPort.ReadAtLeast error: ", err)
		dcaResult.Error = err
		return dcaResult
	}
	mdev.locker.Unlock()
	//
	if core.GlobalConfig.AppDebugMode {
		log.Println("[AppDebugMode] Write data:", wp.In)
		log.Println("[AppDebugMode] Read data:", result[:wp.BufferSize])
	}
	// 返回值
	dataMap := map[string]string{}
	checkOk := false
	if wp.CheckAlgorithm == "CRC16" || wp.CheckAlgorithm == "crc16" {
		glogger.GLogger.Debug("checkCRC:", result[:wp.BufferSize],
			int(result[:wp.BufferSize][wp.ChecksumValuePos]))
		checkOk = mdev.checkCRC(result[:wp.BufferSize],
			int(result[:wp.BufferSize][wp.ChecksumValuePos]))
	}
	//
	if wp.CheckAlgorithm == "XOR" || wp.CheckAlgorithm == "xor" {
		glogger.GLogger.Debug("checkCRC:", result[:wp.BufferSize],
			int(result[:wp.BufferSize][wp.ChecksumValuePos]))
		checkOk = mdev.checkCRC(result[:wp.BufferSize],
			int(result[:wp.BufferSize][wp.ChecksumValuePos]))
	}
	// NONECHECK: 不校验
	if wp.CheckAlgorithm == "NONECHECK" {
		checkOk = true
	}
	if checkOk {
		// 返回给lua参数是十六进制大写字符串
		dataMap["name"] = wp.Name
		dataMap["in"] = wp.In
		dataMap["out"] = hex.EncodeToString(result[:wp.BufferSize])
		bytes, _ := json.Marshal(dataMap)
		dcaResult.Data = string(bytes)
		return dcaResult
	} else {
		dcaResult.Error = fmt.Errorf("check failed")
	}
	return typex.DCAResult{}
}

// --------------------------------------------------------------------------------------------------
// 内部函数
// --------------------------------------------------------------------------------------------------
func (mdev *CustomProtocolDevice) checkXOR(b []byte, v int) bool {
	return utils.XOR(b) == v
}
func (mdev *CustomProtocolDevice) checkCRC(b []byte, v int) bool {

	return int(utils.CRC16(b)) == v
}
