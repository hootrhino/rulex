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

package shellengine

import "context"

/*
*
* 脚本执行引擎
*
 */
type ShellEngine interface {
	JustRun(ctx context.Context, cmd string) error                                 // 一次性指令,不等返回
	JustRunDetach(ctx context.Context, cmd string) error                           // 分离进程
	RunAndWaitResult(ctx context.Context, cmd string) ([]byte, error)              // 运行指令然后等待结果
	RunWithInteractive(ctx context.Context, cmd string, in, out chan []byte) error // 交互运行，自带stdin
	ReadFD(fd int, p []byte) (n int, err error)                                    // Syscall read
	IOctl(fd int, request, arg uintptr) error                                      // UNIX ioctl
}
