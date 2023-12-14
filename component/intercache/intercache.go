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

package intercache

/*
*
* 内部缓存器
*
 */
type InterCache interface {
	RegisterSlot(Slot string)   // 存储槽位, 释放资源的时候调用
	UnRegisterSlot(Slot string) // 注销存储槽位, 释放资源的时候调用
	Size() uint64               // 存储器当前长度
	Flush()                     // 释放存储器空间
}
