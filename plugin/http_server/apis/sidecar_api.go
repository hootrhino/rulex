package apis

import (
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/service"
	"github.com/hootrhino/rulex/trailer"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	"github.com/gin-gonic/gin"
)

/*
*
* Goods
*
 */
func Goods(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if uuid == "" {
		data := []*typex.Goods{}
		trailer.AllGoods().Range(func(key, value interface{}) bool {
			v := value.(*typex.Goods)
			data = append(data, v)
			return true
		})
		c.JSON(common.HTTP_OK, common.OkWithData(data))
	} else {
		c.JSON(common.HTTP_OK, common.OkWithData(trailer.Get(uuid)))
	}
}

/*
*
* 删除外挂
*
 */
func DeleteGoods(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	goods, err := service.GetGoodsWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
	} else {
		// 数据库和内存都要删除
		service.DeleteGoods(goods.UUID)
		trailer.Remove(goods.UUID)
		c.JSON(common.HTTP_OK, common.Ok())
	}
}

/*
*
* CreateGood
*
 */
func CreateGoods(c *gin.Context, ruleEngine typex.RuleX) {
	type Form struct {
		Addr        string   `json:"addr" binding:"required"` // TCP or Unix Socket
		Description string   `json:"description"`             // Description text
		Args        []string `json:"args"`                    // Additional Args
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	mGoods := model.MGoods{
		UUID:        utils.GoodsUuid(),
		Addr:        form.Addr,
		Description: form.Description,
		Args:        form.Args,
	}

	if err := service.InsertGoods(&mGoods); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	goods := typex.Goods{
		UUID:        mGoods.UUID,
		Addr:        mGoods.Addr,
		Description: mGoods.Description,
		Args:        mGoods.Args,
	}
	if err := trailer.Fork(goods); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 更新操作
*
 */
func UpdateGoods(c *gin.Context, ruleEngine typex.RuleX) {
	c.JSON(common.HTTP_OK, common.Error("暂不支持更新"))
}
