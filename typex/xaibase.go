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
	Infer([][]float64) [][]float64
	Stop()
}
type AiBase struct {
	UUID        string                 // UUID
	Name        string                 // 名称
	Type        string                 // 类型
	Filepath    string                 // 文件路径, 是相对于main的aispace目录
	Config      map[string]interface{} // 内部配置
	AiBase      XAiBase                // AI工作模型
	Description string                 // 描述文字
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
