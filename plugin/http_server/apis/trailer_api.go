package apis

import (
	"context"
	"debug/elf"
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
	Uuid          string      `json:"uuid"`
	Pid           int         `json:"pid"`
	Running       bool        `json:"running"`
	AutoStart     bool        `json:"autoStart"`
	GoodsType     string      `json:"goodsType"`   // LOCAL, EXTERNAL
	ExecuteType   string      `json:"executeType"` // exe,elf,js,py....
	LocalPath     string      `json:"local_path"`
	NetAddr       string      `json:"net_addr"`
	Description   string      `json:"description"`
	Args          []string    `json:"args"`
	ProcessDetail interface{} `json:"processDetail"`
}

/*
*
* 停止正在运行的进程
*
 */
func StopGoods(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if goods := trailer.Get(uuid); goods != nil {
		if goods.PsRunning() {
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
		GoodsType:   mGood.GoodsType,
		ExecuteType: mGood.ExecuteType,
		AutoStart:   *mGood.AutoStart,
		LocalPath:   mGood.LocalPath,
		NetAddr:     mGood.NetAddr,
		Args:        []string{mGood.Args},
		Description: mGood.Description,
	}
	if goods := trailer.Get(mGood.UUID); goods != nil {
		vo.Running = goods.PsRunning()
		vo.Pid = goods.Pid()
		detail, _ := trailer.RunningProcessDetail(goods.Pid())
		vo.ProcessDetail = detail
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
			GoodsType:   mGood.GoodsType,
			ExecuteType: mGood.ExecuteType,
			LocalPath:   mGood.LocalPath,
			NetAddr:     mGood.NetAddr,
			Args:        []string{mGood.Args},
			Description: mGood.Description,
		}
		if goods := trailer.Get(mGood.UUID); goods != nil {
			vo.Running = goods.PsRunning()
			vo.Pid = goods.Pid()
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
	mGoods, err := service.GetGoodsWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	if goods := trailer.Get(uuid); goods != nil {
		trailer.RemoveBy(goods.Info.UUID, "RULEX")
	}
	// 数据库和内存都要删除
	if err := service.DeleteGoods(mGoods.UUID); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 删除文件
	os.Remove(mGoods.LocalPath)
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

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		// 出现错误
		return false, err
	}
}

/*
*
* 新建一个扩展
*
 */
func CreateGoods(c *gin.Context, ruleEngine typex.RuleX) {
	fileHeader, err := c.FormFile("file")
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
	Os := runtime.GOOS
	if Os == "linux" {
		if fileExt == ".exe" {
			c.JSON(common.HTTP_OK, common.Error("Linux not support windows execute format"))
			return
		}
	}

	fileName := fmt.Sprintf("goods_%d%s", time.Now().UnixMicro(), fileExt)
	// 目录是否存在
	path := filepath.Dir(__TrailerGoodsUploadDir)
	if exists, _ := PathExists(path); !exists {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}

	localSavePath := __TrailerGoodsUploadDir + fileName
	if err := c.SaveUploadedFile(fileHeader, localSavePath); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	ExeType := getExecuteType(localSavePath)
	if ExeType == "" {
		c.JSON(common.HTTP_OK,
			common.Error("Invalid file:"+localSavePath))
		os.Remove(localSavePath)
		return
	}
	if ExeType == "ELF" {
		if Os == "linux" {
			localSavePath += ".elfx" // 标记是Linux可执行
		}
	}

	if Os == "windows" {
		if ExeType == "ELF" {
			c.JSON(common.HTTP_OK, common.Error("Windows not support linux ELF format"))
			os.Remove(localSavePath)
			return
		}
	}
	if Os == "linux" {
		if ExeType == "EXE" {
			c.JSON(common.HTTP_OK, common.Error("Linux not support Windows EXE format"))
			os.Remove(localSavePath)
			return
		}
	}
	NetAddr := c.PostForm("net_addr")
	Description := c.PostForm("description")
	AutoStart := c.PostForm("autoStart")
	Args := c.PostFormArray("args")
	mGoods := model.MGoods{
		UUID:      utils.GoodsUuid(),
		LocalPath: localSavePath,
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
		ExecuteType: ExeType,
		GoodsType:   "LOCAL", // 默认是LOCAL, 未来根据前端参数决定
		Args: func() string {
			if len(Args) > 0 {
				return Args[0]
			}
			return ""
		}(),
		Description: Description,
	}

	if err := service.InsertGoods(&mGoods); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	goods := trailer.GoodsInfo{
		UUID:        mGoods.UUID,
		AutoStart:   mGoods.AutoStart,
		LocalPath:   mGoods.LocalPath,
		GoodsType:   mGoods.GoodsType,
		ExecuteType: mGoods.ExecuteType,
		NetAddr:     mGoods.NetAddr,
		Args:        mGoods.Args,
		Description: mGoods.Description,
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
	if fileHeader.Size > (100 << 20) {
		c.JSON(common.HTTP_OK, common.Error("File too large"))
		return
	}
	Uuid := c.PostForm("uuid")
	NetAddr := c.PostForm("net_addr")
	Description := c.PostForm("description")
	Args := c.PostFormArray("args")
	AutoStart := c.PostForm("autoStart")
	OriginFileName := filepath.Base(fileHeader.Filename)
	fileExt := filepath.Ext(OriginFileName)
	Os := runtime.GOOS
	if Os == "linux" {
		if fileExt == ".exe" {
			c.JSON(common.HTTP_OK, common.Error("Linux not support windows execute format"))
			return
		}
	}
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
	ExeType := getExecuteType(localSavePath)
	if ExeType == "" {
		c.JSON(common.HTTP_OK,
			common.Error("Invalid file:"+localSavePath))
		os.Remove(localSavePath)
		return
	}
	if ExeType == "ELF" {
		if Os == "linux" {
			localSavePath += ".elfx" // 标记是Linux可执行
		}
	}
	if Os == "windows" {
		if ExeType == "ELF" {
			c.JSON(common.HTTP_OK, common.Error("Windows not support linux ELF format"))
			os.Remove(localSavePath)
			return
		}
	}
	if Os == "linux" {
		if ExeType == "EXE" {
			c.JSON(common.HTTP_OK, common.Error("Linux not support Windows EXE format"))
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
		ExecuteType: ExeType,
		Args: func() string {
			if len(Args) > 0 {
				return Args[0]
			}
			return ""
		}(),
		Description: Description,
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
		grpcConnection, err1 := grpc.Dial(goods.Info.NetAddr,
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
	goods := trailer.GoodsInfo{
		UUID:        mGoods.UUID,
		AutoStart:   mGoods.AutoStart,
		LocalPath:   mGoods.LocalPath,
		NetAddr:     mGoods.NetAddr,
		Args:        mGoods.Args,
		GoodsType:   mGoods.GoodsType,
		ExecuteType: mGoods.ExecuteType,
		Description: mGoods.Description,
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
	goods := trailer.GoodsInfo{
		UUID:        mGoods.UUID,
		AutoStart:   mGoods.AutoStart,
		LocalPath:   mGoods.LocalPath,
		NetAddr:     mGoods.NetAddr,
		Args:        mGoods.Args,
		GoodsType:   mGoods.GoodsType,
		ExecuteType: mGoods.ExecuteType,
		Description: mGoods.Description,
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
	if IsUnixElf(filePath) {
		ChangeX(filePath)
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
	return true
}

/*
*
* 是否是可执行脚本语言
*
 */
func IsExecutableScript(fileExt string) bool {
	return true
}

/*
*
* 是否是可执行Linux文件
*
 */
func IsUnixElf(filePath string) bool {
	file, err := elf.Open(filePath)
	if err != nil {
		glogger.GLogger.Error(err)
		return false
	}
	defer file.Close()
	return true
}
func IsWinPE(filePath string) bool {
	file, err := pe.Open(filePath)
	if err != nil {
		glogger.GLogger.Error(err)
		return false
	}
	defer file.Close()
	return true
}

/*
*
* 给Linux下 ELF 文件增加可执行权限
*
 */
func ChangeX(filePath string) error {
	// 打开 ELF 文件
	file, err := elf.Open(filePath)
	if err != nil {
		return err
	}
	elfHeader := file.FileHeader
	file.Close()
	if elfHeader.Type == elf.ET_EXEC {
		// 设置可执行权限 (0700 表示读、写、执行权限)
		if err := os.Chmod(filePath, 0755); err != nil {
			return err
		}
		return nil
	}
	return nil
}

/*
*
* 获取文件类型
*
 */
func getExecuteType(OriginFileName string) string {
	fileExt := filepath.Ext(OriginFileName)
	if v, ok := trailer.ExecuteType[fileExt]; ok {
		return v
	}
	if IsUnixElf(OriginFileName) {
		return "ELF" // Maybe ELF
	}
	return "" // error
}
