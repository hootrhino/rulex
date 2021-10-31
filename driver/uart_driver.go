package driver

import (
	"bytes"
	"context"
	"encoding/binary"
	"rulex/typex"
	"time"

	"github.com/goburrow/serial"
	"github.com/ngaut/log"
)

//------------------------------------------------------------------------
// 内部函数
//------------------------------------------------------------------------

//
// 正点原子的 Lora 模块封装
//
type UartDriver struct {
	serialPort serial.Port
	channel    chan byte
	ctx        context.Context
	In         *typex.InEnd
	RuleEngine typex.RuleX
}

//
// 初始化一个驱动
//
func NewUartDriver(serialPort serial.Port, in *typex.InEnd, e typex.RuleX) typex.XExternalDriver {
	m := &UartDriver{}
	// 缓冲区： 4KB
	m.channel = make(chan byte, 4096)
	m.In = in
	m.RuleEngine = e
	m.serialPort = serialPort
	m.ctx = context.Background()
	return m
}

//
//
//
func (a *UartDriver) Init() error {
	return nil
}
func (a *UartDriver) Work() error {

	go func(context.Context) {
		log.Debug("UartDriver Start Listening")
		for {
			time.Sleep(400 * time.Millisecond)
			data := make([]byte, 256) // byte
			size, err0 := a.serialPort.Read(data)
			if err0 != nil {
				log.Error("UartDriver error: ", err0)
				continue
			}
			for i := 0; i < size; i++ {
				a.channel <- data[i]
			}
			// 数据包头长度位 3 字节
			if len(a.channel) > 3 && len(a.channel) <= 256 {
				// 包头：
				// -----------------------------------
				// | 包长1 | 包长2 | 类型 | 数据······|
				// -===============+++++++############
				//
				dataBytes := [3]byte{}
				// 前两个字节保存数据长度，最大长256个字节
				dataBytes[0] = <-a.channel
				dataBytes[1] = <-a.channel
				// 第三个字节表示数据包类型
				dataBytes[2] = <-a.channel
				// log.Info(dataBytes)
				var dataLen uint16
				if err := binary.Read(bytes.NewReader([]byte{dataBytes[0], dataBytes[1]}), binary.BigEndian, &dataLen); err != nil {
					log.Error(err)
					continue
				}
				// 读数据包类型
				var dataType uint8
				if err := binary.Read(bytes.NewReader([]byte{dataBytes[2]}), binary.BigEndian, &dataType); err != nil {
					log.Error(err)
					continue
				}
				//
				// log.Infof("len(channel):%d  dataLen:%d Type is: %d", len(a.channel), dataLen, dataType)
				// 允许最大可读 256 字节
				if dataLen <= (256) && dataLen > 0 {
					var buffer = make([]byte, dataLen)
					// 当前的数据够不够读
					// log.Infof("len(channel):%d  dataLen is: %d", len(a.channel), (dataLen))
					if len(a.channel) >= int(dataLen) {
						for i := 0; i < int(dataLen); i++ {
							buffer = append(buffer, <-a.channel)
						}
						log.Info("SerialPort Received:", string(buffer))
						a.RuleEngine.PushQueue(typex.QueueData{
							In:   a.In,
							Out:  nil,
							E:    a.RuleEngine,
							Data: string(buffer),
						})
					}

				}

			}
		}

	}(a.ctx)
	return nil

}
func (a *UartDriver) State() typex.DriverState {
	return typex.RUNNING

}
func (a *UartDriver) Stop() error {
	a.ctx.Done()
	return a.serialPort.Close()
}

func (a *UartDriver) Test() error {
	return nil
}

//
func (a *UartDriver) Read([]byte) (int, error) {

	return 0, nil
}

//
func (a *UartDriver) Write(b []byte) (int, error) {
	n, err := a.serialPort.Write(b)
	if err != nil {
		log.Error(err)
		return 0, err
	} else {
		return n, nil
	}

}
