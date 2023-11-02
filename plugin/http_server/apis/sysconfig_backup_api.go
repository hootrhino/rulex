package apis

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 备份Sqlite文件
*
 */
func BackupSqlite(c *gin.Context, ruleEngine typex.RuleX) {
	fileName := "backup.sql"
	dir := "./upload/backup/"
	fileBytes, err := os.ReadFile(fmt.Sprintf("%s%s", dir, fileName))
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.Header().Set("Content-Length", strconv.Itoa(len(fileBytes)))
	c.Writer.Write(fileBytes)
	c.Writer.Flush()
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
	fileName := "backup.sql"
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
}
