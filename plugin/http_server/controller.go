package httpserver

import (
	"net/http"
	"rulex/core"
	"rulex/statistics"
	"rulex/utils"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ngaut/log"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"gopkg.in/square/go-jose.v2/json"
)

//
// Render dashboard index
//
func Index(c *gin.Context, hh *HttpApiServer, e core.RuleX) {
	c.HTML(http.StatusOK, "dashboard.html", gin.H{})
}

//
// Get all plugins
//
func Plugins(c *gin.Context, hh *HttpApiServer, e core.RuleX) {
	cros(c)
	c.PureJSON(http.StatusOK, gin.H{
		"plugins": e.AllPlugins(),
	})
}

//
// Get system infomation
//
func System(c *gin.Context, hh *HttpApiServer, e core.RuleX) {
	cros(c)
	//
	percent, _ := cpu.Percent(time.Second, false)
	memInfo, _ := mem.VirtualMemory()
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	c.JSON(http.StatusOK, gin.H{
		"diskInfo":   diskInfo.UsedPercent,
		"memInfo":    memInfo.UsedPercent,
		"cpuPercent": percent[0],
		"os":         runtime.GOOS,
		"arch":       runtime.GOARCH,
		"cpus":       runtime.GOMAXPROCS(0)})
}

//
// Get all inends
//
func InEnds(c *gin.Context, hh *HttpApiServer, e core.RuleX) {
	cros(c)
	c.JSON(http.StatusOK, gin.H{"inends": e.AllInEnd()})
}

//
// Get all outends
//
func OutEnds(c *gin.Context, hh *HttpApiServer, e core.RuleX) {
	cros(c)
	c.JSON(http.StatusOK, gin.H{"outends": e.AllOutEnd()})
}

//
// Get all rules
//
func Rules(c *gin.Context, hh *HttpApiServer, e core.RuleX) {
	cros(c)
	c.JSON(http.StatusOK, gin.H{"rules": e.AllRule()})
}

//
// Get statistics data
//
func Statistics(c *gin.Context, hh *HttpApiServer, e core.RuleX) {
	cros(c)
	c.JSON(http.StatusOK, gin.H{"statistics": statistics.AllStatistics()})
}

//
//
//
func Users(c *gin.Context, hh *HttpApiServer, e core.RuleX) {
	cros(c)
	type Form struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	form := Form{}
	err0 := c.ShouldBindJSON(&form)
	if err0 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err0.Error()})
	} else {
		user, err1 := hh.GetMUser(form.Username, form.Password)
		if err1 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err1.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"msg": "success", "data": user.ID})
		}
	}
}

//
// Create InEnd
//
func CreateInend(c *gin.Context, hh *HttpApiServer, e core.RuleX) {
	cros(c)
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
	} else {
		configJson, err1 := json.Marshal(form.Config)
		if err1 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err1.Error()})
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
			} else {
				c.JSON(http.StatusOK, gin.H{"msg": "create success"})
			}
		}
	}
}

//
// Create OutEnd
//
func CreateOutEnd(c *gin.Context, hh *HttpApiServer, e core.RuleX) {
	cros(c)
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
	} else {
		configJson, err1 := json.Marshal(form.Config)
		if err1 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err1.Error()})
		} else {
			// TODO : Load newest OutEnd
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
			} else {
				c.JSON(http.StatusOK, gin.H{"msg": "create success"})
			}
		}
	}
}

//
// Create rule
//
func CreateRule(c *gin.Context, hh *HttpApiServer, e core.RuleX) {
	cros(c)
	type Form struct {
		From        string `json:"from" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Actions     string `json:"actions"`
		Success     string `json:"success"`
		Failed      string `json:"failed"`
	}
	form := Form{}
	err0 := c.ShouldBindJSON(&form)
	if err0 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err0.Error()})
	} else {
		rule := core.NewRule(nil,
			form.Name,
			form.Description,
			nil,
			form.Success,
			form.Actions,
			form.Failed)
		if len(strings.Split(form.From, ",")) > 0 {
			for _, id := range strings.Split(form.From, ",") {
				// must be: 111,222,333... style
				if id != "" {
					if e.GetInEnd(id) == nil {
						c.JSON(http.StatusBadRequest, gin.H{"msg": "inend not exists:" + id})
						return
					}
				} else {
					c.JSON(http.StatusOK, gin.H{"msg": "invalid 'from' string format:" + form.From})
					return
				}
			}
			if err1 := core.VerifyCallback(rule); err1 != nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": err1.Error()})
			} else {
				mRule := &MRule{
					Name:        form.Name,
					Description: form.Description,
					From:        form.From,
					Success:     form.Success,
					Failed:      form.Failed,
					Actions:     form.Actions,
				}
				hh.InsertMRule(mRule)
				rule := core.NewRule(hh.ruleEngine,
					mRule.Name,
					mRule.Description,
					strings.Split(mRule.From, ","),
					mRule.Success,
					mRule.Actions,
					mRule.Failed)
				if err := e.LoadRule(rule); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
				} else {
					c.JSON(http.StatusOK, gin.H{"msg": "create success"})
				}
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "from can't empty"})
		}
	}
}

//
// Delete inend by UUID
//
func DeleteInend(c *gin.Context, hh *HttpApiServer, e core.RuleX) {
	cros(c)
	uuid, exists := c.GetQuery("uuid")
	if exists {
		// Important !!!!!
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
func DeleteOutend(c *gin.Context, hh *HttpApiServer, e core.RuleX) {
	cros(c)
	uuid, exists := c.GetQuery("uuid")
	if exists {
		// Important !!!!!
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
func DeleteRule(c *gin.Context, hh *HttpApiServer, e core.RuleX) {
	cros(c)
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
// Auth
//
func Auth(c *gin.Context, hh *HttpApiServer, e core.RuleX) {
	cros(c)
	type Form struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	form := Form{}
	err0 := c.ShouldBindJSON(&form)
	if err0 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err0.Error()})
	} else {
		user, err1 := hh.GetMUser(form.Username, form.Password)
		if err1 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err1.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"msg": "success", "data": user.ID})
		}
	}
}
