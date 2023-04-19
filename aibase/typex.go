package aibase

/*
*
* 算法模型
*
 */
type Algorithm struct {
	UUID        string // UUID
	Type        string // 模型类型: ANN_APP1 RNN_APP2 CNN_APP3 ....
	Name        string // 名称
	State       int    // 0开启;1关闭
	Document    string // 文档连接
	Description string // 概述
}

/*
*
* AI 接口
*
 */
type AlgorithmResource interface {
	Init(map[string]interface{}) error // 初始化环境
	// Type , Sample, ExpectOut
	Train(string, [][]float64, [][]float64) error      // 训练模型
	Load() error                                       // 加载模型
	OnCall(string, [][]float64) map[string]interface{} // 用数据去执行
	Unload() error                                     // 卸载模型
	AiDetail() Algorithm                               // 获取信息
}
