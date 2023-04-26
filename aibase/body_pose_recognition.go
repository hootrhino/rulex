package aibase

import (
	"github.com/hootrhino/rulex/typex"
)

type BodyPoseRecognition struct {
}

func NewBodyPoseRecognition(re typex.RuleX) typex.XAi {
	return &BodyPoseRecognition{}
}
func (ba *BodyPoseRecognition) Start(map[string]interface{}) error {

	return nil
}
func (ba *BodyPoseRecognition) Infer(input [][]float64) [][]float64 {
	return [][]float64{
		{110000, 120000, 130000},
		{210000, 220000, 230000},
		{310000, 320000, 330000},
	}
}
func (ba *BodyPoseRecognition) Stop() {

}
