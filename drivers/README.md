# 注意
这个目录下的文件是专门针对树莓派开发的一些模块，在不同的硬件上需要调整规范，并不通用。
## 生成序列号
```c
#include <stdio.h>
#include <stdlib.h>
#include <time.h>
int main(int argc, char const *argv[])
{
    FILE *fp = NULL;
    fp = fopen("serial_number.csv", "w");
    if (fp == NULL)
    {
        printf("file can't be opened\n");
        exit(1);
        return 0;
    }
    srand((unsigned)time(NULL));
    fprintf(fp, "Address , Channel\n");
    for (size_t i = 1; i < 0xFF; i++)
    {
        fprintf(fp, "%d , %d\n", i, rand() % 0xFFFF + 10000);
    }
    fclose(fp);
    return 0;
}

```
## 创建一个Lora模块
Config:
```json
{
	"name": "LoraATKModule",
	"type": "LoraATK",
	"config": {
		"name": "COM2",
		"baud": 115200,
		"readTimeout": "test",
		"size": 8,
		"parity": "N",
		"stopbits": 1
	},
	"description": "Lora ATK Module"
}
```
Test:

```sh
go run ./rulexc.go inend-create --config  '{\"name\":\"LoraATKModule\",\"type\":\"LoraATK\",\"config\":{\"name\":\"COM1\",\"baud\":\"115200\",\"readTimeout\":\"0\",\"size\":\"8\",\"parity\":\"N\",\"stopbits\":\"1\"},\"description\":\"Lora ATK Module\"}'
```