package httpserver

import (
	"net/http"
	"rulex/core"
	"rulex/typex"

	"github.com/gin-gonic/gin"
)

//
// All Users
//
func Users(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	users := hh.AllMUser()
	c.JSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Msg:  SUCCESS,
		Data: users,
	})
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
		c.JSON(http.StatusOK, Result{
			Code: http.StatusBadGateway,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}

	if user, err := hh.GetMUser(form.Username, form.Password); err != nil {
		c.JSON(http.StatusOK, Result{
			Code: http.StatusBadGateway,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	} else {
		if user.ID > 0 {
			c.JSON(http.StatusOK, Result{
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
			c.JSON(http.StatusOK, Result{
				Code: http.StatusOK,
				Msg:  "用户创建成功",
				Data: form.Username,
			})
			return
		}
	}
}

//
// Login
// TODO: 下个版本实现用户基础管理
//
func Login(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	c.JSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Msg:  "Auth Success",
		Data: map[string]interface{}{
			"token":  "token",
			"avatar": "rulex",
			"name":   "rulex",
		},
	})
}

/*
*
* 日志管理
*
 */
func Logs(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	type Data struct {
		Id      int    `json:"id" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	logs := []Data{}
	for i, s := range core.GLOBAL_LOGGER.Slot() {
		if s != "" {
			logs = append(logs, Data{i, s})
		}
	}
	c.JSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Msg:  SUCCESS,
		Data: logs,
	})
}

func LogOut(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	c.JSON(http.StatusOK, Ok())
}

/*
*
* TODO：用户信息, 当前版本写死 下个版本实现数据库查找
*
 */
func Info(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	c.JSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Msg:  "Auth Success",
		Data: map[string]interface{}{
			"token":  "token",
			"avatar": "rulex",
			"name":   "rulex",
		},
	})
}
