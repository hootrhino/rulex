package httpserver

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	common "github.com/hootrhino/rulex/plugin/http_server/common"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	SECRETKEY = "you-can-not-get-this-secret"
)

// All Users
type user struct {
	Role        string `json:"role"`
	Username    string `json:"username"`
	Description string `json:"description"`
}

func UserDetail(c *gin.Context, hh *HttpApiServer) {
	Info(c, hh)
}
func Users(c *gin.Context, hh *HttpApiServer) {
	users := []user{}
	for _, u := range hh.AllMUser() {
		users = append(users, user{
			Role:        u.Role,
			Username:    u.Username,
			Description: u.Description,
		})
	}
	c.JSON(common.HTTP_OK, common.OkWithData(users))
}

// CreateUser
func CreateUser(c *gin.Context, hh *HttpApiServer) {
	type Form struct {
		Role        string `json:"role" binding:"required"`
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required"`
		Description string `json:"description"`
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	if _, err := hh.GetMUser(form.Username, md5Hash(form.Password)); err != nil {
		hh.InsertMUser(&MUser{
			Role:        form.Role,
			Username:    form.Username,
			Password:    md5Hash(form.Password),
			Description: form.Description,
		})
		c.JSON(common.HTTP_OK, common.Ok())
		return
	}
	c.JSON(common.HTTP_OK, common.Error("用户名已存在:"+form.Username))
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

// Login
// TODO: 下个版本实现用户基础管理
func Login(c *gin.Context, hh *HttpApiServer) {
	type _user struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var u _user
	if err := c.BindJSON(&u); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if _, err := hh.GetMUser(u.Username, md5Hash(u.Password)); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if token, err := generateToken(u.Username); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	} else {
		c.JSON(common.HTTP_OK, common.OkWithData(token))
	}
}

/*
*
* 日志管理
*
 */
func Logs(c *gin.Context, hh *HttpApiServer) {
	type Data struct {
		Id      int    `json:"id" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	//TODO 日志暂时不记录
	logs := []Data{}
	c.JSON(common.HTTP_OK, common.OkWithData(logs))
}

func LogOut(c *gin.Context, hh *HttpApiServer) {
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* TODO：用户信息, 当前版本写死 下个版本实现数据库查找
*
 */
func Info(c *gin.Context, hh *HttpApiServer) {
	token := c.GetHeader("token")
	if claims, err := parseToken(token); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	} else {
		c.JSON(common.HTTP_OK, common.OkWithData(map[string]interface{}{
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
