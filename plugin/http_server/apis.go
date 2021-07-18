package httpserver

import (
	"net/http"
	"rulex/core"

	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2/json"
)

// Http Return
type R struct {
	Code int
	Msg  string
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
func cros(c *gin.Context) {
	method := c.Request.Method
	origin := c.Request.Header.Get("Origin")
	if origin != "" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		c.Header("Access-Control-Max-Age", "172800")
		c.Header("Access-Control-Allow-Credentials", "true")
	}

	if method == "OPTIONS" {
		c.JSON(http.StatusOK, "ok!")
	}
}

//
// LoadNewestInEnd
//
func (hh *HttpApiServer) LoadNewestInEnd(uuid string) error {
	mInEnd, _ := hh.GetMInEndWithUUID(uuid)
	config := map[string]interface{}{}
	if err := json.Unmarshal([]byte(mInEnd.Config), &config); err != nil {
		return err
	} else {
		in := core.NewInEnd(mInEnd.Type, mInEnd.Name, mInEnd.Description, &config)
		// Important !!!!!!!!
		in.Id = mInEnd.UUID
		if err := hh.ruleEngine.LoadInEnd(in); err != nil {
			return err
		} else {
			return nil
		}
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
	} else {
		out := core.NewOutEnd(mOutEnd.Type, mOutEnd.Name, mOutEnd.Description, &config)
		// Important !!!!!!!!
		out.Id = mOutEnd.UUID
		if err := hh.ruleEngine.LoadOutEnd(out); err != nil {
			return err
		} else {
			return nil
		}
	}
}
