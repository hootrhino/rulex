# 设备文档框架

每个设备、资源都需要有对应的文档，前期先暂时放在源码目录下，未来通过CI提取到文档项目中。下面是规范文档.

---

## 简介

TC-S200 系列空气质量监测仪内置 PM2.5、TVOC、甲醛、CO2，温湿度等高精度传感器套件，可通过吸顶式或壁挂安装，RS-485 接口通过 Modbus-RTU 协议进行数据输出，通过网关组网，或配合联动模块可以用于新风联动控制。

## 配置

```json
{
	"host":       "127.0.0.1",
	"port":       1883
}
```

## 数据

### 读
#### 1. CMD

| 值  | 含义 |
| --- | ---- |
| 0   | 开门 |
| 1   | 关门 |

#### 1. Args

| 值      | 含义 |
| ------- | ---- |
| "OPEN"  | 开门 |
| "CLOSE" | 关门 |

#### 3. 数据样例
```json
{
	"host":       "127.0.0.1",
	"port":       1883
}
```

### 写
#### 1. CMD

| 值  | 含义 |
| --- | ---- |
| 0   | 开门 |
| 1   | 关门 |

#### 2. Args

| 值      | 含义 |
| ------- | ---- |
| "OPEN"  | 开门 |
| "CLOSE" | 关门 |

#### 3. 数据样例
```json
{
	"host":       "127.0.0.1",
	"port":       1883
}
```

## 案例

```lua
-- Success
function Success()
    print("Success")
end
-- Failed
function Failed(error)
    print("Error:", error)
end

-- Actions
Actions = {
    function(data)
        local _, err = rulexlib:WriteDevice(device, rulexlib:T2J({ cmd = "open"}))
        if (err ~= nil) then
            log('WriteDevice open err: ', err)
            return false, data
        end
        return true, data
    end
}

```

## 说明

- XXX
- YYY
- or 暂无

## 社区

### 维护者

- User1
- Admin2
