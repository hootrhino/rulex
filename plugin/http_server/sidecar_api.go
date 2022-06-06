package httpserver

import (
	"rulex/typex"
	"rulex/utils"

	"github.com/gin-gonic/gin"
)

/*
*
* Goods
*
 */
func Goodss(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {

	uuid, _ := c.GetQuery("uuid")
	if uuid == "" {
		c.JSON(200, Ok())
	} else {
		e.GetGoods(uuid)
	}
}

/*
*
* CreateGood
*
 */
func CreateGood(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	type Form struct {
		Addr        string   `json:"addr" binding:"required"` // TCP or Unix Socket
		Description string   `json:"description"`             // Description text
		Args        []string `json:"args"`                    // Additional Args
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(200, Error400(err))
		return
	}
	mGoods := MGoods{
		UUID:        utils.GoodsUuid(),
		Addr:        form.Addr,
		Description: form.Description,
		Args:        form.Args,
	}

	if err := hh.InsertGoods(&mGoods); err != nil {
		c.JSON(200, Error400(err))
		return
	}
	goods := typex.Goods{
		UUID:        mGoods.UUID,
		Addr:        mGoods.Addr,
		Description: mGoods.Description,
		Args:        mGoods.Args,
	}
	if err := e.LoadGoods(goods); err != nil {
		c.JSON(200, Error400(err))
		return
	}
}
