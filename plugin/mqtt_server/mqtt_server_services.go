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
			Remote:       v.Remote,
			Username:     string(v.Username),
			CleanSession: v.CleanSession,
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
	if _, ok := s.clients[clientid]; ok {
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
	return typex.ServiceResult{}
}
