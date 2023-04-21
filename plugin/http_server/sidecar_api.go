package httpserver

import (
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"

	"github.com/gin-gonic/gin"
)

/*
*
* Goods
*
 */
func Goods(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if uuid == "" {
		data := []*typex.Goods{}
		e.AllGoods().Range(func(key, value interface{}) bool {
			v := value.(*typex.Goods)
			data = append(data, v)
			return true
		})
		c.JSON(200, OkWithData(data))
	} else {
		c.JSON(200, OkWithData(e.GetGoods(uuid)))
	}
}

/*
*
* 删除外挂
*
 */
func DeleteGoods(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	goods, err := hh.GetGoodsWithUUID(uuid)
	if err != nil {
		c.JSON(200, Error400(err))
	} else {
		// 数据库和内存都要删除
		hh.DeleteGoods(goods.UUID)
		e.RemoveGoods(goods.UUID)
		c.JSON(200, Ok())
	}
}

/*
*
* CreateGood
*
 */
func CreateGoods(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
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
	c.JSON(200, Ok())
}

/*
*
* 更新操作
*
 */
func UpdateGoods(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	c.JSON(200, Error("暂不支持更新"))
}
