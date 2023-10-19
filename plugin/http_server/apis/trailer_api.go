package apis

import (
	"context"
	"debug/pe"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
* Goods
*
 */
type goodsVo struct {
	AutoStart   bool     `json:"autoStart"`
	Pid         int      `json:"pid"`
	Running     bool     `json:"running"`
	Uuid        string   `json:"uuid"`
	LocalPath   string   `json:"local_path"`
	NetAddr     string   `json:"net_addr"`
	Description string   `json:"description"`
	Args        []string `json:"args"`
}

/*
*
* 停止正在运行的进程
*
 */
func StopGoods(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if goods := trailer.Get(uuid); goods != nil {
		if goods.PsRunning {
			goods.StopBy("RULEX")
			c.JSON(common.HTTP_OK, common.Ok())
			return
		}
		c.JSON(common.HTTP_OK, common.Error("Already stopped"))
		return
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
		AutoStart:   *mGood.AutoStart,
		LocalPath:   mGood.LocalPath,
		NetAddr:     mGood.NetAddr,
		Description: mGood.Description,
		Args:        []string{mGood.Args},
	}
	if goods := trailer.Get(mGood.UUID); goods != nil {
		vo.Running = goods.PsRunning
		vo.Pid = goods.Pid
	}
	c.JSON(common.HTTP_OK, common.OkWithData(vo))

}

func GoodsList(c *gin.Context, ruleEngine typex.RuleX) {
	data := []goodsVo{}
	Goods := service.AllGoods()
	for _, mGood := range Goods {
		vo := goodsVo{
			Running:     false,
			Uuid:        mGood.UUID,
			AutoStart:   *mGood.AutoStart,
			LocalPath:   mGood.LocalPath,
			NetAddr:     mGood.NetAddr,
			Description: mGood.Description,
			Args:        []string{mGood.Args},
		}
		if goods := trailer.Get(mGood.UUID); goods != nil {
			vo.Running = goods.PsRunning
			vo.Pid = goods.Pid
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

	if goods := trailer.Get(uuid); goods != nil {
		if goods.PsRunning {
			trailer.RemoveBy(goods.Uuid, "RULEX")
		}
	}
	// 数据库和内存都要删除
	if err := service.DeleteGoods(goods.UUID); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 删除文件
	os.Remove(goods.LocalPath)
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
var __TrailerGoodsUploadDir = "./upload/TrailerGoods/"

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
	OriginFileName := filepath.Base(fileHeader.Filename)
	fileExt := filepath.Ext(OriginFileName)
	fileName := fmt.Sprintf("goods_%d%s", time.Now().UnixMicro(), fileExt)
	if err := os.MkdirAll(filepath.Dir(__TrailerGoodsUploadDir), os.ModePerm); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	localSavePath := __TrailerGoodsUploadDir + fileName
	if err := c.SaveUploadedFile(fileHeader, localSavePath); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 现阶段先只支持这俩系统
	if runtime.GOOS == "windows" {
		if !IsExecutableFileWin(localSavePath) {
			c.JSON(common.HTTP_OK,
				common.Error("Is not windows Executable File:"+localSavePath))
			os.Remove(localSavePath)
			return
		}
	}
	if runtime.GOOS == "linux" {
		if !IsExecutableFileUnix(localSavePath) {
			c.JSON(common.HTTP_OK,
				common.Error("Is not Linux(Unix) Executable File:"+localSavePath))
			os.Remove(localSavePath)
			return
		}
	}

	NetAddr := c.PostForm("net_addr")
	Description := c.PostForm("description")
	AutoStart := c.PostForm("AutoStart")
	Args := c.PostFormArray("args")
	mGoods := model.MGoods{
		UUID:      utils.GoodsUuid(),
		LocalPath: __TrailerGoodsUploadDir + fileName,
		NetAddr:   NetAddr,
		AutoStart: func() *bool {
			if AutoStart == "1" ||
				AutoStart == "true" {
				r := true
				return &r
			}
			r := false
			return &r
		}(),
		Description: Description,
		Args: func() string {
			if len(Args) > 0 {
				return Args[0]
			}
			return ""
		}(),
	}

	if err := service.InsertGoods(&mGoods); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	goods := trailer.Goods{
		UUID:        mGoods.UUID,
		AutoStart:   mGoods.AutoStart,
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
	AutoStart := c.PostForm("AutoStart")
	OriginFileName := filepath.Base(fileHeader.Filename)
	fileExt := filepath.Ext(OriginFileName)
	fileName := fmt.Sprintf("goods_%d%s", time.Now().UnixMicro(), fileExt)
	if err := os.MkdirAll(filepath.Dir(__TrailerGoodsUploadDir), os.ModePerm); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	localSavePath := __TrailerGoodsUploadDir + fileName
	if err := c.SaveUploadedFile(fileHeader, localSavePath); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 现阶段先只支持这俩系统
	if runtime.GOOS == "windows" {
		if !IsExecutableFileWin(localSavePath) {
			c.JSON(common.HTTP_OK,
				common.Error("Is not windows Executable File:"+localSavePath))
			os.Remove(localSavePath)
			return
		}
	}
	if runtime.GOOS == "linux" {
		if !IsExecutableFileUnix(localSavePath) {
			c.JSON(common.HTTP_OK,
				common.Error("Is not Linux(Unix) Executable File:"+localSavePath))
			os.Remove(localSavePath)
			return
		}
	}

	mGoods := model.MGoods{
		UUID: Uuid,
		AutoStart: func() *bool {
			if AutoStart == "1" ||
				AutoStart == "true" {
				r := true
				return &r
			}
			r := false
			return &r
		}(),
		LocalPath:   localSavePath,
		NetAddr:     NetAddr,
		Description: Description,
		Args: func() string {
			if len(Args) > 0 {
				return Args[0]
			}
			return ""
		}(),
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
		AutoStart:   mGoods.AutoStart,
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
* 尝试启动已经停止的进程
*
 */
func StartGoods(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	mGoods, err := service.GetGoodsWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if oldPs := trailer.Get(uuid); oldPs != nil {
		c.JSON(common.HTTP_OK, common.Error("Already started:"+uuid))
		return
	}
	// 开新进程
	goods := trailer.Goods{
		UUID:        mGoods.UUID,
		AutoStart:   mGoods.AutoStart,
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

/*
*
* 读取PE头判断是否可执行
*
 */
func IsExecutableFileWin(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		glogger.GLogger.Error(err)
		return false
	}
	defer file.Close()

	if _, err := pe.NewFile(file); err != nil {
		glogger.GLogger.Error(err)
		return false
	}
	return true
}
