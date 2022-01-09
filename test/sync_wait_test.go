package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

//
// 起一个测试环境
// docker run -it --link test-redis:redis --rm redis redis-cli -h redis -p 6379
//
/*
*
* 测试开关打开或者关闭后状态同步机制
*
 */
func Test_Open_Switch(t *testing.T) {
	var redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	defer redisClient.Close()
	//
	//
	//
	requestId := "request-id-001"
	sendCmd(ctx, redisClient, requestId)
	waitResult(ctx, redisClient, requestId)
	time.Sleep(5 * time.Second)

}

/*
*
* 发送指令:当指令下发后马上给redis保存一个指令id, 用于等待后期同步
 */
func sendCmd(ctx context.Context, redisClient *redis.Client, requestId string) {
	fmt.Println("Send open cmd to rulex")
	redisClient.Set(ctx, requestId, 0, 5*time.Second)
}

/*
*
* 等待执行结果
*
 */
func waitResult(ctx context.Context, redisClient *redis.Client, requestId string) {
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				{
					failedCmd(ctx, redisClient, requestId)
					return
				}
			default:
				{
					s := redisClient.Get(ctx, requestId)
					if s.Err() != nil && s.Val() != "" {
						if ok, _ := s.Bool(); ok {
							finishCmd(ctx, redisClient, requestId)
							return
						}
					}
				}
			}
		}
	}(ctx)

}

/*
*
*监听rulex的反馈, 如果  rulexlib:finishCmd(CmdId) 调用了 这里就把redis的值更新
*
 */
func finishCmd(ctx context.Context, redisClient *redis.Client, requestId string) {
	println("finished:" + requestId)

}

/*
*
*
*
 */
func failedCmd(ctx context.Context, redisClient *redis.Client, requestId string) {
	println("failed:" + requestId)
}
