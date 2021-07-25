package core

import (
	"time"

	"github.com/ngaut/log"
	"github.com/tarm/serial"
)

type SerialResource struct {
	XStatus
	serialPort *serial.Port
}

func NewSerialResource(inEndId string, e *RuleEngine) *SerialResource {
	s := SerialResource{}
	s.PointId = inEndId
	s.ruleEngine = e
	//
	return &s
}

func (mm *SerialResource) DataModels() *map[string]XDataModel {
	return &map[string]XDataModel{}
}

func (s *SerialResource) Test(inEndId string) bool {
	return true
}

func (s *SerialResource) Register(inEndId string) error {
	return nil
}

func (s *SerialResource) Start() error {
	config := s.ruleEngine.GetInEnd(s.PointId).Config
	name := (*config)["name"]
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

func (s *SerialResource) Status() State {
	return UP
}

func (s *SerialResource) Stop() {
}
