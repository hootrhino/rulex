package httpserver

import (
	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/service"
	"github.com/hootrhino/rulex/utils"
)

type VisualVo struct {
	Gid     string `json:"gid"`                         // 分组ID
	UUID    string `json:"uuid"`                        // 名称
	Name    string `json:"name" validate:"required"`    // 名称
	Type    string `json:"type"`                        // 类型
	Content string `json:"content" validate:"required"` // 大屏的内容
}

/*
*
* 新建大屏
*
 */

func CreateVisual(c *gin.Context, hh *HttpApiServer) {
	vvo := VisualVo{}
	if err := c.ShouldBindJSON(&vvo); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	_, err0 := service.GetGenericGroupWithUUID(vvo.Gid)
	if err0 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err0))
		return
	}
	MVisual := model.MVisual{
		UUID:    utils.VisualUuid(),
		Name:    vvo.Name,
		Type:    vvo.Type,
		Content: vvo.Content,
	}
	if err := service.InsertVisual(MVisual); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 新建大屏的时候必须给一个分组
	if err := service.BindResource(vvo.Gid, MVisual.UUID); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 更新大屏
*
 */
func UpdateVisual(c *gin.Context, hh *HttpApiServer) {
	vvo := VisualVo{}
	if err := c.ShouldBindJSON(&vvo); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	MVisual := model.MVisual{
		UUID:    vvo.UUID,
		Name:    vvo.Name,
		Type:    vvo.Type,
		Content: vvo.Content,
	}
	if err := service.UpdateVisual(MVisual); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 删除大屏
*
 */
func DeleteVisual(c *gin.Context, hh *HttpApiServer) {
	uuid, _ := c.GetQuery("uuid")
	if err := service.DeleteVisual(uuid); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 大屏列表
*
 */
func ListVisual(c *gin.Context, hh *HttpApiServer) {
	visuals := []VisualVo{}
	for _, vv := range service.AllVisual() {
		visuals = append(visuals, VisualVo{
			UUID:    vv.UUID,
			Name:    vv.Name,
			Type:    vv.Type,
			Content: vv.Content,
		})
	}
	c.JSON(common.HTTP_OK, common.OkWithData(visuals))

}

/*
*
* 大屏详情
*
 */
func VisualDetail(c *gin.Context, hh *HttpApiServer) {
	uuid, _ := c.GetQuery("uuid")
	mVisual, err := service.GetVisualWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
	}
	vo := VisualVo{
		UUID:    mVisual.UUID,
		Name:    mVisual.Name,
		Type:    mVisual.Type,
		Content: mVisual.Content,
	}
	c.JSON(common.HTTP_OK, common.OkWithData(vo))
}

/*
*
* 生成随机数
*
 */
func GenComponentUUID(c *gin.Context, hh *HttpApiServer) {
	c.JSON(common.HTTP_OK, common.OkWithData(utils.MakeLongUUID("WEIGHT")))
}
