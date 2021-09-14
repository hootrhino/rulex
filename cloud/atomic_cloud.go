// Atomic 平台云端服务
//
//
package cloud

import (
	"encoding/json"
	"rulex/core"
)

type httpResult struct {
	Name        string      `json:"name"`
	Doc         string      `json:"doc"`
	Description interface{} `json:"description"`
}

//
//
//
func ListService(pageIndex int, pageSize int) []httpResult {
	r := GetCloud(
		core.GlobalConfig.Token,
		core.GlobalConfig.Secret,
		core.GlobalConfig.Path+"/data.json",
	)
	results := []httpResult{}
	json.Unmarshal([]byte(r), &results)
	return results
}

//
//
//
func CallService(id string, args []ServiceArg) CallResult {
	return CallResult{}
}
