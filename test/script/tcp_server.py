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
    server_address = ("0.0.0.0", 5556)

    # 创建TCP套接字
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    try:
        # 绑定地址和端口
        server_socket.bind(server_address)

        # 开始监听
        server_socket.listen(5)
        print("Server is listening on port 5556...")

        # 接受客户端连接
        client_socket, client_address = server_socket.accept()
        print("Connected to client:", client_address)

        # 持续读取数据并打印
        while True:
            data = client_socket.recv(1024)
            if not data:
                break
            print("Received:", data.decode())

    except Exception as e:
        print("Error:", e)

    finally:
        # 关闭客户端套接字
        client_socket.close()

        # 关闭服务器套接字
        server_socket.close()


if __name__ == "__main__":
    main()
