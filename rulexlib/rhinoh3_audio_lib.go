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

package rulexlib

import (
	"fmt"
	"os/exec"
	"time"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 调用系统自带的MPV播放器
*
 */
var __mkv_playing bool = false

func __MPVPlay(filePath string, duration time.Duration) error {
	if __mkv_playing {
		return fmt.Errorf("Audio output device busying now")
	}
	__mkv_playing = true
	cmd := exec.Command("mpv", filePath)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Error starting MP3 playback: %v", err)
	}
	go func(cmd *exec.Cmd) {
		timer := time.NewTimer(duration)
		<-timer.C
		cmd.Process.Kill()
	}(cmd)
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("Error starting MP3 playback: %v", err)
	}
	__mkv_playing = false
	return nil
}

/*
*
* applib:PlayMusic('001.mp3')
*
 */
func PlayMusic(rx typex.RuleX) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		filename := l.ToString(2)
		if filename != "" {
			err := __MPVPlay(filename, 10*time.Second)
			if err != nil {
				l.Push(lua.LString(err.Error()))
			}
		} else {
			l.Push(lua.LString("Invalid filename"))
		}
		l.Push(lua.LNil)
		return 1
	}
}
