// Copyright (C) 2024 wwhai
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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package aibase

import "fmt"

type Yolo8ObjectDetectionCpu struct {
	Path string
	Mode string
}

func NewYolo8ObjectDetectionCpu() XAlgorithm {
	return &Yolo8ObjectDetectionCpu{
		Path: "yolov8n.onnx",
		Mode: "GPU",
	}
}
func (Yolo8 *Yolo8ObjectDetectionCpu) Init(Config map[string]interface{}) error {
	return fmt.Errorf("NOT support Arm32")
}
func (Yolo8 *Yolo8ObjectDetectionCpu) Load() error {
	return nil
}

/*
*
* 推断
*
 */
func (Yolo8 *Yolo8ObjectDetectionCpu) Forward(Input []byte) (map[string]interface{}, error) {

	return map[string]interface{}{}, nil
}
func (Yolo8 *Yolo8ObjectDetectionCpu) Unload() error {
	return nil
}
func (Yolo8 *Yolo8ObjectDetectionCpu) AlgorithmDetail() Algorithm {
	return Algorithm{
		UUID:        "Yolo8ObjectDetectionCpu",
		Type:        "BUILDIN", // 内置
		Name:        "Yolo8 Object Detection Cpu",
		State:       1,
		Document:    "https://docs.ultralytics.com",
		Description: "OpenCv DNN Module: Yolo8 Object Detection Cpu Version",
	}
}
