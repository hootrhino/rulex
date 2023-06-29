package httpserver

import (
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2/json"
)

func OutEnds(c *gin.Context, hh *HttpApiServer) {
	uuid, _ := c.GetQuery("uuid")
	if uuid == "" {
		outends := []typex.OutEnd{}
		for _, mOut := range hh.AllMOutEnd() {
			outEnd := hh.ruleEngine.GetOutEnd(mOut.UUID)
			if outEnd == nil {
				tOut := typex.OutEnd{}
				tOut.UUID = mOut.UUID
				tOut.Name = mOut.Name
				tOut.Type = typex.TargetType(mOut.Type)
				tOut.Description = mOut.Description
				tOut.Config = mOut.GetConfig()
				tOut.State = typex.SOURCE_STOP
				outends = append(outends, tOut)
			}
			if outEnd != nil {
				outEnd.State = outEnd.Target.Status()
				outends = append(outends, *outEnd)
			}
		}
		c.JSON(common.HTTP_OK, common.OkWithData(outends))
		return
	}
	mOut, err := hh.GetMOutEndWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	outEnd := hh.ruleEngine.GetOutEnd(mOut.UUID)
	if outEnd == nil {
		// 如果内存里面没有就给安排一个死设备
		tOut := typex.OutEnd{}
		tOut.UUID = mOut.UUID
		tOut.Name = mOut.Name
		tOut.Type = typex.TargetType(mOut.Type)
		tOut.Description = mOut.Description
		tOut.Config = mOut.GetConfig()
		tOut.State = typex.SOURCE_STOP
		c.JSON(common.HTTP_OK, common.OkWithData(tOut))
		return
	}
	outEnd.State = outEnd.Target.Status()
	c.JSON(common.HTTP_OK, common.OkWithData(outEnd))
}

// Get all outends
func OutEndDetail(c *gin.Context, hs *HttpApiServer) {
	uuid, _ := c.GetQuery("uuid")
	mOut, err := hs.GetMOutEndWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400EmptyObj(err))
		return
	}
	outEnd := hs.ruleEngine.GetOutEnd(mOut.UUID)
	if outEnd == nil {
		// 如果内存里面没有就给安排一个死设备
		tOutEnd := new(typex.OutEnd)
		tOutEnd.UUID = mOut.UUID
		tOutEnd.Name = mOut.Name
		tOutEnd.Type = typex.TargetType(mOut.Type)
		tOutEnd.Description = mOut.Description
		tOutEnd.Config = mOut.GetConfig()
		tOutEnd.State = typex.SOURCE_STOP
		c.JSON(common.HTTP_OK, common.OkWithData(tOutEnd))
		return
	}
	outEnd.State = outEnd.Target.Status()
	c.JSON(common.HTTP_OK, common.OkWithData(outEnd))
}

// Delete outEnd by UUID
func DeleteOutEnd(c *gin.Context, hs *HttpApiServer) {
	uuid, _ := c.GetQuery("uuid")
	_, err := hs.GetMOutEndWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	if err := hs.DeleteMOutEnd(uuid); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	old := hs.ruleEngine.GetOutEnd(uuid)
	if old != nil {
		if old.Target.Status() == typex.SOURCE_UP {
			old.Target.Details().State = typex.SOURCE_STOP
			old.Target.Stop()
		}
	}
	hs.ruleEngine.RemoveOutEnd(uuid)
	c.JSON(common.HTTP_OK, common.Ok())
}

// Create or Update OutEnd
func CreateOutEnd(c *gin.Context, hh *HttpApiServer) {
	type Form struct {
		UUID        string                 `json:"uuid"` // 如果空串就是新建, 非空就是更新
		Type        string                 `json:"type" binding:"required"`
		Name        string                 `json:"name" binding:"required"`
		Description string                 `json:"description"`
		Config      map[string]interface{} `json:"config" binding:"required"`
	}
	form := Form{}

	if err0 := c.ShouldBindJSON(&form); err0 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err0))
		return
	}
	configJson, err1 := json.Marshal(form.Config)
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	newUUID := utils.OutUuid()
	if err := hh.InsertMOutEnd(&MOutEnd{
		UUID:        newUUID,
		Type:        form.Type,
		Name:        form.Name,
		Description: form.Description,
		Config:      string(configJson),
	}); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err := hh.LoadNewestOutEnd(newUUID); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

// 更新
func UpdateOutEnd(c *gin.Context, hs *HttpApiServer) {
	type Form struct {
		UUID        string                 `json:"uuid"` // 如果空串就是新建, 非空就是更新
		Type        string                 `json:"type" binding:"required"`
		Name        string                 `json:"name" binding:"required"`
		Description string                 `json:"description"`
		Config      map[string]interface{} `json:"config" binding:"required"`
	}
	form := Form{}

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	configJson, err := json.Marshal(form.Config)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if form.UUID == "" {
		c.JSON(common.HTTP_OK, common.Error("missing 'uuid' fields"))
		return
	}
	// 更新的时候从数据库往外面拿
	OutEnd, err := hs.GetMOutEndWithUUID(form.UUID)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	if err := hs.UpdateMOutEnd(OutEnd.UUID, &MOutEnd{
		UUID:        form.UUID,
		Type:        form.Type,
		Name:        form.Name,
		Description: form.Description,
		Config:      string(configJson),
	}); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	if err := hs.LoadNewestOutEnd(form.UUID); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	c.JSON(common.HTTP_OK, common.Ok())
}
