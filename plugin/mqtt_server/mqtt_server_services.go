package mqttserver

import (
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 获取当前连接进来的MQTT客户端
*
 */
type Client struct {
	ID           string `json:"id"`
	Remote       string `json:"remote"`
	Listener     string `json:"listener"`
	Username     string `json:"username"`
	CleanSession bool   `json:"cleanSession"`
}

func (s *MqttServer) ListClients(offset, count int) []Client {
	result := []Client{}
	for _, v := range s.clients {
		c := Client{
			ID:           v.ID,
			Remote:       v.Net.Remote,
			Username:     string(v.Properties.Username),
			CleanSession: v.Properties.Clean,
			Listener:     v.Net.Listener,
		}
		result = append(result, c)
	}
	return result
}

/*
*
* 把某个客户端给踢下线
*
 */
func (s *MqttServer) KickOut(clientid string) bool {
	if client, ok := s.clients[clientid]; ok {
		client.Stop(nil)
		return true
	}
	return false

}

/*
*
* 服务调用接口
*
 */
func (s *MqttServer) Service(arg typex.ServiceArg) typex.ServiceResult {
	if arg.Name == "clients" {
		return typex.ServiceResult{Out: s.ListClients(0, 100)}
	}
	if arg.Name == "kickout" {
		switch tt := arg.Args.(type) {
		case []interface{}:
			{
				for _, id := range tt {
					switch idt := id.(type) {
					case string:
						{
							s.KickOut(idt)
						}
					}
				}
			}
		default:
			{
				return typex.ServiceResult{Out: []string{"unsupported args:" + arg.Name}}
			}
		}
		return typex.ServiceResult{Out: []string{"kick out success"}}
	}
	return typex.ServiceResult{Out: []string{"no such service name:" + arg.Name}}
}
