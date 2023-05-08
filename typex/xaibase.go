package typex

type AIType string

// 内建类型
const BODY_POSE_RECOGNITION AIType = "BODY_POSE_RECOGNITION"

/*
*
* AI 应用管理器接口
*
 */
type XAiRuntime interface {
	GetRuleX() RuleX
	ListAi() []*AI
	LoadAi(Ai *AI) error
	GetAi(uuid string) *AI
	RemoveAi(uuid string) error
	UpdateAi(Ai *AI) error
	StartAi(uuid string) error
	StopAi(uuid string) error
	Stop()
}

/*
*
* AI 应用层接口
*
 */
type XAi interface {
	Start(map[string]interface{}) error
	Infer([][]float64) [][]float64
	Stop()
}

/*
*
* 内建AI
*
 */
type AI struct {
	UUID        string                 `json:"uuid"`        // UUID
	Name        string                 `json:"name"`        // 名称
	Type        AIType                 `json:"type"`        // 类型
	IsBuildIn   bool                   `json:"isBuildIn"`   // 是否内建
	Filepath    string                 `json:"filepath"`    // 文件路径, 是相对于main的aispace目录
	Config      map[string]interface{} `json:"config"`      // 内部配置
	Description string                 `json:"description"` // 描述文字
	XAI         XAi                    `json:"-"`
}

/*
*
* 生成typex.AI应用
*
 */
func NewAI(UUID, Name, Type, Filepath, Description string) *AI {
	ai := new(AI)
	ai.UUID = UUID
	ai.Name = Name
	ai.Type = AIType(Type)
	ai.Filepath = Filepath
	ai.Description = Description
	return ai
}
