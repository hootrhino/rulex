// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
