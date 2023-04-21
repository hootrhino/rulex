package typex

type XAiBase interface {
	GetRuleX() RuleX
	ListAi() []*AiBase
	LoadAi(Ai *AiBase) error
	GetAi(uuid string) *AiBase
	RemoveAi(uuid string) error
	UpdateAi(Ai AiBase) error
	StartAi(uuid string) error
	StopAi(uuid string) error
	Stop()
}
type AiBase struct {
	UUID        string // 名称
	Name        string // 名称
	Type        string // 类型
	Filepath    string // 文件路径, 是相对于main的aibase目录
	Description string
}

/*
*
* 生成AI应用
*
 */
func NewAiBase(UUID, Name, Type, Filepath, Description string) *AiBase {
	ai := new(AiBase)
	ai.UUID = UUID
	ai.Name = Name
	ai.Type = Type
	ai.Filepath = Filepath
	ai.Description = Description
	return ai
}
