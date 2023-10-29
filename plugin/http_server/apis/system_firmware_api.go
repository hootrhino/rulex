// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package apis

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/component/trailer"
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/plugin/http_server/service"
	"github.com/hootrhino/rulex/typex"
)

/*
*
  - 上传最新固件, 必须是ZIP包, 固件保存在:./upload/Firmware/Firmware.zip
    压缩包内就是rulex发布的最新版本

*
*/
func UploadFirmWare(c *gin.Context, ruleEngine typex.RuleX) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	saveDir := "./upload/Firmware/"
	fileName := "Firmware.zip" // 固定名称
	if err := os.MkdirAll(filepath.Dir(saveDir), os.ModePerm); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err := c.SaveUploadedFile(file, saveDir+fileName); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	c.JSON(common.HTTP_OK, common.OkWithData(saveDir+fileName))
}

/*
*
* 解压、升级
*
 */
func UpgradeFirmWare(c *gin.Context, ruleEngine typex.RuleX) {
	dir := "./upload/Firmware/"
	fileName := "Firmware.zip" // 固定名称
	tempPath, err := os.MkdirTemp(dir+fileName+"/temp", "temp*")
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err := trailer.Unzip(dir+fileName, tempPath); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 检查 /tmp/temp227209938/rulex 的Md5
	md51, err1 := sumMD5(tempPath + "/rulex")
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	// 从解压后的目录提取Md5
	readBytes, err2 := os.ReadFile(tempPath + "/md5.sum")
	if err2 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err2))
		return
	}
	if md51 != string(readBytes) {
		c.JSON(common.HTTP_OK, common.Error("invalid sum md5!"))
		return
	}
	// Rulex 默认安装在 /usr/local/ 目录
	// cp -> rulex-arm32-linux "/usr/local/"
	if err := os.Rename(tempPath+"/rulex*", "/usr/local/"); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err := service.ReStartRulex(); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
  - 检查包, 一般包里面会有一个可执行文件和 MD5 SUM 值。要对比一下。
    文件列表:
  - rulex
  - rulex.ini
  - md5.sum

*
*/

func sumMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	_, _ = io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil)), nil
}
