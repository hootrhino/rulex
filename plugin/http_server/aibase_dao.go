package httpserver

// 获取AiBase列表
func (s *HttpApiServer) AllAiBase() []MAiBase {
	m := []MAiBase{}
	s.sqliteDb.Find(&m)
	return m

}
func (s *HttpApiServer) GetAiBaseWithUUID(uuid string) (*MAiBase, error) {
	m := MAiBase{}
	if err := s.sqliteDb.Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除AiBase
func (s *HttpApiServer) DeleteAiBase(uuid string) error {
	return s.sqliteDb.Where("uuid=?", uuid).Delete(&MAiBase{}).Error
}

// 创建AiBase
func (s *HttpApiServer) InsertAiBase(AiBase *MAiBase) error {
	return s.sqliteDb.Create(AiBase).Error
}

// 更新AiBase
func (s *HttpApiServer) UpdateAiBase(AiBase *MAiBase) error {
	m := MAiBase{}
	if err := s.sqliteDb.Where("uuid=?", AiBase.UUID).First(&m).Error; err != nil {
		return err
	} else {
		s.sqliteDb.Model(m).Updates(*AiBase)
		return nil
	}
}
