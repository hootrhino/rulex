package apis

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/component/interdb"
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 备份Sqlite文件
*
 */
func BackupSqlite(c *gin.Context, ruleEngine typex.RuleX) {
	wd, err := os.Getwd()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	fileName := "rulex.db"
	dir := wd
	c.Writer.WriteHeader(http.StatusOK)
	c.FileAttachment(fmt.Sprintf("%s/%s", dir, fileName),
	fmt.Sprintf("backup_%d_.db",time.Now().UnixNano()))
}

/*
*
* 上传恢复
*
 */
func UploadSqlite(c *gin.Context, ruleEngine typex.RuleX) {
	// single file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	fileName := "recovery.db"
	dir := "./upload/backup/"
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
	interdb.DB().Migrator()
}
