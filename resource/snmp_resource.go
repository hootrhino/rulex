package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"rulex/typex"
	"rulex/utils"
	"strings"
	"sync"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/ngaut/log"
)

//----------------------------------------------------------------------------------

type SNMPResource struct {
	sync.Mutex
	typex.XStatus
	snmpClients []*gosnmp.GoSNMP
}

func (s *SNMPResource) GetClient(i int) *gosnmp.GoSNMP {
	return s.snmpClients[i]
}
func (s *SNMPResource) SetClient(i int, c *gosnmp.GoSNMP) {
	s.snmpClients[i] = c
}

func (s *SNMPResource) SystemDescrption(i int) string {
	r := ""
	s.GetClient(i).Walk(".1.3.6.1.2.1.1.1.0", func(variable gosnmp.SnmpPDU) error {
		if variable.Type == gosnmp.OctetString {
			r = string(variable.Value.([]byte))
		}
		return nil
	})
	return r
}
func (s *SNMPResource) PCName(i int) string {
	r := ""
	s.GetClient(i).Walk(".1.3.6.1.2.1.1.5.0", func(variable gosnmp.SnmpPDU) error {
		if variable.Type == gosnmp.OctetString {
			r = string(variable.Value.([]byte))
		}
		return nil
	})
	return r
}
func (s *SNMPResource) TotalMemory(i int) int {
	v := 0
	s.GetClient(i).Walk(".1.3.6.1.2.1.25.2.2.0", func(variable gosnmp.SnmpPDU) error {
		if variable.Type == gosnmp.Integer {
			v = int(variable.Value.(int))
		}
		return nil
	})
	return v

}
func (s *SNMPResource) CPUs(i int) map[string]int {
	oid := ".1.3.6.1.2.1.25.3.3.1.2"
	r := map[string]int{}
	s.GetClient(i).Walk(oid, func(variable gosnmp.SnmpPDU) error {
		if variable.Type == gosnmp.Integer {
			k := strings.Replace(variable.Name, ".1.3.6.1.2.1.25.3.3.1.2.", "", 1)
			r[k] = variable.Value.(int)
		}
		return nil
	})
	return r
}
func (s *SNMPResource) ProcessList(i int) []string {
	ss := []string{}
	s.GetClient(i).Walk(".1.3.6.1.2.1.25.4.2.1.2", func(variable gosnmp.SnmpPDU) error {
		if variable.Type == gosnmp.OctetString {
			ss = append(ss, string(variable.Value.([]byte)))
		}
		return nil
	})

	return ss
}
func (s *SNMPResource) InterfaceIPs(i int) []string {
	oid := "1.3.6.1.2.1.4.20.1.2"
	r := []string{}
	s.GetClient(i).Walk(oid, func(variable gosnmp.SnmpPDU) error {
		if variable.Type == gosnmp.Integer {
			ip := strings.Replace(variable.Name, ".1.3.6.1.2.1.4.20.1.2.", "", 1)
			if ip != "127.0.0.1" {
				r = append(r, ip)
			}
		}
		return nil
	})
	return r
}

func (s *SNMPResource) HardwareNetInterfaceMac(i int) []string {
	oid := ".1.3.6.1.2.1.2.2.1.6"
	maps := map[string]string{}
	s.GetClient(i).Walk(oid, func(variable gosnmp.SnmpPDU) error {
		if variable.Type == gosnmp.OctetString {
			macByte := variable.Value.([]byte)
			if len(macByte) == 6 {
				mac := fmt.Sprintf("%0x-%0x-%0x-%0x-%0x-%0x", macByte[0], macByte[1], macByte[2], macByte[3], macByte[4], macByte[5])
				maps[mac] = ""
			}
		}
		return nil
	})
	result := make([]string, 0)
	for k := range maps {
		result = append(result, k)
	}
	return result
}

//----------------------------------------------------------------------------------
type target struct {
	Target     string             `json:"target" validate:"required"`
	Port       uint16             `json:"port" validate:"required"`
	Transport  string             `json:"transport" validate:"required"`
	Community  string             `json:"community" validate:"required"`
	Version    uint8              `json:"version" validate:"required"`
	DataModels []typex.XDataModel `json:"dataModels" validate:"required"`
}

// GoSNMP represents GoSNMP library state.
type SNMPConfig struct {
	Frequency int64    `json:"frequency" validate:"required,gte=1,lte=10000"`
	Targets   []target `json:"targets" validate:"required"`
}

//--------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------

func NewSNMPInEndResource(inEndId string, e typex.RuleX) *SNMPResource {
	s := SNMPResource{}
	s.RuleEngine = e
	s.PointId = inEndId
	return &s
}

func (s *SNMPResource) Test(inEndId string) bool {
	r := []bool{}
	for i := 0; i < len(s.snmpClients); i++ {
		if err := s.GetClient(i).Connect(); err != nil {
			log.Errorf("SnmpClient [%v] Connect err: %v", s.GetClient(i).Target, err)
		} else {
			r = append(r, true)
		}
	}
	return len(r) == len(s.snmpClients)

}

func (s *SNMPResource) Register(inEndId string) error {
	s.PointId = inEndId
	return nil
}

func (s *SNMPResource) Start() error {
	config := s.RuleEngine.GetInEnd(s.PointId).Config
	mainConfig := SNMPConfig{}
	if err := utils.BindResourceConfig(config, &mainConfig); err != nil {
		return err
	}
	s.snmpClients = make([]*gosnmp.GoSNMP, len(mainConfig.Targets))
	for i, target := range mainConfig.Targets {
		s.SetClient(i, gosnmp.Default)
		s.GetClient(i).Target = target.Target
		s.GetClient(i).Community = target.Community

		if err := s.GetClient(i).Connect(); err != nil {
			log.Errorf("SnmpClient Connect err: %v", err)
			return err
		}

		go func(ctx context.Context, idx int) {
			log.Info("SnmpClient start working:", s.GetClient(i).Target)
			ticker := time.NewTicker(time.Duration(mainConfig.Frequency) * time.Second)
			for {
				select {
				case t := <-ticker.C:
					data := map[string]interface{}{
						"cpus":        s.CPUs(idx),
						"netsMac":     s.HardwareNetInterfaceMac(idx),
						"memory":      s.TotalMemory(idx),
						"ips":         s.InterfaceIPs(idx),
						"name":        s.PCName(idx),
						"description": s.SystemDescrption(idx),
					}
					dataBytes, _ := json.Marshal(data)
					if err0 := s.RuleEngine.PushQueue(typex.QueueData{
						In:   s.Details(),
						Out:  nil,
						E:    s.RuleEngine,
						Data: string(dataBytes),
					}); err0 != nil {
						log.Error("SNMPResource PushQueue error: ", err0, t)
					}
				default:
					{
					}
				}
			}
		}(context.Background(), i)
		log.Info("SNMPResource start successfully!")
	}

	return nil
}

func (s *SNMPResource) Enabled() bool {
	return s.Enable
}

func (s *SNMPResource) Details() *typex.InEnd {
	return s.RuleEngine.GetInEnd(s.PointId)
}

func (s *SNMPResource) DataModels() []typex.XDataModel {
	return []typex.XDataModel{}
}

func (s *SNMPResource) Reload() {

}

func (s *SNMPResource) Pause() {

}

func (s *SNMPResource) Status() typex.ResourceState {
	r := []bool{}
	for i := 0; i < len(s.snmpClients); i++ {
		if err := s.GetClient(i).Connect(); err != nil {
			log.Errorf("SnmpClient [%v] Connect err: %v", s.GetClient(i).Target, err)
		} else {
			r = append(r, true)
		}
	}

	if len(r) == len(s.snmpClients) {
		return typex.UP
	} else {
		return typex.DOWN
	}
}

func (s *SNMPResource) OnStreamApproached(data string) error {
	return nil
}

func (s *SNMPResource) Stop() {

}