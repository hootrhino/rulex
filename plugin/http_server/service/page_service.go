package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"gorm.io/gorm"
	"strconv"
)

func Paginate(page model.PageRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		current := page.Current
		size := page.Size

		offset := (current - 1) * size
		return db.Offset(offset).Limit(size)
	}
}

func WrapPageResult(page model.PageRequest, records any, count int64) model.PageResult {
	return model.PageResult{
		Current: page.Current,
		Size:    page.Size,
		Total:   int(count),
		Records: records,
	}
}

func ReadPageRequest(c *gin.Context) (*model.PageRequest, error) {
	page := &model.PageRequest{}
	var err error
	page.Current, err = strconv.Atoi(c.DefaultQuery("current", "1"))
	if err != nil {
		return nil, err
	}
	page.Size, err = strconv.Atoi(c.DefaultQuery("size", "25"))
	if err != nil {
		return nil, err
	}
	page.SearchCount, err = strconv.Atoi(c.DefaultQuery("searchCount", "1"))
	if err != nil {
		return nil, err
	}
	return page, nil
}
