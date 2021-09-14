// Atomic 平台云端服务
//
//
package cloud

import (
	"encoding/json"
	"rulex/core"
)

type httpResult struct {
	Code int
	Msg  string
	Data interface{}
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
