package typex

/*
*
* TODO: 定义一系列接口，以后每个涉及到CURD管理的都实现这个接口
* 暂时先放在这里作为参考，短期内不会重构
*
 */
type XService interface {
	Load()
	Get()
	List()
	Remove()
	Stop()
}
