package test

import (

	"fmt"
	"io"
	"os"
	"testing"

)

func Test_read_line(t *testing.T) {
	t.Log(readLastLine("conf/rulex.ini"))
}

func readLastLine(filepath string) string {
	fileHandle, err := os.Open(filepath)
	if err != nil {
		return ""
	}
	defer fileHandle.Close()
	line := ""
	var cursor int64 = 0
	stat, _ := fileHandle.Stat()
	filesize := stat.Size()
	for {
		cursor -= 1
		fileHandle.Seek(cursor, io.SeekEnd)
		char := make([]byte, 1)
		fileHandle.Read(char)
		if cursor != -1 && (char[0] == 10 || char[0] == 13) {
			break
		}
		line = fmt.Sprintf("%s%s", string(char), line)
		if cursor == -filesize {
			break
		}
	}
	return line
}
