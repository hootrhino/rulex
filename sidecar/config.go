package sidecar

/*
*
* 子进程的配置
*
 */
type GoodsConfig struct {
	// TCP or Unix Socket
	SocketAddr string
	// Description text
	Description string
}
