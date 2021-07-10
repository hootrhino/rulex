package x

import (
	"time"

	"github.com/ngaut/log"
	"github.com/tarm/serial"
)

type SerialResource struct {
	*XStatus
	serialPort *serial.Port
}

func NewSerialResource(inEndId string) *SerialResource {
	s := SerialResource{}
	s.InEndId = inEndId
	//
	return &s
}

func (s *SerialResource) Test(inEndId string) bool {
	return true
}

func (s *SerialResource) Register(inEndId string) error {
	return nil
}

func (s *SerialResource) Start(e *RuleEngine) error {
	config := e.GetInEnd(s.InEndId).Config
	name := (*config)["port"]
	baud := (*config)["baud"]
	readTimeout := (*config)["read_timeout"]
	size := (*config)["size"]
	parity := (*config)["parity"]
	stopbits := (*config)["stopbits"]

	serialPort, err := serial.OpenPort(&serial.Config{
		Name:        name.(string),
		Baud:        baud.(int),
		ReadTimeout: time.Duration(readTimeout.(int64)),
		Size:        size.(byte),
		Parity:      serial.Parity(parity.(int)),
		StopBits:    serial.StopBits(stopbits.(int)),
	})
	if err != nil {
		log.Error("SerialResource start failed:", err)
		return err
	} else {
		s.serialPort = serialPort
		log.Info("SerialResource start success:")
		return nil
	}
}

func (s *SerialResource) Enabled() bool {
	return true
}

func (s *SerialResource) Reload() {
}

func (s *SerialResource) Pause() {

}

func (s *SerialResource) Status(e *RuleEngine) State {
	return UP
}

func (s *SerialResource) Stop() {
}
