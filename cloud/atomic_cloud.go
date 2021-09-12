// Atomic 平台云端服务
//
//
package cloud

type AtomicCloud struct {
}

func (a *AtomicCloud) ListService(pageIndex int, pageSize int) []CloudService {

	return []CloudService{}
}
func (a *AtomicCloud) CallService(id string, args []ServiceArg) CallResult {
	return CallResult{}
}
