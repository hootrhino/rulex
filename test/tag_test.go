package test

import (
	"encoding/binary"
	"encoding/json"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/typex"
)

func Test_RenderConfig(t *testing.T) {
	//
	type __config struct {
		A    int32   `json:"a" validate:"required" title:"title--a" info:"aaaa" hidden:"true"`
		B    int64   `json:"b" validate:"required" title:"title--b" info:"bbbb" placeholder:"BBBBBBBBBBBB"`
		C    string  `json:"c" validate:"required" title:"title--c" info:"cccc" options:"cv1,cv2|cv3,cv4"`
		D    float32 `json:"d" validate:"required" title:"title--d" info:"dddd"`
		F    []int   `json:"f" validate:"required" title:"title--f" info:"ffff"`
		File string  `json:"file" validate:"required" title:"File" info:"File" file:"uploadfile"`
	}

	xcfgs, err := core.RenderInConfig(typex.MQTT, "MQTT", __config{})
	if err != nil {
		t.Fatal(err)
	} else {
		b, _ := json.MarshalIndent(xcfgs, "", " ")
		t.Log(string(b))
	}
}
func Test_binary_to_int(t *testing.T) {
	data := []byte{
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x04,
		0x04, 0x03, 0x02, 0x01,
	}
	Address := binary.BigEndian.Uint32(data[0:4])
	Start := binary.BigEndian.Uint32(data[4:8])
	Size := binary.BigEndian.Uint32(data[8:12])
	assert.Equal(t, uint32(1), (Address))
	assert.Equal(t, uint32(1), (Start))
	assert.Equal(t, uint32(4), (Size))
}
