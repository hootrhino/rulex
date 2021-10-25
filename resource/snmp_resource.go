package resource

import (
	"context"
	"encoding/json"
	"rulex/typex"
	"rulex/utils"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/ngaut/log"
)

type SNMPResource struct {
	typex.XStatus
	snmpClient *gosnmp.GoSNMP
}

// GoSNMP represents GoSNMP library state.
type SNMPConfig struct {
	Target     string             `json:"target" validate:"required"`
	Port       uint16             `json:"port" validate:"required"`
	Transport  string             `json:"transport" validate:"required"`
	Community  string             `json:"community" validate:"required"`
	Version    uint8              `json:"version" validate:"required"`
	DataModels []typex.XDataModel `json:"dataModels" validate:"required"`
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
	if err := s.snmpClient.Connect(); err != nil {
		log.Errorf("SnmpClient Connect err: %v", err)
		return false
	} else {
		return true
	}

}

func (s *SNMPResource) Register(inEndId string) error {
	s.PointId = inEndId
	return nil
}

func (s *SNMPResource) Start() error {
	config := s.RuleEngine.GetInEnd(s.PointId).Config
	configBytes, err0 := json.Marshal(config)
	if err0 != nil {
		return err0
	}
	var mainConfig SNMPConfig
	if err1 := json.Unmarshal(configBytes, &mainConfig); err1 != nil {
		return err1
	}
	if err2 := utils.TransformConfig(configBytes, &mainConfig); err2 != nil {
		return err2
	}
	s.snmpClient = gosnmp.Default
	s.snmpClient.Target = mainConfig.Target
	s.snmpClient.Community = mainConfig.Community
	// s.snmpClient.Version = gosnmp.SnmpVersion(mainConfig.Version)
	ticker := time.NewTicker(5 * time.Second)

	if err := s.snmpClient.Connect(); err != nil {
		log.Errorf("SnmpClient Connect err: %v", err)
		return err
	} else {
		go func(ctx context.Context, snmpClient *gosnmp.GoSNMP) {
			defer ticker.Stop()
			for {
				<-ticker.C
				result, err2 := s.snmpClient.Get([]string{".1.3.6.1.2.1.1.1.0"})
				if err2 != nil {
					log.Error(err2)
				}
				for i, variable := range result.Variables {
					log.Infof("%d: oid: %s ", i, variable.Name)

					switch variable.Type {
					case gosnmp.OctetString:
						log.Infof("string: %s\n", string(variable.Value.([]byte)))
					default:
						log.Infof("number: %d\n", gosnmp.ToBigInt(variable.Value))
					}
				}
			}

		}(context.Background(), s.snmpClient)
		log.Info("SNMPResource start successfully!")
		return nil
	}
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
	if err := s.snmpClient.Connect(); err != nil {
		log.Errorf("SnmpClient Connect err: %v", err)
		return typex.DOWN
	} else {
		return typex.UP
	}
}

func (s *SNMPResource) OnStreamApproached(data string) error {
	return nil
}

func (s *SNMPResource) Stop() {

}
