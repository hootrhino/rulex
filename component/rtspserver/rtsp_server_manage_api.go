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

package rtspserver

func RegisterFlvStreamSource(liveId string) error {
	return __DefaultRtspServer.RegisterFlvStreamSource(liveId)
}
func GetFlvStreamSource(liveId string) (*FlvStream, error) {
	return __DefaultRtspServer.GetFlvStreamSource(liveId)
}
func Exists(liveId string) bool {
	return __DefaultRtspServer.Exists(liveId)
}
func FlvStreamSourceList() []FlvStream {
	return __DefaultRtspServer.FlvStreamSourceList()
}
