package aibase

import (
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
)

type Mnist struct {
}

func NewMnist(re typex.RuleX) typex.XAi {
	return &Mnist{}
}
func (ba *Mnist) Start(map[string]interface{}) error {

	return nil
}
func (ba *Mnist) Infer(input [][]float64) [][]float64 {
	glogger.GLogger.Debug("Mnist.Infer:", input)
	return [][]float64{
		{110000, 120000, 130000},
		{210000, 220000, 230000},
		{310000, 320000, 330000},
	}
}
func (ba *Mnist) Stop() {

}
