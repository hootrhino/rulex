package typex

/*
*
* 子进程的配置, 将 SocketAddr 传入 GRPC 客户端, Args 传入外挂的启动参数
*  $> /test_driver Args
 */
type Goods struct {
	UUID string
	// TCP or Unix Socket
	Addr string
	// Description text
	Description string
	// Additional Args
	Args []string
}
