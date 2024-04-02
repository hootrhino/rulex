// Copyright (C) 2023 wwhai
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

//	if Duration, err := pingQ(ip, 2000*time.Millisecond); err != nil {
//		glogger.GLogger.WithFields(Fields).Info(fmt.Sprintf(
//			"[Count:%d] Ping Error:%s", i,
//			err.Error()))
//	} else {
//
//		glogger.GLogger.WithFields(Fields).Info(fmt.Sprintf(
//			"[Count:%d] Ping Reply From %s: time=%v ms TTL=128", i,
//			tt, Duration))
//	}
package rulexlib

import (
	"net"
	"time"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* Ping
*
 */
func PingIp(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		ip := l.ToString(2)
		Duration, err := pingQ(ip, 5000*time.Millisecond)
		if err != nil {
			l.Push(lua.LNumber(0))
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNumber(Duration))
			l.Push(lua.LNil)
		}
		return 2
	}
}

// --------------------------------------------------------------------------------------------------
// private
// --------------------------------------------------------------------------------------------------
func pingQ(ip string, timeout time.Duration) (time.Duration, error) {
	const IcmpLen = 8
	msg := [32]byte{
		8, 0, 0, 0, 0, 13, 0, 37,
	}
	check := checkSum(msg[:IcmpLen])
	msg[2] = byte(check >> 8)
	msg[3] = byte(check & 255)

	remoteAddr, err := net.ResolveIPAddr("ip", ip)
	if err != nil {
		return 0, err
	}
	conn, err := net.DialIP("ip:icmp", nil, remoteAddr)
	if err != nil {
		return 0, err
	}
	start := time.Now()
	if _, err := conn.Write(msg[:IcmpLen]); err != nil {
		return 0, err
	}
	conn.SetReadDeadline(time.Now().Add(timeout))
	_, err1 := conn.Read(msg[:])
	conn.SetReadDeadline(time.Time{})
	if err1 != nil {
		return 0, err1
	}
	return time.Since(start), nil
}

func checkSum(msg []byte) uint16 {
	sum := 0
	for n := 0; n < len(msg); n += 2 {
		sum += int(msg[n])<<8 + int(msg[n+1])
	}
	sum = (sum >> 16) + sum&0xffff
	sum += sum >> 16
	return uint16(^sum)
}
