## 文档
支持 GRPC 用户自定义协议接入。
### 协议定义
```proto
syntax = "proto3";

option go_package = "./;xstream";
option java_multiple_files = false;
option java_package = "xstream";
option java_outer_classname = "RulexXStream";

package xstream;

service XStream {
  rpc OnStreamApproached (stream Data) returns (Response) {}
}

message Data {
  string value = 1;
}

message Response {
  int32 code = 1;
  string message = 2;
}

```