package typex

type TargetRegistry interface {
	Register(string, func(RuleX) XTarget)

	Find(string) func(RuleX) XTarget
}
