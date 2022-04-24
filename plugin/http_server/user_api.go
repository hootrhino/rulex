package httpserver

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"rulex/core"
	"rulex/typex"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	SECRETKEY = "you-can-not-get-this-secret"
)

//
// All Users
//
func Users(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	users := hh.AllMUser()
	c.JSON(http.StatusOK, OkWithData(users))
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
		c.JSON(http.StatusBadRequest, Error400(err))
		return
	}

	if _, err := hh.GetMUser(form.Username, md5Hash(form.Password)); err != nil {
		// 查不到说明没有，可以新建
		hh.InsertMUser(&MUser{
			Role:        form.Role,
			Username:    form.Username,
			Password:    md5Hash(form.Password),
			Description: form.Description,
		})
		c.JSON(http.StatusOK, Ok())
		return
	}
	c.JSON(http.StatusBadRequest, Error("用户名已存在:"+form.Username))
}

/*
*
* Md5 计算
*
 */
func md5Hash(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//
// Login
// TODO: 下个版本实现用户基础管理
//
func Login(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	type _user struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var u _user
	if err := c.BindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, Error400(err))
		return
	}
	if _, err := hh.GetMUser(u.Username, md5Hash(u.Password)); err != nil {
		c.JSON(http.StatusBadRequest, Error400(err))
		return
	}
	if token, err := generateToken(u.Username); err != nil {
		c.JSON(http.StatusBadRequest, Error400(err))
		return
	} else {
		c.JSON(http.StatusOK, OkWithData(token))
	}
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
	token := c.GetHeader("token")
	if claims, err := parseToken(token); err != nil {
		c.JSON(http.StatusBadRequest, Error400(err))
		return
	} else {
		c.JSON(http.StatusOK, OkWithData(map[string]interface{}{
			"token":  token,
			"avatar": "rulex",
			"name":   claims.Username,
		}))
	}

}

type JwtClaims struct {
	Username string
	jwt.StandardClaims
}

/*
*
* 生成Token
*
 */
func generateToken(username string) (string, error) {
	claims := &JwtClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(60*60*24) * time.Second).Unix(),
			Issuer:    username,
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRETKEY))
	return token, err
}

/*
*
* 解析Token
*
 */
func parseToken(tokenString string) (*JwtClaims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("expected token string on headers")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(SECRETKEY), nil
		})
	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
