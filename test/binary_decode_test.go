package test

import (
	"bytes"
	"encoding/binary"
	"testing"
)

type message struct {
	Flag    uint8
	Version uint8
	Type    uint64
}

// 定义编码函数
func EncodeMessage(msg message) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, msg)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 定义解码函数
func DecodeMessage(data []byte) (message, error) {
	var msg message
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.BigEndian, &msg)
	if err != nil {
		return message{}, err
	}
	return msg, nil
}

// go test -timeout 30s -run ^TestOk github.com/hootrhino/rulex/test -v -count=1
func TestEncodeMessage(t *testing.T) {
	// 创建一个消息
	msg := message{
		Flag:    1,
		Version: 2,
		Type:    3,
	}

	// 编码消息
	encoded, err := EncodeMessage(msg)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(encoded)

	// 解码消息
	decoded, err := DecodeMessage(encoded)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(decoded)
}
