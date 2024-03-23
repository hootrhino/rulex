//go:build amd64
// +build amd64

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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package utils

import (
	"fmt"
	"image"
	"image/color"
	"sync"
	"time"

	"gocv.io/x/gocv"
)

/*
*
* OpenCV处理数据
*
 */
var __CGoMutex sync.Mutex

/*
*
* Jpeg Stream 帧
*
 */
type Resolution struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (O Resolution) String() string {
	return fmt.Sprintf("%dx%d", O.Width, O.Height)
}

func CvMatToImageBytes(FrameBuffer []byte) ([]byte, Resolution, error) {
	__CGoMutex.Lock()
	defer __CGoMutex.Unlock()
	imgMat := gocv.NewMat()
	defer imgMat.Close()
	err0 := gocv.IMDecodeIntoMat(FrameBuffer, gocv.IMReadFlag(gocv.IMReadColor), &imgMat)
	Resolution := Resolution{
		imgMat.Cols(), imgMat.Rows(),
	}
	if err0 != nil {
		return nil, Resolution, err0
	}
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	gocv.PutText(&imgMat, fmt.Sprintf("%s(%d*%d)", formattedTime, Resolution.Width, Resolution.Height), image.Point{5, 25},
		gocv.FontHersheyPlain, 2, color.RGBA{255, 0, 0, 0}, 2)
	NewImgMat := gocv.NewMat()
	defer NewImgMat.Close()
	if imgMat.Cols() > 1920 {
		ImgBytes, err1 := gocv.IMEncode(".jpg", imgMat)
		if err1 != nil {
			return nil, Resolution, err0
		}
		return ImgBytes.GetBytes(), Resolution, nil
	}
	gocv.Resize(imgMat, &NewImgMat, image.Point{}, 2, 2, gocv.InterpolationArea)
	Resolution.Width = imgMat.Cols() * 2
	Resolution.Height = imgMat.Rows() * 2
	ImgBytes, err1 := gocv.IMEncode(".jpg", NewImgMat)
	if err1 != nil {
		return nil, Resolution, err0
	}
	return ImgBytes.GetBytes(), Resolution, nil
}

/*
*
* 使用DNN来做AI处理, 返回值根据模型不同而不同，需要结合AIBase里面的规则
*
 */
func DNNForward(blob gocv.Mat, DnnNet gocv.Net) gocv.Mat {
	DnnNet.SetInput(blob, "")
	// outs.Size() 返回矩阵的形状
	// [a11 a12 a13 ...a1n]
	// [a21 a22 a23 ...a2n]
	// [a31 a32 a33 ...a3n]
	// z    := outs.Size()[0]
	// cols := outs.Size()[1]
	// rows := outs.Size()[2]
	return DnnNet.Forward("")
}
