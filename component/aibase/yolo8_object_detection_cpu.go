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

import (
	"image"
	"math"

	"gocv.io/x/gocv"
)

type Yolo8ObjectDetectionCpu struct {
	Path   string
	Mode   string
	DnnNet gocv.Net
}

func NewYolo8ObjectDetectionCpu() XAlgorithm {
	return &Yolo8ObjectDetectionCpu{
		Path: "yolov8n.onnx",
		Mode: "GPU",
	}
}
func (Yolo8 *Yolo8ObjectDetectionCpu) Init(Config map[string]interface{}) error {
	// Yolo8.Path = Config["path"].(string)
	// Yolo8.Mode = Config["mode"].(string)
	Yolo8.Path = "./component/aibase/models/yolov8n.onnx"
	Yolo8.Mode = "GPU"
	return nil
}
func (Yolo8 *Yolo8ObjectDetectionCpu) Load() error {
	Yolo8.DnnNet = gocv.ReadNetFromONNX(Yolo8.Path)
	if Yolo8.Mode == "CPU" {
		Yolo8.DnnNet.SetPreferableBackend(gocv.NetBackendOpenCV)
		Yolo8.DnnNet.SetPreferableTarget(gocv.NetTargetCPU)
	}
	if Yolo8.Mode == "GPU" {
		Yolo8.DnnNet.SetPreferableBackend(gocv.NetBackendCUDA)
		Yolo8.DnnNet.SetPreferableTarget(gocv.NetTargetCUDA)
	}
	return nil
}

/*
*
* 推断
*
 */
func (Yolo8 *Yolo8ObjectDetectionCpu) Forward(Input []byte) (map[string]interface{}, error) {
	ImgMat, err := gocv.NewMatFromBytes(640, 640, gocv.MatTypeCV16S, Input)
	if err != nil {
		return nil, err
	}
	blob := gocv.BlobFromImage(ImgMat, 1/255.0, image.Point{640, 640}, gocv.Scalar{}, true, false)
	Yolo8.DnnNet.SetInput(blob, "")
	Result := GetYolo8ForwardResult(Yolo8.DnnNet.Forward(""))
	return map[string]interface{}{
		"Result": Result,
	}, nil
}
func (Yolo8 *Yolo8ObjectDetectionCpu) Unload() error {
	return Yolo8.DnnNet.Close()
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

/*
*
* 缩放图片大小，适配Yolo8模型输入
*
 */
func ResizeImgToFitYolo8(src gocv.Mat, dst *gocv.Mat, size image.Point) {
	// 计算缩放比例
	k := math.Min(float64(size.X)/float64(src.Cols()), float64(size.Y)/float64(src.Rows()))
	// 计算新尺寸
	newSize := image.Pt(int(k*float64(src.Cols())), int(k*float64(src.Rows())))
	// 调整图像尺寸
	gocv.Resize(src, dst, newSize, 0, 0, gocv.InterpolationLinear)

	// 如果目标尺寸不等于指定尺寸，则重新创建目标矩阵
	if dst.Cols() != size.X || dst.Rows() != size.Y {
		*dst = gocv.NewMatWithSize(size.Y, size.X, src.Type())
	}
	// 将调整后的图像复制到目标图像的对应区域
	rectOfDst := image.Rect(0, 0, newSize.X, newSize.Y)
	regionOfDst := (*dst).Region(rectOfDst)
	tmp := src.Region(rectOfDst)
	tmp.CopyTo(&regionOfDst)
}

type Result struct {
	Class             int
	X, Y, W, H, Score float32
}

func GetYolo8ForwardResult(Mat gocv.Mat) []Result {
	//
	return []Result{}
}
