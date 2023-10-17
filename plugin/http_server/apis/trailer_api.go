package apis

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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

func GoodsDetail(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if goods := trailer.Get(uuid); goods != nil {
		c.JSON(common.HTTP_OK, common.OkWithData(goodsVo{
			Running:     goods.Running,
			Uuid:        goods.Uuid,
			LocalPath:   goods.LocalPath,
			NetAddr:     goods.NetAddr,
			Description: goods.Description,
			Args:        goods.Args,
		}))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(goodsVo{}))
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
	} else {
		// 数据库和内存都要删除
		service.DeleteGoods(goods.UUID)
		trailer.Remove(goods.UUID)
		c.JSON(common.HTTP_OK, common.Ok())
	}
}

/*
*
* CreateGood
*
 */
func CreateGoods(c *gin.Context, ruleEngine typex.RuleX) {
	form := goodsVo{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	mGoods := model.MGoods{
		UUID:        utils.GoodsUuid(),
		LocalPath:   form.LocalPath,
		NetAddr:     form.NetAddr,
		Description: form.Description,
		Args:        form.Args,
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
	if err := trailer.Fork(goods); err != nil {
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
	form := goodsVo{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	mGoods := model.MGoods{
		UUID:        form.Uuid,
		LocalPath:   form.LocalPath,
		NetAddr:     form.NetAddr,
		Description: form.Description,
		Args:        form.Args,
	}
	err := service.UpdateGoods(mGoods)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
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
	if err := trailer.Fork(goods); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 上传缩略图
*
 */
func UploadGoodsFile(c *gin.Context, ruleEngine typex.RuleX) {
	// single file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	fileName := fmt.Sprintf("goods_%d", time.Now().UnixMicro())
	dir := "./resource/goods/"
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
