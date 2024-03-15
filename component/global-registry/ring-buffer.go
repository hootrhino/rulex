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

package globalregistry

import (
	"errors"
)

// RingBuffer 是环形缓冲区结构体
type RingBuffer struct {
	buffer []byte
	size   int
	head   int
	tail   int
	count  int
}

// NewRingBuffer 创建一个新的环形缓冲区
func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		buffer: make([]byte, size),
		size:   size,
		head:   0,
		tail:   0,
		count:  0,
	}
}

// Write 向环形缓冲区写入数据
func (rb *RingBuffer) Write(data []byte) (int, error) {
	if len(data) > rb.size-rb.count {
		return 0, errors.New("not enough space in buffer")
	}

	var written int
	for _, b := range data {
		rb.buffer[rb.tail] = b
		rb.tail = (rb.tail + 1) % rb.size
		written++
		rb.count++
	}
	return written, nil
}

// Read 从环形缓冲区读取数据
func (rb *RingBuffer) Read() (byte, error) {
	if rb.count == 0 {
		return 0, errors.New("buffer is empty")
	}

	b := rb.buffer[rb.head]
	rb.head = (rb.head + 1) % rb.size
	rb.count--
	return b, nil
}

// func main() {
// 	rb := NewRingBuffer(5)

// 	// 写入数据
// 	data := []byte{'a', 'b', 'c', 'd', 'e'}
// 	written, err := rb.Write(data)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	fmt.Println("Written:", written)

// 	// 读取数据
// 	for i := 0; i < 5; i++ {
// 		b, err := rb.Read()
// 		if err != nil {
// 			fmt.Println("Error:", err)
// 			return
// 		}
// 		fmt.Printf("Read: %c\n", b)
// 	}
// }
