package dto

import (
	"mime/multipart"
)

type CronTaskCreateDTO struct {
	Name     string                `form:"name" binding:"required"`
	CronExpr string                `form:"cronExpr" binding:"required"`
	TaskType int                   `form:"taskType" binding:"required"` // 1-shell 2-cmd
	Args     string                `form:"args"`                        // "param1 param2 param3"
	IsRoot   string                `form:"isRoot"`                      // 0-false 1-true
	Env      []string              `form:"env"`                         // ["A=e1", "B=e2", "C=e3"]
	File     *multipart.FileHeader `form:"file" binding:"required"`
}

type CronTaskUpdateDTO struct {
	ID       int                   `form:"id" binding:"required"`
	Name     string                `form:"name"`
	CronExpr string                `form:"cronExpr"`
	TaskType int                   `form:"taskType"` // 1-shell 2-cmd
	Args     *string               `form:"args"`     // "param1 param2 param3"
	IsRoot   string                `form:"isRoot"`   // 0-false 1-true
	Env      []string              `form:"env"`      // ["A=e1", "B=e2", "C=e3"]
	File     *multipart.FileHeader `form:"file"`
}
