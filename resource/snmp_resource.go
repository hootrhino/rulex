package resource

//-----------------------------------------------------------------------------
//                              Warning
// github.com/gosnmp/gosnmp:
//    这个库有点问题, 效率不高, 因为大量用了循环导致监控点数量多了以后, 就会大量吃CPU
// 因此这个SNMP的监控[不推荐正式场景]使用, 后期可能会找个好点的库重写这个功能.
//
//------------------------------------------------------------------------------

import (
	"context"
	"encoding/json"
	"fmt"
	"rulex/core"
	"rulex/typex"
	"rulex/utils"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/ngaut/log"
)

//----------------------------------------------------------------------------------

type SNMPResource struct {
	typex.XStatus
	snmpClients []*gosnmp.GoSNMP
}

func (s *SNMPResource) GetClient(i int) *gosnmp.GoSNMP {
	return s.snmpClients[i]
}
func (s *SNMPResource) SetClient(i int, c *gosnmp.GoSNMP) {
	s.snmpClients[i] = c
}

func (s *SNMPResource) SystemInfo(i int) map[string]interface{} {
	results, err := s.GetClient(i).Get([]string{
		".1.3.6.1.2.1.1.1.0",    // 信息
		".1.3.6.1.2.1.1.5.0",    // PCName
		".1.3.6.1.2.1.25.2.2.0", // TotalMemory
	})
	if err != nil {
		log.Error(err)
		return map[string]interface{}{
			"info":        "",
			"pcName":      "",
			"totalMemory": 0,
		}
	}
	if len(results.Variables) == 3 {
		Info := string(results.Variables[0].Value.([]byte))
		PCName := string(results.Variables[1].Value.([]byte))
		TotalMemory := (results.Variables[2].Value.(int))
		return map[string]interface{}{
			"info":        Info,
			"pcName":      PCName,
			"totalMemory": TotalMemory,
		}
	} else {
		return map[string]interface{}{
			"info":        "",
			"pcName":      "",
			"totalMemory": 0,
		}
	}

}

func (s *SNMPResource) CPUs(i int) map[string]int {
	oid := ".1.3.6.1.2.1.25.3.3.1.2"
	r := map[string]int{}
	err := s.GetClient(i).Walk(oid, func(variable gosnmp.SnmpPDU) error {
		if variable.Type == gosnmp.Integer {
			k := strings.Replace(variable.Name, ".1.3.6.1.2.1.25.3.3.1.2.", "", 1)
			r[k] = variable.Value.(int)
		}
		return nil
	})
	if err != nil {
		log.Error(err)
		return r
	}
	return r
}
func (s *SNMPResource) InterfaceIPs(i int) []string {
	oid := "1.3.6.1.2.1.4.20.1.2"
	var r []string
	err := s.GetClient(i).Walk(oid, func(variable gosnmp.SnmpPDU) error {
		if variable.Type == gosnmp.Integer {
			ip := strings.Replace(variable.Name, ".1.3.6.1.2.1.4.20.1.2.", "", 1)
			if ip != "127.0.0.1" {
				r = append(r, ip)
			}
		}
		return nil
	})
	if err != nil {
		log.Error(err)
		return r
	}
	return r
}

func (s *SNMPResource) HardwareNetInterfaceMac(i int) []string {
	oid := ".1.3.6.1.2.1.2.2.1.6"
	maps := map[string]string{}
	result := make([]string, 0)

	err := s.GetClient(i).Walk(oid, func(variable gosnmp.SnmpPDU) error {
		if variable.Type == gosnmp.OctetString {
			macByte := variable.Value.([]byte)
			if len(macByte) == 6 {
				mac := fmt.Sprintf("%0x-%0x-%0x-%0x-%0x-%0x", macByte[0], macByte[1], macByte[2], macByte[3], macByte[4], macByte[5])
				maps[mac] = ""
			}
		}
		return nil
	})
	if err != nil {
		log.Error(err)
		return result
	}
	for k := range maps {
		result = append(result, k)
	}
	return result
}

//----------------------------------------------------------------------------------
type target struct {
	Target     string             `json:"target" validate:"required" title:"目标IP" info:""`
	Port       uint16             `json:"port" validate:"required" title:"目标端口" info:""`
	Transport  string             `json:"transport" validate:"required" title:"传输形式" info:""`
	Community  string             `json:"community" validate:"required" title:"社区名称" info:""`
	Version    uint8              `json:"version" validate:"required" title:"SNMP版本" info:""`
	DataModels []typex.XDataModel `json:"dataModels" validate:"required" title:"数据模型" info:""`
}

// SNMPConfig
// GoSNMP represents GoSNMP library state.
type SNMPConfig struct {
	Frequency int64    `json:"frequency" validate:"required" title:"采集频率" info:""`
	Timeout   int64    `json:"timeout" validate:"required" title:"超时时间" info:""`
	Targets   []target `json:"targets" validate:"required" title:"采集目标" info:""`
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
func (*SNMPResource) Driver() typex.XExternalDriver {
	return nil
}
func (s *SNMPResource) Test(inEndId string) bool {
	var r []bool
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
		s.GetClient(i).Timeout = time.Duration(time.Duration(mainConfig.Timeout) * time.Second)

		if err := s.GetClient(i).Connect(); err != nil {
			log.Errorf("SnmpClient Connect err: %v", err)
			return err
		}
		ticker := time.NewTicker(time.Duration(mainConfig.Frequency) * time.Second)
		go func(ctx context.Context, idx int, sr *SNMPResource) {
			for {
				select {
				case t := <-ticker.C:
					data := map[string]interface{}{
						"systemInfo": sr.SystemInfo(i), // Waining: CPU maybe used 100%
						"time":       t.Format("2006-01-02 15:04:05"),
					}
					dataBytes, err := json.Marshal(data)
					if err != nil {
						log.Error("SNMPResource json Marshal error: ", err)
					} else {
						if _, err0 := sr.RuleEngine.Work(sr.Details(), string(dataBytes)); err0 != nil {
							log.Error("SNMPResource PushQueue error: ", err0)
						}

					}
				default:
					{
					}
				}

			}
		}(context.Background(), i, s)
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
	var r []bool
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
func (*SNMPResource) Configs() typex.XConfig {
	config, err := core.RenderConfig("SNMP_SERVER", "", SNMPConfig{})
	if err != nil {
		log.Error(err)
		return typex.XConfig{}
	} else {
		return config
	}
}

//
// 拓扑
//
func (*SNMPResource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}
