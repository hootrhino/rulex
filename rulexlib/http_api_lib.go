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
	"io"
	"net/http"
	"strings"
	"time"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* HTTP GET
*
 */
func HttpGet(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		url := l.ToString(2)
		l.Push(lua.LString(__HttpGet(url)))
		return 1
	}
}
func HttpPost(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		url := l.ToString(2)
		body := l.ToString(3)
		l.Push(lua.LString(__HttpPost(url, body)))
		return 1
	}
}

/*
*
* GET
*
 */
func __HttpGet(url string) string {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

/*
*
* POST
 */
func __HttpPost(url string, body string) string {
	resp, err := http.Post(url, "application/json", strings.NewReader(body))
	if err != nil {
		glogger.GLogger.Error(err)
		return ""
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		glogger.GLogger.Error(err)
		return ""
	}
	return string(responseBody)
}
