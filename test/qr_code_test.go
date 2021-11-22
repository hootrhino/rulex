package test

import (
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

func Test_gen_QR_code(t *testing.T) {
	enc := qrcode.NewQRCodeWriter()
	img, _ := enc.Encode("Hello, World!", gozxing.BarcodeFormat_QR_CODE, 250, 250, nil)
	file, _ := os.Create("data/qrcode.png")
	_ = png.Encode(file, img)
	defer file.Close()
}
func Test_read_QR_code(t *testing.T) {
	file, _ := os.Open("data/qrcode.png")
	img, _, _ := image.Decode(file)
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
	qrReader := qrcode.NewQRCodeReader()
	result, _ := qrReader.Decode(bmp, nil)
	t.Log(result)
	defer file.Close()

}
