package test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/hootrhino/rulex/utils"
)

// go test -timeout 30s -run ^TestOk github.com/hootrhino/rulex/test -v -count=1
func Test_CheckSumCRC16(t *testing.T) {
	m_data := []byte{0x01, 0x02, 0x03, 0x04}
	checksum := utils.CRC16(m_data)
	fmt.Printf("check sum:%X \n", checksum)
	int16buf := new(bytes.Buffer)
	binary.Write(int16buf, binary.LittleEndian, checksum)
	fmt.Printf("write buf is: %+X \n", int16buf.Bytes())
	fmt.Printf("output-before:%X \n", m_data)
	m_data = append(m_data, int16buf.Bytes()...)
	fmt.Printf("output-after:%X \n", m_data)
}
func Test_CheckXOR(t *testing.T) {
	data := []byte{0xFF, 0xFF, 0xFF, 0x50, 0x03, 0xFA, 0x01, 0x01} //56
	t.Log(utils.XOR(data) == 0x56)
}
