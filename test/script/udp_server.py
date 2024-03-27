# Copyright (C) 2024 wwhai
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.
import socket


def main():
    # 服务器地址和端口
    server_address = ("0.0.0.0", 5557)

    # 创建UDP套接字
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

    try:
        # 绑定地址和端口
        server_socket.bind(server_address)
        print("Server is listening on port 5557...")

        # 持续读取数据并打印
        while True:
            data, client_address = server_socket.recvfrom(1024)
            print("Received:", data.decode(), "from", client_address)

    except Exception as e:
        print("Error:", e)

    finally:
        # 关闭服务器套接字
        server_socket.close()


if __name__ == "__main__":
    main()
