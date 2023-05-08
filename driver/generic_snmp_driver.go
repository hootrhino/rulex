package driver

/*
*
* SNMP 数据读取:SNMP 是专门设计用于在IP 网络管理网络节点（服务器、工作站、路由器、交换机及HUBS等）的一种标准协议，
* 它是一种应用层协议。 SNMP 使网络管理员能够管理网络效能，发现并解决网络问题以及规划网络增长。
* 通过SNMP 接收随机消息（及事件报告）网络管理系统获知网络出现问题。
* 更多文档可参考这里: https://info.support.huawei.com/info-finder/encyclopedia/zh/SNMP.html
*
 */
import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gosnmp/gosnmp"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

/*
* Notice:
*  RULEX对于SNMP当前仅仅支持获取硬件信息、内存总量、用户名。
*
 */
type snmpDriver struct {
	state      typex.DriverState
	client     *gosnmp.GoSNMP
	RuleEngine typex.RuleX
	device     *typex.Device
}

func NewSnmpDriver(
	d *typex.Device,
	e typex.RuleX,
	client *gosnmp.GoSNMP,
) typex.XExternalDriver {
	sd := new(snmpDriver)
	sd.client = client
	sd.RuleEngine = e
	sd.device = d
	sd.state = typex.DRIVER_DOWN
	return sd
}
func (sd *snmpDriver) Test() error {
	return nil
}

func (sd *snmpDriver) Init(_ map[string]string) error {
	sd.state = typex.DRIVER_UP
	return nil
}

func (sd *snmpDriver) Work() error {
	return nil
}

func (sd *snmpDriver) State() typex.DriverState {
	return sd.state
}

type _snmp_data struct {
	PCHost        string
	PCDescription string
	PCUserName    string
	PCHardIFaces  []string
	PCTotalMemory int
}

func (sd *snmpDriver) Read(cmd []byte, data []byte) (int, error) {
	bites, err := json.Marshal(_snmp_data{
		PCHost:        sd.client.Target,
		PCDescription: sd.systemDescription(),
		PCUserName:    sd.pCUserName(),
		PCTotalMemory: sd.totalMemory(),
		PCHardIFaces:  sd.hardwareNetInterfaceMac(),
	})
	copy(data, bites)
	return len(bites), err
}

func (sd *snmpDriver) Write(cmd []byte, _ []byte) (int, error) {
	return 0, nil
}

func (sd *snmpDriver) DriverDetail() typex.DriverDetail {
	return typex.DriverDetail{
		Name:        "SNMP-DRIVER",
		Type:        "SNMP",
		Description: "通用SNMP协议客户端",
	}
}

func (sd *snmpDriver) Stop() error {
	return sd.client.Conn.Close()
}

// ----------------------------------------------------------------
func (sd *snmpDriver) connect() error {
	err1 := sd.client.Connect()
	if err1 != nil {
		glogger.GLogger.Error("Connect() err: %v", err1)
		sd.state = typex.DRIVER_DOWN
		return err1
	}
	return nil
}

/*
*
* 获取PC的描述信息
*
 */
func (sd *snmpDriver) systemDescription() string {

	s := ""

	if err1 := sd.connect(); err1 != nil {
		return s
	}

	err := sd.client.Walk(".1.3.6.1.2.1.1.1.0", func(variable gosnmp.SnmpPDU) error {
		if variable.Type == gosnmp.OctetString {
			s = string(variable.Value.([]byte))
		}
		return nil
	})
	if err != nil {
		glogger.GLogger.Error(err)
		sd.state = typex.DRIVER_DOWN
	}
	return s
}
func (sd *snmpDriver) pCUserName() string {
	s := ""
	if err1 := sd.connect(); err1 != nil {
		return s
	}

	err := sd.client.Walk(".1.3.6.1.2.1.1.5.0", func(variable gosnmp.SnmpPDU) error {
		if variable.Type == gosnmp.OctetString {
			s = string(variable.Value.([]byte))
		}
		return nil
	})
	if err != nil {
		glogger.GLogger.Error(err)
		sd.state = typex.DRIVER_DOWN
	}
	return s
}

/*
*
* 获取PC的内存
*
 */
func (sd *snmpDriver) totalMemory() int {
	v := 0
	if err1 := sd.connect(); err1 != nil {
		return v
	}
	err := sd.client.Walk(".1.3.6.1.2.1.25.2.2.0", func(variable gosnmp.SnmpPDU) error {
		if variable.Type == gosnmp.Integer {
			v = int(variable.Value.(int))
		}
		return nil
	})
	if err != nil {
		glogger.GLogger.Error(err)
		sd.state = typex.DRIVER_DOWN
	}
	return v

}

/*
*
* 获取硬件网卡MAC地址
*
 */
func (sd *snmpDriver) hardwareNetInterfaceMac() []string {
	if err1 := sd.connect(); err1 != nil {
		return []string{}
	}
	oid := ".1.3.6.1.2.1.2.2.1.6"
	result := []string{}
	macMaps := map[string]string{}
	err := sd.client.Walk(oid, func(variable gosnmp.SnmpPDU) error {
		if variable.Type == gosnmp.OctetString {
			macHexs := variable.Value.([]uint8)

			if len(macHexs) > 0 {
				if macHexs[0] > 0 {
					hexs := []string{}
					for _, macHex := range macHexs {
						hexs = append(hexs, fmt.Sprintf("%X", macHex))
					}
					macMaps[strings.Join(hexs, ":")] = strings.Join(hexs, ":")
				}

			}
		}
		return nil
	})
	for k := range macMaps {
		result = append(result, k)
	}
	if err != nil {
		glogger.GLogger.Error(err)
		sd.state = typex.DRIVER_DOWN
	}
	return result
}
