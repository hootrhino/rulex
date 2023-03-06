package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/i4de/rulex/typex"
)

// 列表
func Apps(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if uuid == "" {
		c.JSON(200, OkWithData(e.AllApp()))
		return
	}
	c.JSON(200, OkWithData(e.GetApp(uuid)))

}

// 新建
func CreateApp(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	c.JSON(200, OkWithData("ok"))
}

// 更新
func UpdateApp(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	c.JSON(200, OkWithData("ok"))
}

// 停止
func StopApp(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	c.JSON(200, e.StopApp(uuid))
}

// 删除
func RemoveApp(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	c.JSON(200, e.RemoveApp(uuid))
}

//-------------------------------------------------------------------------------------
// App Dao
//-------------------------------------------------------------------------------------

// 获取App列表
func (s *HttpApiServer) AllApp() []MApp {
	m := []MApp{}
	s.sqliteDb.Find(&m)
	return m

}
func (s *HttpApiServer) GetAppWithUUID(uuid string) (*MApp, error) {
	m := MApp{}
	if err := s.sqliteDb.Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除App
func (s *HttpApiServer) DeleteApp(uuid string) error {
	return s.sqliteDb.Where("uuid=?", uuid).Delete(&MApp{}).Error
}

// 创建App
func (s *HttpApiServer) InsertApp(goods *MApp) error {
	return s.sqliteDb.Table("m_goods").Create(goods).Error
}

// 更新App
func (s *HttpApiServer) UpdateApp(uuid string, goods *MApp) error {
	m := MApp{}
	if err := s.sqliteDb.Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		s.sqliteDb.Model(m).Updates(*goods)
		return nil
	}
}
