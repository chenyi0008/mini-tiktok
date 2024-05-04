# Mini-tiktok

> 迷你抖音，基于go-zero gorm

用到的命令
```text
# 创建API服务
goctl api new core
# 启动服务
go run core.go -f etc/core-api.yaml
# 使用api文件生成代码
goctl api go -api core.api -dir . -style go_zero
# proto生成代码
goctl rpc protoc xx.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=.
```

腾讯云COS后台地址：https://console.cloud.tencent.com/cos/bucket

腾讯云COS帮助文档：https://cloud.tencent.com/document/product/436/31215




