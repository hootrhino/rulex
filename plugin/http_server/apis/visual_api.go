package apis

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/service"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

type VisualVo struct {
	Gid       string `json:"gid"`                         // 分组ID
	UUID      string `json:"uuid"`                        // 名称
	Name      string `json:"name" validate:"required"`    // 名称
	Type      string `json:"type"`                        // 类型
	Content   string `json:"content" validate:"required"` // 大屏的内容
	Thumbnail string `json:"thumbnail"`
	Status    *bool  `json:"status"`
}

/*
*
* 新建大屏
*
 */

func CreateVisual(c *gin.Context, ruleEngine typex.RuleX) {
	form := VisualVo{Type: "BUILDIN"}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	_, err0 := service.GetGenericGroupWithUUID(form.Gid)
	if err0 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err0))
		return
	}
	MVisual := model.MVisual{
		UUID:    utils.VisualUuid(),
		Name:    form.Name,
		Type:    form.Type,
		Content: form.Content,
		Status:  false,
	}
	if err := service.InsertVisual(MVisual); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 新建大屏的时候必须给一个分组
	if err := service.BindResource(form.Gid, MVisual.UUID); err != nil {
		c.JSON(common.HTTP_OK, common.Error("Group not found"))
		return
	}
	// 返回新建的大屏字段 用来跳转编辑器
	c.JSON(common.HTTP_OK, common.OkWithData(map[string]string{
		"uuid": MVisual.UUID,
	}))

}

/*
*
* 更新大屏
*
 */
func UpdateVisual(c *gin.Context, ruleEngine typex.RuleX) {
	form := VisualVo{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	MVisual := model.MVisual{
		UUID:    form.UUID,
		Name:    form.Name,
		Type:    form.Type,
		Content: form.Content,
	}

	if err := service.UpdateVisual(MVisual); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 取消绑定分组,删除原来旧的分组
	Group := service.GetVisualGroup(MVisual.UUID)
	if err := service.UnBindResource(Group.UUID, MVisual.UUID); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 重新绑定分组
	if err := service.BindResource(form.Gid, MVisual.UUID); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 返回新建的大屏字段 用来跳转编辑器
	c.JSON(common.HTTP_OK, common.OkWithData(map[string]string{
		"uuid": MVisual.UUID,
	}))
}

/*
*
* 发布
*
 */
func PublishVisual(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	MVisual, err := service.GetVisualWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if MVisual.Status {
		c.JSON(common.HTTP_OK, common.Error("Already published:"+MVisual.Name))
		return
	}
	MVisual.Status = true
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
func DeleteVisual(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if err := service.DeleteVisual(uuid); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 大屏列表
*
 */
func ListVisual(c *gin.Context, ruleEngine typex.RuleX) {
	visuals := []VisualVo{}
	for _, vv := range service.AllVisual() {
		Vo := VisualVo{
			UUID:      vv.UUID,
			Name:      vv.Name,
			Type:      vv.Type,
			Content:   vv.Content,
			Status:    &vv.Status,
			Thumbnail: vv.Thumbnail,
		}
		Group := service.GetVisualGroup(vv.UUID)
		if Group.UUID != "" {
			Vo.Gid = Group.UUID
		} else {
			Vo.Gid = ""
		}
		visuals = append(visuals, Vo)
	}
	c.JSON(common.HTTP_OK, common.OkWithData(visuals))

}

/*
*
* 大屏分组查看
*
 */
func ListVisualByGroup(c *gin.Context, ruleEngine typex.RuleX) {
	Gid, _ := c.GetQuery("uuid")
	visuals := []VisualVo{}
	MVisuals, _ := service.FindByType(Gid, "VISUAL")
	for _, vv := range MVisuals {
		Vo := VisualVo{
			UUID:      vv.UUID,
			Name:      vv.Name,
			Type:      vv.Type,
			Content:   vv.Content,
			Status:    &vv.Status,
			Thumbnail: vv.Thumbnail,
		}
		Group := service.GetVisualGroup(vv.UUID)
		Vo.Gid = Group.UUID
		visuals = append(visuals, Vo)
	}
	c.JSON(common.HTTP_OK, common.OkWithData(visuals))
}

/*
*
* 大屏详情
*
 */
func VisualDetail(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	mVisual, err := service.GetVisualWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	Vo := VisualVo{
		UUID:      mVisual.UUID,
		Name:      mVisual.Name,
		Type:      mVisual.Type,
		Content:   mVisual.Content,
		Status:    &mVisual.Status,
		Thumbnail: mVisual.Thumbnail,
	}
	Group := service.GetVisualGroup(mVisual.UUID)
	if Group.UUID != "" {
		Vo.Gid = Group.UUID
	} else {
		Vo.Gid = ""
	}
	c.JSON(common.HTTP_OK, common.OkWithData(Vo))
}

/*
*
* 上传缩略图
*
 */
func UploadFile(c *gin.Context, ruleEngine typex.RuleX) {
	// single file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	fileName := fmt.Sprintf("file_%d.png", time.Now().UnixMicro())
	dir := "./upload/VisualThumbnail/"
	if err := os.MkdirAll(filepath.Dir(dir), os.ModePerm); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err := c.SaveUploadedFile(file, dir+fileName); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(map[string]string{
		"url": fileName,
	}))
}

/*
*
* 加载缩略图
*
 */
func GetThumbnail(c *gin.Context, ruleEngine typex.RuleX) {
	fileName, _ := c.GetQuery("fileName")

	dir := "./upload/VisualThumbnail/"
	fileBytes, err := os.ReadFile(fmt.Sprintf("%s%s", dir, fileName))
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Header().Set("Content-Type", "image/jpeg")
	c.Writer.Header().Set("Content-Length", strconv.Itoa(len(fileBytes)))
	c.Writer.Write(fileBytes)
	c.Writer.Flush()
}
