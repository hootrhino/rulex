package httpserver

import (
	"errors"
	"net/http"

	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"

	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2/json"
)

const SUCCESS string = "Success"

// Http Return
type R struct {
	Code int    `json:"code" binding:"required"`
	Msg  string `json:"msg" binding:"required"`
}

//
//
//
type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Ok() R {
	return R{200, "操作成功"}
}
func OkWithEmpty() Result {
	return Result{200, "操作成功", []interface{}{}}
}
func OkWithData(data interface{}) Result {
	return Result{200, "操作成功", data}
}
func Error(s string) R {
	return R{4001, s}
}
func Error400(e error) R {
	return R{4001, e.Error()}
}
func Error500(e error) R {
	return R{5001, e.Error()}
}
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token != "" {
			// TODO add jwt Authorize support
			c.Next()
		} else {
			// c.Abort()
			// c.JSON(http.StatusUnauthorized, gin.H{"msg": "No authority operate"})
			c.Next()
			return
		}
	}
}

//
//
//
func Cros() gin.HandlerFunc {
	return func(c *gin.Context) {
		cros(c)
	}
}

//
func cros(c *gin.Context) {
	c.Header("Cache-Control", "private, max-age=10")
	method := c.Request.Method
	origin := c.Request.Header.Get("Origin")

	c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
	c.Header("Access-Control-Max-Age", "172800")
	c.Header("Access-Control-Allow-Credentials", "true")

	if method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Request.Header.Del("Origin")
	c.Next()
}

//
// Add api route
//
func (h *HttpApiServer) addRoute(f func(*gin.Context, *HttpApiServer, typex.RuleX)) func(*gin.Context) {

	return func(c *gin.Context) {
		f(c, h, h.ruleEngine)
	}
}

//
// LoadNewestInEnd
//
func (hh *HttpApiServer) LoadNewestInEnd(uuid string) error {
	mInEnd, _ := hh.GetMInEndWithUUID(uuid)
	if mInEnd == nil {
		return errors.New("Inend not exists:" + uuid)
	}
	config := map[string]interface{}{}
	if err1 := json.Unmarshal([]byte(mInEnd.Config), &config); err1 != nil {
		glogger.GLogger.Error(err1)
		return err1
	}
	// :mInEnd: {k1 :{k1:v1}, k2 :{k2:v2}} --> InEnd: [{k1:v1}, {k2:v2}]
	var dataModelsMap map[string]typex.XDataModel
	if err1 := json.Unmarshal([]byte(mInEnd.XDataModels), &dataModelsMap); err1 != nil {
		glogger.GLogger.Error(err1)
		return err1
	}
	in := typex.NewInEnd(mInEnd.Type, mInEnd.Name, mInEnd.Description, config)
	// Important !!!!!!!! in.Id = mInEnd.UUID
	in.UUID = mInEnd.UUID
	in.DataModelsMap = dataModelsMap
	if err2 := hh.ruleEngine.LoadInEnd(in); err2 != nil {
		glogger.GLogger.Error(err2)
		return err2
	} else {
		return nil
	}

}

//
// LoadNewestOutEnd
//
func (hh *HttpApiServer) LoadNewestOutEnd(uuid string) error {
	mOutEnd, _ := hh.GetMOutEndWithUUID(uuid)
	config := map[string]interface{}{}
	if err := json.Unmarshal([]byte(mOutEnd.Config), &config); err != nil {
		return err
	}
	out := typex.NewOutEnd(typex.TargetType(mOutEnd.Type), mOutEnd.Name, mOutEnd.Description, config)
	// Important !!!!!!!!
	out.UUID = mOutEnd.UUID
	if err := hh.ruleEngine.LoadOutEnd(out); err != nil {
		return err
	} else {
		return nil
	}

}

//
// LoadNewestDevice
//
func (hh *HttpApiServer) LoadNewestDevice(uuid string) error {
	mDevice, _ := hh.GetDeviceWithUUID(uuid)
	config := map[string]interface{}{}
	if err := json.Unmarshal([]byte(mDevice.Config), &config); err != nil {
		return err
	}
	dev := typex.NewDevice(typex.DeviceType(mDevice.Type), mDevice.Name, mDevice.Description, mDevice.ActionScript, config)
	// Important !!!!!!!!
	dev.UUID = mDevice.UUID
	if err := hh.ruleEngine.LoadDevice(dev); err != nil {
		return err
	} else {
		return nil
	}

}
