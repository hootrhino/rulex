package apis

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"unicode/utf8"

	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/service"
	"github.com/hootrhino/rulex/typex"

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

func UserDetail(c *gin.Context, ruleEngine typex.RuleX) {
	Info(c, ruleEngine)
}
func Users(c *gin.Context, ruleEngine typex.RuleX) {
	users := []user{}
	for _, u := range service.AllMUser() {
		users = append(users, user{
			Role:        u.Role,
			Username:    u.Username,
			Description: u.Description,
		})
	}
	c.JSON(common.HTTP_OK, common.OkWithData(users))
}
func isLengthBetween8And16(str string) bool {
	length := utf8.RuneCountInString(str)
	return length >= 8 && length <= 16
}

// CreateUser
func CreateUser(c *gin.Context, ruleEngine typex.RuleX) {
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
	if !isLengthBetween8And16(form.Username) {
		c.JSON(common.HTTP_OK, common.Error("Username Length must Between 8 ~ 16"))
		return
	}
	if !isLengthBetween8And16(form.Password) {
		c.JSON(common.HTTP_OK, common.Error("Password Length must Between 8 ~ 16"))
		return
	}
	if _, err := service.GetMUser(form.Username); err != nil {
		service.InsertMUser(&model.MUser{
			Role:        form.Role,
			Username:    form.Username,
			Password:    md5Hash(form.Password),
			Description: form.Description,
		})
		c.JSON(common.HTTP_OK, common.Ok())
		return
	}
	c.JSON(common.HTTP_OK, common.Error("user already exists:"+form.Username))
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
func Login(c *gin.Context, ruleEngine typex.RuleX) {
	type _user struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var u _user
	if err := c.BindJSON(&u); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if _, err := service.Login(u.Username, md5Hash(u.Password)); err != nil {
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
func Logs(c *gin.Context, ruleEngine typex.RuleX) {
	type Data struct {
		Id      int    `json:"id" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	//TODO 日志暂时不记录
	logs := []Data{}
	c.JSON(common.HTTP_OK, common.OkWithData(logs))
}

func LogOut(c *gin.Context, ruleEngine typex.RuleX) {
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* TODO：用户信息, 当前版本写死 下个版本实现数据库查找
*
 */
func Info(c *gin.Context, ruleEngine typex.RuleX) {
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

/*
*
* 上传头像
*
 */
func UploadSysLogo(c *gin.Context, ruleEngine typex.RuleX) {
	// single file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	fileName := "logo.png"
	dir := "./upload/Logo/"
	if err := os.MkdirAll(filepath.Dir(dir), os.ModePerm); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err := c.SaveUploadedFile(file, dir+fileName); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	err1 := service.UpdateSiteConfig(model.MSiteConfig{
		Logo: dir + fileName,
	})
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(map[string]string{
		"url": "/api/v1/site/logo",
	}))
}

/*
*
* 加载头像
*
 */
func GetSysLogo(c *gin.Context, ruleEngine typex.RuleX) {
	MSiteConfig, err1 := service.GetSiteConfig()
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	var Binary []byte
	// data:image/png;base64,
	if len(MSiteConfig.Logo) < 22 {
		var err1 error
		Binary, err1 = base64.StdEncoding.DecodeString(model.SysDefaultLogo[22:])
		if err1 != nil {
			c.JSON(common.HTTP_OK, common.Error400(err1))
			return
		}
	} else {
		var err2 error
		Binary, err2 = base64.StdEncoding.DecodeString(MSiteConfig.Logo[22:])
		if err2 != nil {
			c.JSON(common.HTTP_OK, common.Error400(err2))
			return
		}
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Header().Set("Content-Type", "image/jpeg")
	c.Writer.Header().Set("Content-Length", strconv.Itoa(len(Binary)))
	c.Writer.Write(Binary)
	c.Writer.Flush()
}

/*
*
* 重置站点
*
 */
func ResetSiteConfig(c *gin.Context, ruleEngine typex.RuleX) {
	err1 := service.UpdateSiteConfig(model.MSiteConfig{
		SiteName: "Rhino EEKit",
		Logo:     "/logo.png",
		AppName:  "Rhino EEKit",
	})
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}
