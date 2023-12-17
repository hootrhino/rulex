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

package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hootrhino/rulex/component/eventbus"
)

// @ go test -timeout 30s -run ^TestEventBus github.com/hootrhino/rulex/test -v -count=1
func TestEventBus(t *testing.T) {

	eventbus.InitEventBus()
	eventbus.Subscribe("hello", &eventbus.Subscriber{
		Callback: func(Topic string, Msg eventbus.EventMessage) {
			t.Log("hello:", Msg)
		},
	})
	start := time.Now()
	for i := 0; i < 100; i++ {
		eventbus.Publish("hello", eventbus.EventMessage{
			Payload: fmt.Sprintf("world:%d", i),
		})
	}
	duration := time.Since(start)
	t.Log("time.Since(start):", duration)
	time.Sleep(3 * time.Second)
	eventbus.Flush()
}
