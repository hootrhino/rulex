package dto

import (
	"mime/multipart"
)

type CronTaskCreateDTO struct {
	ID       int
	Name     string                `form:"name" binding:"required"`
	CronExpr string                `form:"cronExpr" binding:"required"` // Quartz standard
	TaskType int                   `form:"taskType" binding:"required"` // 1-shell 2-cmd
	Args     string                `form:"args"`                        // "param1 param2 param3"
	IsRoot   string                `form:"isRoot"`                      // 0-false 1-true
	Env      []string              `form:"env"`                         // ["A=e1", "B=e2", "C=e3"]
	File     *multipart.FileHeader `form:"file"`
}
