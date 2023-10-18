package apis

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/hootrhino/rulex/component/trailer"
	"github.com/hootrhino/rulex/glogger"
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/service"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gin-gonic/gin"
)

/*
*
* 停止正在运行的进程
*
 */
func StopGoods(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if goods := trailer.Get(uuid); goods != nil {
		if goods.Running {
			goods.Stop()
			c.JSON(common.HTTP_OK, common.Ok())
			return
		} else {
			c.JSON(common.HTTP_OK, common.Error("Already stopped"))
			return
		}
	}
	c.JSON(common.HTTP_OK, common.Error("Not exists:"+uuid))
}

/*
*
* 详情
*
 */
func GoodsDetail(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	mGood, err := service.GetGoodsWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	vo := goodsVo{
		Running:     false,
		Uuid:        mGood.UUID,
		LocalPath:   mGood.LocalPath,
		NetAddr:     mGood.NetAddr,
		Description: mGood.Description,
		Args:        mGood.Args,
	}
	if goods := trailer.Get(mGood.UUID); goods != nil {
		vo.Running = goods.Running
	}
	c.JSON(common.HTTP_OK, common.OkWithData(vo))

}

/*
*
* Goods
*
 */
type goodsVo struct {
	Running     bool     `json:"running"`
	Uuid        string   `json:"uuid"`
	LocalPath   string   `json:"local_path"`
	NetAddr     string   `json:"net_addr"`
	Description string   `json:"description"`
	Args        []string `json:"args"`
}

func GoodsList(c *gin.Context, ruleEngine typex.RuleX) {
	data := []goodsVo{}
	Goods := service.AllGoods()
	for _, mGood := range Goods {
		vo := goodsVo{
			Running:     false,
			Uuid:        mGood.UUID,
			LocalPath:   mGood.LocalPath,
			NetAddr:     mGood.NetAddr,
			Description: mGood.Description,
			Args:        mGood.Args,
		}
		if goods := trailer.Get(mGood.UUID); goods != nil {
			vo.Running = goods.Running
			data = append(data, vo)
		} else {
			data = append(data, vo)
		}
	}
	c.JSON(common.HTTP_OK, common.OkWithData(data))

}

/*
*
* 删除外挂
*
 */
func DeleteGoods(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	goods, err := service.GetGoodsWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 数据库和内存都要删除
	if err := service.DeleteGoods(goods.UUID); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	trailer.Remove(goods.UUID)
	// 删除文件
	if err := os.Remove(goods.LocalPath); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 清理垃圾文件: 从数据库里面筛选出所有的路径检查是不是和本地文件匹配，没用的直接删了
*
 */
func CleanGoodsUpload(c *gin.Context, ruleEngine typex.RuleX) {
	// 清理进程包
	if err := service.CleanGoodsUpload(); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* CreateGood
*
 */
func CreateGoods(c *gin.Context, ruleEngine typex.RuleX) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if fileHeader.Size > (100 << 20) {
		c.JSON(common.HTTP_OK, common.Error("File too large"))
		return
	}
	dir := "./upload/TrailerGoods/"
	OriginFileName := filepath.Base(fileHeader.Filename)
	fileExt := filepath.Ext(OriginFileName)
	fileName := fmt.Sprintf("goods_%d%s", time.Now().UnixMicro(), fileExt)
	if err := os.MkdirAll(filepath.Dir(dir), os.ModePerm); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err := c.SaveUploadedFile(fileHeader, dir+fileName); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	NetAddr := c.PostForm("net_addr")
	Description := c.PostForm("description")
	Args := c.PostFormArray("args")
	mGoods := model.MGoods{
		UUID:        utils.GoodsUuid(),
		LocalPath:   dir + fileName,
		NetAddr:     NetAddr,
		Description: Description,
		Args:        Args,
	}

	if err := service.InsertGoods(&mGoods); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	goods := trailer.Goods{
		UUID:        mGoods.UUID,
		LocalPath:   mGoods.LocalPath,
		NetAddr:     mGoods.NetAddr,
		Description: mGoods.Description,
		Args:        mGoods.Args,
	}
	if err := trailer.StartProcess(goods); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 更新操作
*
 */
func UpdateGoods(c *gin.Context, ruleEngine typex.RuleX) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if fileHeader.Size > (100 << 20) {
		c.JSON(common.HTTP_OK, common.Error("File too large"))
		return
	}
	Uuid := c.PostForm("uuid")
	NetAddr := c.PostForm("net_addr")
	Description := c.PostForm("description")
	Args := c.PostFormArray("args")
	//
	dir := "./upload/TrailerGoods/"
	OriginFileName := filepath.Base(fileHeader.Filename)
	fileExt := filepath.Ext(OriginFileName)
	fileName := fmt.Sprintf("goods_%d%s", time.Now().UnixMicro(), fileExt)
	if err := os.MkdirAll(filepath.Dir(dir), os.ModePerm); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err := c.SaveUploadedFile(fileHeader, dir+fileName); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	mGoods := model.MGoods{
		UUID:        Uuid,
		LocalPath:   dir + fileName,
		NetAddr:     NetAddr,
		Description: Description,
		Args:        Args,
	}
	err1 := service.UpdateGoods(mGoods)
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	// 把正在运行的给停了
	if goods := trailer.Get(mGoods.UUID); goods != nil {
		if err != nil {
			glogger.GLogger.Error(err)
			return
		}
		glogger.GLogger.Debug("Already running, ready to stop:", mGoods.UUID)
		grpcConnection, err1 := grpc.Dial(goods.NetAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err1 != nil {
			return
		}
		defer grpcConnection.Close()
		client := trailer.NewTrailerClient(grpcConnection)
		client.Stop(context.Background(), &trailer.Request{})
		trailer.Remove(mGoods.UUID)
	}

	// 开新进程
	goods := trailer.Goods{
		UUID:        mGoods.UUID,
		LocalPath:   mGoods.LocalPath,
		NetAddr:     mGoods.NetAddr,
		Description: mGoods.Description,
		Args:        mGoods.Args,
	}
	if err := trailer.StartProcess(goods); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 上传文件，保存在 "./upload/goods/" 路径
*
 */
func UploadGoodsFile(c *gin.Context, ruleEngine typex.RuleX) {
	// single file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if runtime.GOOS == "windows" {
		if !IsExecutableFileWin(file.Filename) {
			c.JSON(common.HTTP_OK, common.Error("Invalid execute file"))
			return
		}
	}
	if runtime.GOOS == "linux" {
		if !IsExecutableFileUnix(file.Filename) ||
			IsExecutableFileWin(file.Filename) {
			c.JSON(common.HTTP_OK, common.Error("Invalid execute file"))
			return
		}
	}
	dir := "./upload/TrailerGoods/"
	OriginFileName := filepath.Base(file.Filename)
	fileExt := filepath.Ext(OriginFileName)
	fileName := fmt.Sprintf("goods_%d_%s", time.Now().UnixMicro(), fileExt)
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
* 判断是否可执行(Linux Only)
*
 */
func IsExecutableFileUnix(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	if fileInfo.Mode()&0111 != 0 {
		return true
	}

	return false
}
func IsExecutableFileWin(filePath string) bool {
	filePath = strings.ToLower(filePath)
	return strings.HasSuffix(filePath, ".exe") ||
		strings.HasSuffix(filePath, ".jar") ||
		strings.HasSuffix(filePath, ".py") ||
		strings.HasSuffix(filePath, ".js") ||
		strings.HasSuffix(filePath, ".lua")

}
