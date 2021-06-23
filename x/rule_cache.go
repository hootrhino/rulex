//
// Cache local
//
package x

var ruleCache map[string]*Rule

func init() {
	ruleCache = map[string]*Rule{}
}

//
//
//
func GetRule(id string) *Rule {
	return ruleCache[id]
}

//
//
//
func SaveRule(Rule *Rule) {
	ruleCache[Rule.Id] = Rule

}

//
//
//
func RemoveRule(Rule *Rule) {
	delete(ruleCache, Rule.Id)
}
