package rulexlib

// 以腾讯云为目标
// 参考文档：https://cloud.tencent.com/document/product/1081/34916#.E8.AE.BE.E5.A4.87.E8.A1.8C.E4.B8.BA.E8.B0.83.E7.94.A8
import (
	"encoding/json"

	"github.com/hootrhino/rulex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

/*
*
* PropertyReply
* {
*     "method":"property_reply",
*     "id": "20a4ccfd",
*     "code":0,
*     "status":"*****"
* }
*
 */
type _reply struct {
	Method string      `json:"method"`
	Code   int         `json:"code"`
	Id     string      `json:"id"`
	Status string      `json:"status"`
	Out    interface{} `json:"out,omitempty"`
}

func returnR(inend *typex.InEnd, bytes []byte, l *lua.LState) int {
	_, err := inend.Source.DownStream(bytes)
	if err != nil {
		l.Push(lua.LString(err.Error()))
		return 1
	} else {
		l.Push(lua.LNil)
		return 1
	}
}

/*
*
* 属性下发到设备，回复成功
*
 */
func PropertyReplySuccess(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		uuid := l.ToString(2)      // Iothub 资源ID
		requestId := l.ToString(3) // 请求ID
		if inend := rx.GetInEnd(uuid); inend != nil {
			bytes, _ := json.Marshal(_reply{
				Method: "property_reply",
				Code:   0,
				Id:     requestId,
				Status: "Success",
			})
			return returnR(inend, bytes, l)
		} else {
			l.Push(lua.LString("IotHUB resource not exists"))
			return 1
		}
	}
}

/*
*
* 属性下发到设备，回复失败
*
 */
func PropertyReplyFailed(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		uuid := l.ToString(2)      // Iothub 资源ID
		requestId := l.ToString(3) // 请求ID
		if inend := rx.GetInEnd(uuid); inend != nil {
			bytes, _ := json.Marshal(_reply{
				Method: "property_reply",
				Code:   500,
				Id:     requestId,
				Status: "Failed",
			})
			return returnR(inend, bytes, l)
		} else {
			l.Push(lua.LString("IotHUB resource not exists"))
			return 1
		}
	}

}

/*
* 设备行为调用成功
* {
*     "method": "action_reply",
*     "id": "20a4ccfd",
*     "code": 0,
*     "status": "message",
*     "out": {属性}
* }
*
 */
func ActionReplySuccess(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		uuid := l.ToString(2)      // Iothub 资源ID
		requestId := l.ToString(3) // 请求ID
		outParam := l.ToString(4)  // 请求ID

		if inend := rx.GetInEnd(uuid); inend != nil {
			bytes, _ := json.Marshal(_reply{
				Method: "action_reply",
				Code:   0,
				Id:     requestId,
				Status: "Success",
				Out:    outParam,
			})
			return returnR(inend, bytes, l)
		} else {
			l.Push(lua.LString("IotHUB resource not exists"))
			return 1
		}
	}
}

/*
* 设备行为调用失败
* {
*     "method": "action_reply",
*     "id": "20a4ccfd",
*     "code": 0,
*     "status": "message",
*     "out": {属性}
* }
*
 */
func ActionReplyFailed(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		uuid := l.ToString(2)      // Iothub 资源ID
		requestId := l.ToString(3) // 请求ID
		if inend := rx.GetInEnd(uuid); inend != nil {
			bytes, _ := json.Marshal(_reply{
				Method: "action_reply",
				Code:   500,
				Id:     requestId,
				Status: "Failed",
			})
			return returnR(inend, bytes, l)
		} else {
			l.Push(lua.LString("IotHUB resource not exists"))
			return 1
		}
	}
}
