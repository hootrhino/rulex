//
// Cache local
//
package x

var ruleCache map[string]*rule

func init() {
	ruleCache = map[string]*rule{}
}

//
//
//
func GetRule(id string) *rule {
	return ruleCache[id]
}

//
//
//
func SaveRule(r *rule) {
	ruleCache[r.Id] = r

}

//
//
//
func RemoveRule(r *rule) {
	delete(ruleCache, r.Id)
}

//
//
//
func AllRule() map[string]*rule {
	return ruleCache
}
