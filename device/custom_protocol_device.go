package device

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"
	serial "github.com/tarm/serial"
)

// 读出来的字节缓冲默认大小
const __DEFAULT_BUFFER_SIZE = 100

// 传输形式：
// `rawtcp`, `rawudp`, `rs485rawserial`, `rs485rawtcp`
// const rawtcp string = "rawtcp"
// const rawudp string = "rawudp"
// const rs485rawserial string = "rs485rawserial"
// const rs485rawtcp string = "rs485rawtcp"

type _CommonConfig struct {
	Transport string `json:"transport" validate:"required"` // 传输协议
	Frequency int64  `json:"frequency" validate:"required" title:"采集频率" info:""`
	RetryTime int    `json:"retryTime" validate:"required"` // 几次以后重启,0 表示不重启
}
type _UartConfig struct {
	Timeout  int    `json:"timeout" validate:"required"`
	Uart     string `json:"uart" validate:"required"`
	BaudRate int    `json:"baudRate" validate:"required"`
	DataBits int    `json:"dataBits" validate:"required"`
	Parity   string `json:"parity" validate:"required"`
	StopBits int    `json:"stopBits" validate:"required"`
}

// Type=1
type _ProtocolArg struct {
	In  string `json:"in"`  // 十六进制字符串, 只有在静态协议下有用, 动态协议下就是""
	Out string `json:"out"` // 十六进制字符串, 用来存储返回值
}
type _Protocol struct {
	Name string `json:"name" validate:"required"` // 名称
	// 如果是静态的, 就取in参数; 如果是动态的, 则直接取第三个参数
	Type        int    `json:"type" validate:"required" default:"1"` // 指令类型, 1 静态, 2动态, 3 定时读, 4 定时读写
	Description string `json:"description"`                          // 描述文本
	RW          int    `json:"rw" validate:"required"`               // 1:RO 2:WO 3:RW
	BufferSize  int    `json:"bufferSize" validate:"required"`       // 缓冲区大小
	Timeout     int    `json:"timeout" validate:"required"`          // 指令的等待时间, 在 Timeout 范围读 BufferSize 个字节, 否则就直接失败
	// [Important!] 该参数用来配合定时协议使用, Type== 3、4 时生效
	TimeSlice int `json:"timeSlice" validate:"required"` // 定时请求倒计时,单位毫秒，默认为0
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
	OnCheckError     string `json:"onCheckError" default:"IGNORE"`                          // 当指令操作失败时动作: IGNORE, LOG
	//
	AutoRequest bool `json:"autoRequest" validate:"required"` // 是否开启轮询, 开启轮询后, 每次间隔时间为 Frequency 毫秒
	//---------------------------------------------------------------------
	// 只有在静态协议(Type=1)下有用
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
	errorCount int // 记录最大容错数，默认5次，出错超过5此就重启
}

func NewCustomProtocolDevice(e typex.RuleX) typex.XDevice {
	mdev := new(CustomProtocolDevice)
	mdev.RuleEngine = e
	mdev.mainConfig = _CustomProtocolConfig{
		CommonConfig: _CommonConfig{},
		UartConfig:   _UartConfig{},
		DeviceConfig: map[string]_Protocol{},
	}
	mdev.Busy = false
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

	// 频率不能太快
	if mdev.mainConfig.CommonConfig.Frequency < 50 {
		return errors.New("'frequency' must grate than 50 millisecond")
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
			ticker := time.NewTicker(time.Duration(mdev.mainConfig.CommonConfig.Frequency) * time.Millisecond)
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
				if mdev.serialPort == nil {
					mdev.status = typex.DEV_DOWN
					return
				}
				//----------------------------------------------------------------------------------
				if mdev.Busy {
					continue
				}
				mdev.Busy = true
				for _, p := range pp {
					if !p.AutoRequest {
						continue
					}
					// 1: 读
					if p.RW != 1 {
						continue
					}
					// 只针对静态协议
					if p.Type != 1 {
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
					if _, err1 := mdev.serialPort.Write(hexs); err1 != nil {
						glogger.GLogger.Error("serialPort.Write error: ", err1)
						mdev.errorCount++
						continue
					}
					result := [__DEFAULT_BUFFER_SIZE]byte{} // 全局buf, 默认是100字节, 应该能覆盖绝大多数报文了
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					if _, err2 := utils.ReadAtLeast(ctx, mdev.serialPort, result[:p.BufferSize],
						p.BufferSize); err2 != nil {
						glogger.GLogger.Error("serialPort.ReadAtLeast error: ", err2)
						mdev.errorCount++
						cancel()
						continue
					}
					cancel()
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
						mdev.RuleEngine.WorkDevice(mdev.Details(), string(bytes))
					} else {
						if p.OnCheckError == "LOG" {
							msg := "checkSum error, Algorithm:%s; Begin:%v; End:%v; CheckPos:%v;"
							glogger.GLogger.Error(msg,
								p.CheckAlgorithm,
								p.ChecksumBegin,
								p.ChecksumEnd,
								p.ChecksumValuePos)
							mdev.errorCount++
						}
					}
				}
				mdev.Busy = false
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
func (mdev *CustomProtocolDevice) OnRead(cmd []byte, data []byte) (int, error) {
	// 拿到命令的索引
	p, exists := mdev.mainConfig.DeviceConfig[string(cmd)]
	if exists {

		// 静态协议
		if p.Type == 1 {
			// 判断是不是读权限
			if p.RW != 1 {
				return 0, errors.New("RW permission deny")
			}
			hexs, err0 := hex.DecodeString(p.ProtocolArg.In)
			if err0 != nil {
				glogger.GLogger.Error(err0)
				mdev.errorCount++
				return 0, err0
			}

			_, err1 := mdev.serialPort.Write(hexs)
			if err1 != nil {
				glogger.GLogger.Error("serialPort.Write error: ", err1)
				mdev.errorCount++
				return 0, err1
			}

			result := [__DEFAULT_BUFFER_SIZE]byte{} // 全局buf, 默认是100字节, 应该能覆盖绝大多数报文了
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if _, err2 := utils.ReadAtLeast(ctx, mdev.serialPort, result[:p.BufferSize],
				p.BufferSize); err2 != nil {
				glogger.GLogger.Error("serialPort.ReadAtLeast error: ", err2)
				mdev.errorCount++
				cancel()
				return 0, err2
			}
			cancel()

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
			} else {
				if p.OnCheckError == "IGNORE" {
					// Do Nothing
					return 0, nil
				}
				if p.OnCheckError == "LOG" {
					msg := "checkSum error, Algorithm:%s; Begin:%v; End:%v; CheckPos:%v;"
					glogger.GLogger.Error(msg,
						p.CheckAlgorithm,
						p.ChecksumBegin,
						p.ChecksumEnd,
						p.ChecksumValuePos)
					mdev.errorCount++
					return 0, errors.New(msg)
				}
			}
		}
	}
	return 0, errors.New("unknown read command:" + string(cmd))

}

/*
*
* 写进来的数据格式 参考@Protocol
*
 */

// 把数据写入设备
func (mdev *CustomProtocolDevice) OnWrite(cmd []byte, data []byte) (int, error) {
	// 拿到命令的索引
	p, exists := mdev.mainConfig.DeviceConfig[string(cmd)]
	if exists {

		// 静态协议
		if p.Type == 1 {
			// 判断是不是读权限
			if p.RW != 1 {
				return 0, errors.New("RW permission deny")
			}
			hexs, err0 := hex.DecodeString(p.ProtocolArg.In)
			if err0 != nil {
				glogger.GLogger.Error(err0)
				mdev.errorCount++
				return 0, err0
			}

			_, err1 := mdev.serialPort.Write(hexs)
			if err1 != nil {
				glogger.GLogger.Error("serialPort.Write error: ", err1)
				mdev.errorCount++
				return 0, err1
			}

			result := [__DEFAULT_BUFFER_SIZE]byte{} // 全局buf, 默认是100字节, 应该能覆盖绝大多数报文了
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if _, err2 := utils.ReadAtLeast(ctx, mdev.serialPort, result[:p.BufferSize],
				p.BufferSize); err2 != nil {
				glogger.GLogger.Error("serialPort.ReadAtLeast error: ", err2)
				mdev.errorCount++
				cancel()
				return 0, err2
			}
			cancel()

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
			} else {
				if p.OnCheckError == "IGNORE" {
					// Do Nothing
					return 0, nil
				}
				if p.OnCheckError == "LOG" {
					msg := "checkSum error, Algorithm:%s; Begin:%v; End:%v; CheckPos:%v;"
					glogger.GLogger.Error(msg,
						p.CheckAlgorithm,
						p.ChecksumBegin,
						p.ChecksumEnd,
						p.ChecksumValuePos)
					mdev.errorCount++
					return 0, errors.New(msg)
				}

			}
		}

	}
	return 0, errors.New("unknown write command:" + string(cmd))
}

/*
*
* 外部指令交互, 常用来实现自定义协议等
*
 */
func (mdev *CustomProtocolDevice) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	// 拿到命令的索引
	p, exists := mdev.mainConfig.DeviceConfig[string(cmd)]
	if exists {
		// 动态协议
		// local err = applib:WriteDevice("UUID", "CMD", hex_string)
		// local err = applib:ReadDevice("UUID", "CMD")
		// 如果是动态协议, 不检查RW, 则取第三个参数，然后写入到设备, 第三个参数要求是十六进制字符串
		//
		// 实际上当类型为2的时候其他的参数都无意义了
		//
		if p.Type == 2 {
			glogger.GLogger.Debug("Dynamic protocol:", string(args))
			// 取data参数
			hexs, err := hex.DecodeString(string(args))
			if err != nil {
				glogger.GLogger.Error(err)
				return nil, err
			}
			_, err1 := mdev.serialPort.Write(hexs)
			if core.GlobalConfig.AppDebugMode {
				log.Println("[AppDebugMode] Write data:", p.ProtocolArg.In)
			}
			if err1 != nil {
				glogger.GLogger.Error("Dynamic protocol write error: ", err1)
				mdev.errorCount++
				return nil, err1
			}

			result := [__DEFAULT_BUFFER_SIZE]byte{} // 全局buf, 默认是100字节, 应该能覆盖绝大多数报文了
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if _, err2 := utils.ReadAtLeast(ctx, mdev.serialPort, result[:p.BufferSize],
				p.BufferSize); err2 != nil {
				glogger.GLogger.Error("serialPort.ReadAtLeast error: ", err2)
				mdev.errorCount++
				cancel()
				return nil, err2
			}
			cancel()
			if core.GlobalConfig.AppDebugMode {
				log.Println("[AppDebugMode] Read data:", result[:p.BufferSize])
			}
			// return
			dataMap := map[string]string{}
			dataMap["name"] = p.Name
			dataMap["in"] = string(args)
			dataMap["out"] = hex.EncodeToString(result[:p.BufferSize])
			bytes, _ := json.Marshal(dataMap)
			return (bytes), nil
		}
		//------------------------------------------------------------------------------------------
		// 基于时间片的轮询协议
		//------------------------------------------------------------------------------------------
		// 时间片只读
		if p.Type == 3 {
			glogger.GLogger.Debug("Time slice SliceReceive:", p.TimeSlice)
			result := [__DEFAULT_BUFFER_SIZE]byte{}
			count, err := utils.SliceReceive(context.Background(),
				mdev.serialPort, result[:], time.Duration(p.TimeSlice))
			return (result[:count]), err
		}
		// 时间片读写
		if p.Type == 4 {
			glogger.GLogger.Debug("Time slice SliceRequest:", string(args))
			hexs, err := hex.DecodeString(string(args))
			if err != nil {
				glogger.GLogger.Error(err)
				return nil, err
			}
			result := [__DEFAULT_BUFFER_SIZE]byte{}
			count, err := utils.SliceRequest(context.Background(),
				mdev.serialPort, hexs, result[:], time.Duration(p.TimeSlice))
			return (result[:count]), err
		}
	}
	return nil, errors.New("unknown ctrl command:" + string(cmd))
}

// 设备当前状态
func (mdev *CustomProtocolDevice) Status() typex.DeviceState {
	if mdev.mainConfig.CommonConfig.RetryTime == 0 {
		mdev.status = typex.DEV_UP
	}
	if mdev.mainConfig.CommonConfig.RetryTime > 0 {
		if mdev.errorCount >= mdev.mainConfig.CommonConfig.RetryTime {
			mdev.status = typex.DEV_DOWN
		}
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
func (mdev *CustomProtocolDevice) checkXOR(b []byte, v int) bool {
	return utils.XOR(b) == v
}
func (mdev *CustomProtocolDevice) checkCRC(b []byte, v int) bool {

	return int(utils.CRC16(b)) == v
}
