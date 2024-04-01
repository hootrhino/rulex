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
package jpegstream

import (
	"fmt"

	"github.com/hootrhino/rulex/utils"
)

type JpegStream struct {
	frame         []byte
	frameSize     int
	headerSize    int
	GetFirstFrame bool
	Type          string
	LiveId        string
	Pulled        bool
	Resolution    utils.Resolution
}

func (S JpegStream) String() string {
	return fmt.Sprintf(`{"liveId":%s,"pulled":%v,"resolution":%s}`,
		S.LiveId, S.Pulled, S.Resolution.String())
}
func (s *JpegStream) GetWebJpegFrame() []byte {
	b := s.frame[:s.frameSize]
	return b
}
func (s *JpegStream) GetRawFrame() []byte {
	return s.frame[s.headerSize:s.frameSize]
}
func (s *JpegStream) UpdateJPEG(jpeg []byte) {
	__headerReal := "--MJPEG_BOUNDARY\r\nContent-Type: image/jpeg\r\nContent-Length: %d\r\nX-Timestamp: 0.000000\r\n\r\n"
	size := len(jpeg)
	frameSize := size
	headerReal := fmt.Sprintf(__headerReal, size)
	if len(s.frame) < size+len(headerReal) {
		s.frame = make([]byte, (size+len(headerReal))*2)
	}
	copy(s.frame, headerReal)
	copy(s.frame[len(headerReal):], jpeg[:size])
	s.headerSize = len(headerReal) + 10 // \r\n出现了5次 计算为10个字符
	s.frameSize = frameSize

}
