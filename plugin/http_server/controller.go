package httpserver

import (
	"net/http"
	"rulex/core"
	"rulex/statistics"
	"rulex/typex"
	"rulex/utils"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ngaut/log"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"gopkg.in/square/go-jose.v2/json"
)

//
//
//
type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

//
// Get all plugins
//
func Plugins(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	data := []interface{}{}
	for _, v := range e.AllPlugins() {
		data = append(data, v)
	}
	c.PureJSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Msg:  "Success",
		Data: data,
	})
}

//
// Get system infomation
//
func System(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	percent, _ := cpu.Percent(time.Second, false)
	memInfo, _ := mem.VirtualMemory()
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	c.PureJSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Msg:  "Success",
		Data: gin.H{
			"rulexVersion": e.Version().Version,
			"diskInfo":     diskInfo.UsedPercent,
			"memInfo":      memInfo.UsedPercent,
			"cpuPercent":   percent[0],
			"os":           runtime.GOOS,
			"arch":         runtime.GOARCH,
			"cpus":         runtime.GOMAXPROCS(0)},
	})
}

//
// Get all inends
//
func InEnds(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	data := []interface{}{}
	for _, v := range e.AllInEnd() {
		data = append(data, v)
	}
	c.PureJSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Msg:  "Success",
		Data: data,
	})
}

//
// Get all outends
//
func OutEnds(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	data := []interface{}{}
	for _, v := range e.AllOutEnd() {
		data = append(data, v)
	}
	c.PureJSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Msg:  "Success",
		Data: data,
	})
}

//
// Get all rules
//
func Rules(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	data := []interface{}{}
	for _, v := range e.AllRule() {
		data = append(data, v)
	}
	c.PureJSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Msg:  "Success",
		Data: data,
	})
}

//
// Get statistics data
//
func Statistics(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	c.PureJSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Msg:  "Success",
		Data: statistics.AllStatistics(),
	})
}

//
// All Users
//
func Users(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	users := hh.AllMUser()
	c.PureJSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Msg:  "Success",
		Data: users,
	})
}

//
// Create InEnd
//
func CreateInend(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	type Form struct {
		Type        string                 `json:"type" binding:"required"`
		Name        string                 `json:"name" binding:"required"`
		Description string                 `json:"description"`
		Config      map[string]interface{} `json:"config" binding:"required"`
	}
	form := Form{}
	err0 := c.ShouldBindJSON(&form)
	if err0 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err0.Error()})
		return
	} else {
		configJson, err1 := json.Marshal(form.Config)
		if err1 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err1.Error()})
			return
		} else {
			uuid := utils.MakeUUID("INEND")
			hh.InsertMInEnd(&MInEnd{
				UUID:        uuid,
				Type:        form.Type,
				Name:        form.Name,
				Description: form.Description,
				Config:      string(configJson),
			})
			if err := hh.LoadNewestInEnd(uuid); err != nil {
				log.Error(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
				return
			} else {
				c.PureJSON(http.StatusOK, gin.H{"msg": "create success"})
				return
			}
		}
	}
}

//
// Create OutEnd
//
func CreateOutEnd(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	type Form struct {
		Type        string                 `json:"type" binding:"required"`
		Name        string                 `json:"name" binding:"required"`
		Description string                 `json:"description"`
		Config      map[string]interface{} `json:"config" binding:"required"`
	}
	form := Form{}
	err0 := c.ShouldBindJSON(&form)
	if err0 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err0.Error()})
		return
	} else {
		configJson, err1 := json.Marshal(form.Config)
		if err1 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err1.Error()})
			return
		} else {
			uuid := utils.MakeUUID("OUTEND")
			hh.InsertMOutEnd(&MOutEnd{
				UUID:        uuid,
				Type:        form.Type,
				Name:        form.Name,
				Description: form.Description,
				Config:      string(configJson),
			})
			err := hh.LoadNewestOutEnd(uuid)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
				return
			} else {
				c.JSON(http.StatusOK, gin.H{"msg": "create success"})
				return
			}
		}
	}
}

//
// Create rule
//
func CreateRule(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	type Form struct {
		From        []string `json:"from" binding:"required"`
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description"`
		Actions     string   `json:"actions"`
		Success     string   `json:"success"`
		Failed      string   `json:"failed"`
	}
	form := Form{}
	err0 := c.ShouldBindJSON(&form)
	if err0 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err0.Error()})
		return
	} else {

		if len(form.From) > 0 {
			for _, id := range form.From {
				if id != "" {
					if e.GetInEnd(id) == nil {
						c.JSON(http.StatusBadRequest, gin.H{"msg": "inend not exists:" + id})
						return
					}
				} else {
					c.JSON(http.StatusOK, gin.H{"msg": "invalid 'from'"})
					return
				}
			}
			tmpRule := typex.NewRule(nil,
				form.Name,
				form.Description,
				nil,
				form.Success,
				form.Actions,
				form.Failed)

			if err1 := core.VerifyCallback(tmpRule); err1 != nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": err1.Error()})
				return
			} else {
				mRule := &MRule{
					Name:        form.Name,
					Description: form.Description,
					From:        form.From,
					Success:     form.Success,
					Failed:      form.Failed,
					Actions:     form.Actions,
				}
				if err := hh.InsertMRule(mRule); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
					return
				}
				rule := typex.NewRule(hh.ruleEngine,
					mRule.Name,
					mRule.Description,
					mRule.From,
					mRule.Success,
					mRule.Actions,
					mRule.Failed)
				if err := e.LoadRule(rule); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
				} else {
					c.JSON(http.StatusOK, gin.H{"msg": "create success"})
				}
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "from can't empty"})
			return
		}
	}
}

//
// Delete inend by UUID
//
func DeleteInend(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, exists := c.GetQuery("uuid")
	if exists {
		// Important !!!!!
		e.GetInEnd(uuid).Resource.Stop()
		e.RemoveInEnd(uuid)  //1
		hh.DeleteMRule(uuid) //2
		//
		c.JSON(http.StatusOK, gin.H{"msg": "remove success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "rule not exists"})
	}
}

//
// Delete outend by UUID
//
func DeleteOutend(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, exists := c.GetQuery("uuid")
	if exists {
		// Important !!!!!
		e.GetOutEnd(uuid).Target.Stop()
		e.RemoveRule(uuid)     //1
		hh.DeleteMOutEnd(uuid) //2
		//
		c.JSON(http.StatusOK, gin.H{"msg": "remove success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "rule not exists"})
	}
}

//
// Delete rule by UUID
//
func DeleteRule(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, exists := c.GetQuery("uuid")
	if exists {
		// Important !!!!! e.RemoveRule(uuid)
		e.RemoveRule(uuid)   //1
		hh.DeleteMRule(uuid) //2
		//
		c.JSON(http.StatusOK, gin.H{"msg": "remove success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "rule not exists"})
	}
}

//
// CreateUser
//
func CreateUser(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	type Form struct {
		Role        string `json:"role" binding:"required"`
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required"`
		Description string `json:"description"`
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.PureJSON(http.StatusOK, Result{
			Code: http.StatusBadGateway,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}

	if user, err := hh.GetMUser(form.Username, form.Password); err != nil {
		c.PureJSON(http.StatusOK, Result{
			Code: http.StatusBadGateway,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	} else {
		if user.ID > 0 {
			c.PureJSON(http.StatusOK, Result{
				Code: http.StatusBadGateway,
				Msg:  "用户已存在:" + user.Username,
				Data: nil,
			})
			return
		} else {
			hh.InsertMUser(&MUser{
				Role:        form.Role,
				Username:    form.Username,
				Password:    form.Password,
				Description: form.Description,
			})
			c.PureJSON(http.StatusOK, Result{
				Code: http.StatusOK,
				Msg:  "用户创建成功",
				Data: form.Username,
			})
			return
		}
	}
}

//
// Auth
//
func Auth(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	type Form struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	form := Form{}
	err0 := c.ShouldBindJSON(&form)
	if err0 != nil {
		c.PureJSON(http.StatusOK, Result{
			Code: http.StatusBadGateway,
			Msg:  err0.Error(),
			Data: nil,
		})
	} else {
		user, err1 := hh.GetMUser(form.Username, form.Password)
		if err1 != nil {
			c.PureJSON(http.StatusOK, Result{
				Code: http.StatusBadGateway,
				Msg:  err1.Error(),
				Data: nil,
			})
		} else {
			c.PureJSON(http.StatusOK, Result{
				Code: http.StatusOK,
				Msg:  "Auth Success",
				Data: user.Username,
			})
		}
	}
}
