# filestore-server

go语言开发的仿百度网盘的分布式云存储系统

开发环境说明：
- go语言版本：go1.19 windows/amd64
- 编译环境：Linux
- 开发工具：goland

## 项目初始化

### web框架选型
```bash
    go get -u github.com/gin-gonic/gin@v1.8.1
```

### 配置参数分离
```bash
    go get github.com/spf13/viper@v1.13.0
```
> 参考文档：https://github.com/spf13/viper

## 项目接口开发

### 1. 文件上传服务
- [x] 文件上传 | 下载 | 查询 | 删除

### 2. 用户管理
- [x] 用户注册 | 登录 | 查询