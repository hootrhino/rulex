package dto

type CronTaskCreateDTO struct {
	Name     string   `form:"name" binding:"required" json:"name"`
	CronExpr string   `form:"cronExpr" binding:"required" json:"cronExpr"`
	TaskType int      `form:"taskType" binding:"required" json:"taskType"` // 1-shell 2-cmd
	Args     *string  `form:"args" json:"args"`                            // "param1 param2 param3"
	IsRoot   string   `form:"isRoot" json:"isRoot"`                        // 0-false 1-true
	Env      []string `form:"env" json:"env"`                              // ["A=e1", "B=e2", "C=e3"]
	Script   string   `form:"script" json:"script"`
}

type CronTaskUpdateDTO struct {
	UUID     string   `form:"uuid" binding:"required" json:"uuid"`
	Name     string   `form:"name" json:"name"`
	CronExpr string   `form:"cronExpr" json:"cronExpr"`
	TaskType int      `form:"taskType" json:"taskType"` // 1-shell 2-cmd
	Args     *string  `form:"args" json:"args"`         // "param1 param2 param3"
	IsRoot   string   `form:"isRoot" json:"isRoot"`     // 0-false 1-true
	Env      []string `form:"env" json:"env"`           // ["A=e1", "B=e2", "C=e3"]
	Script   string   `form:"script" json:"script"`
}
