package cloud

type ServiceArg struct {
	Value interface{}
}
type CallResult struct {
	Code int
	Msg  string
	Data []interface{}
}
type CloudService struct {
	Args []ServiceArg
}
type Cloud interface {
	ListService(pageIndex int, pageSize int) []CloudService
	CallService(id string, args []ServiceArg) CallResult
}
