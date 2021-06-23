//
// Cache local
//
package x

var inEndCache map[string]*InEnd

func init() {
	inEndCache = map[string]*InEnd{}
}

//
//
//
func GetInEnd(id string) *InEnd {
	return inEndCache[id]
}

//
//
//
func SaveInEnd(inEnd *InEnd) {
	inEndCache[inEnd.Id] = inEnd
}

//
//
//
func RemoveInEnd(inEnd *InEnd) {
	delete(inEndCache, inEnd.Id)
}
