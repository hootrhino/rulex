import socket
import threading
import tkinter as tk
import json


def update_values(
    root,
    ph_value,
    do_temp_value,
    do_value,
    salinity_conductivity_value,
    salinity_resistance_value,
    salinity_temp_value,
    data,
):
    try:
        if data[0] == "H":
            print("HeartBeat:", data)
            return
        sensor_data = json.loads(data)
        params = sensor_data.get("params", {})
        ph_value.set(params.get("ph", ""))
        do_temp_value.set(params.get("t1", ""))
        do_value.set(params.get("oxygen", ""))
        salinity_conductivity_value.set(params.get("conductivity", ""))
        salinity_resistance_value.set(params.get("resistance", ""))
        salinity_temp_value.set(params.get("t2", ""))
    except json.JSONDecodeError as e:
        print("Error parsing JSON data:", e, data)
    except KeyError as e:
        print("KeyError:", e)


def listen_for_data(
    root,
    ph_value,
    do_temp_value,
    do_value,
    salinity_conductivity_value,
    salinity_resistance_value,
    salinity_temp_value,
):
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

        # 开始处理客户端连接
        while True:
            client_socket, client_address = server_socket.accept()
            print("Connected to client:", client_address)

            # 持续读取数据并解析显示
            while True:
                data = client_socket.recv(1024).decode()
                if not data:
                    break
                root.after(
                    0,
                    update_values,
                    root,
                    ph_value,
                    do_temp_value,
                    do_value,
                    salinity_conductivity_value,
                    salinity_resistance_value,
                    salinity_temp_value,
                    data,
                )

    except Exception as e:
        print("Error:", e)

    finally:
        # 关闭服务器套接字
        server_socket.close()


def main():
    # 创建主窗口
    root = tk.Tk()
    root.title("海水数据")

    # 创建并设置布局
    main_frame = tk.Frame(root, padx=20, pady=20)
    main_frame.pack()

    # 创建指标标签和值显示Entry
    ph_value = tk.DoubleVar()
    do_temp_value = tk.DoubleVar()
    do_value = tk.DoubleVar()
    salinity_conductivity_value = tk.DoubleVar()
    salinity_resistance_value = tk.DoubleVar()
    salinity_temp_value = tk.DoubleVar()

    ph_label = tk.Label(main_frame, text="PH传感器温度:")
    ph_label.grid(row=0, column=0, sticky="w")
    ph_entry = tk.Entry(main_frame, textvariable=ph_value, state="readonly")
    ph_entry.grid(row=0, column=1)

    do_temp_label = tk.Label(main_frame, text="溶解氧传感器温度:")
    do_temp_label.grid(row=1, column=0, sticky="w")
    do_temp_entry = tk.Entry(main_frame, textvariable=do_temp_value, state="readonly")
    do_temp_entry.grid(row=1, column=1)

    do_label = tk.Label(main_frame, text="溶解氧传感器度:")
    do_label.grid(row=2, column=0, sticky="w")
    do_entry = tk.Entry(main_frame, textvariable=do_value, state="readonly")
    do_entry.grid(row=2, column=1)

    salinity_conductivity_label = tk.Label(main_frame, text="盐度传感器电导率:")
    salinity_conductivity_label.grid(row=3, column=0, sticky="w")
    salinity_conductivity_entry = tk.Entry(
        main_frame, textvariable=salinity_conductivity_value, state="readonly"
    )
    salinity_conductivity_entry.grid(row=3, column=1)

    salinity_resistance_label = tk.Label(main_frame, text="盐度传感器电阻率:")
    salinity_resistance_label.grid(row=4, column=0, sticky="w")
    salinity_resistance_entry = tk.Entry(
        main_frame, textvariable=salinity_resistance_value, state="readonly"
    )
    salinity_resistance_entry.grid(row=4, column=1)

    salinity_temp_label = tk.Label(main_frame, text="盐度传感器温度:")
    salinity_temp_label.grid(row=5, column=0, sticky="w")
    salinity_temp_entry = tk.Entry(
        main_frame, textvariable=salinity_temp_value, state="readonly"
    )
    salinity_temp_entry.grid(row=5, column=1)

    # 创建线程来监听数据
    data_thread = threading.Thread(
        target=listen_for_data,
        args=(
            root,
            ph_value,
            do_temp_value,
            do_value,
            salinity_conductivity_value,
            salinity_resistance_value,
            salinity_temp_value,
        ),
    )
    data_thread.daemon = True
    data_thread.start()

    # 启动事件循环
    root.mainloop()


if __name__ == "__main__":
    main()
