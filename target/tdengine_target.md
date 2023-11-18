# Tdengine 支持

## 测试
新建一个数据库，这里用Docker:
```sh
docker run -d --name tdengine -p 6041:6041 tdengine/tdengine
```
新建测试库:
```sql
create database if not exists rhino
```
新建测试表
```sql
create table if not exists tb1 (ts timestamp, a int)
```
插入数据
```sql
use rhino;
insert into tb1 values(now, 0)(now+1s,1)(now+2s,2)(now+3s,3);
insert into tb1 values(now, 0)(now+1s,1)(now+2s,2)(now+3s,3);
insert into tb1 values(now, 0)(now+1s,1)(now+2s,2)(now+3s,3);
```

