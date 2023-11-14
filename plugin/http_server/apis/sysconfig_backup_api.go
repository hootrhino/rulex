package apis

import (
	"encoding/binary"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/ossupport"
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 备份Sqlite文件
*
 */
func DownloadSqlite(c *gin.Context, ruleEngine typex.RuleX) {
	wd, err := os.Getwd()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	fileName := "rulex.db"
	dir := wd
	c.Writer.WriteHeader(http.StatusOK)
	c.FileAttachment(fmt.Sprintf("%s/%s", dir, fileName),
		fmt.Sprintf("backup_%d_.db", time.Now().UnixNano()))
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
	dir := "./upload/Backup/"
	if err := os.MkdirAll(filepath.Dir(dir), os.ModePerm); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err := c.SaveUploadedFile(file, dir+fileName); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if _, err := ReadSQLiteFileMagicNumber(dir + fileName); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
	ossupport.StartRecoverProcess()

}

// https://www.sqlite.org/fileformat.html
func ReadSQLiteFileMagicNumber(filePath string) ([16]byte, error) {
	MagicNumber := [16]byte{}
	file, err := os.Open(filePath)
	if err != nil {
		return MagicNumber, err
	}
	defer file.Close()
	binary.Read(file, binary.BigEndian, &MagicNumber)
	if string(MagicNumber[:]) == "SQLite format 3\x00" {
		return MagicNumber, nil
	}
	return MagicNumber, fmt.Errorf("invalid Sqlite Db ,MagicNumber:%v error", MagicNumber)
}
