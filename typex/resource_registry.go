package typex

type ResourceRegistry interface {
	Register(string, func(RuleX) XResource)

	Find(string) func(RuleX) XResource
}
