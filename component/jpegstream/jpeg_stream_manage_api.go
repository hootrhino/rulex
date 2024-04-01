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
package jpegstream



func RegisterJpegStreamSource(liveId string) error {
	return __DefaultJpegStreamServer.RegisterJpegStreamSource(liveId)
}
func GetJpegStreamSource(liveId string) (*JpegStream, error) {
	return __DefaultJpegStreamServer.GetJpegStreamSource(liveId)
}
func Exists(liveId string) bool {
	return __DefaultJpegStreamServer.Exists(liveId)
}
func JpegStreamSourceList() []JpegStream {
	return __DefaultJpegStreamServer.JpegStreamSourceList()
}
