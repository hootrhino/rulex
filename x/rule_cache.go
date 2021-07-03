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
	lock.Lock()
	defer lock.Unlock()
	return ruleCache[id]
}

//
//
//
func SaveRule(r *rule) {
	lock.Lock()
	defer lock.Unlock()
	ruleCache[r.Id] = r

}

//
//
//
func RemoveRule(ruleId string) {
	lock.Lock()
	defer lock.Unlock()
	delete(ruleCache, ruleId)
}

//
//
//
func AllRule() map[string]*rule {
	lock.Lock()
	defer lock.Unlock()
	return ruleCache
}
