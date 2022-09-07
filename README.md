https://github.com:latelee/go_project.git

## 概述
web服务。带测试用的代码（单独文件）

### 设计思路

命令行使用 cobra 框架。  
系统参数使用 viper，可读取，可保存。  
web 服务使用 gin 框架。    
前端使用 bootstrap，但和golang紧密结合，利用template设计页面。  
用于测试的各种功能见下

## 功能

- [x] 默认整合post接口及内部网站
- [x] post各种测试示例 如发送的json在文件中，直接发送post
- [x] post客户端
- [ ] 转发测试
- [ ] https
- [ ] 
- [ ] WS
- [ ] 页面（优先级低）

其它功能：

- [x] 兼容不同平台的编译（即编译相应平台的go源码文件）
- [ ] golang 调用C++（实为C）
- [ ] 文件打包成golang代码
- [ ] 参数保存

## 目录及文件说明
```
.
|-- README.md
|-- cmd  上层业务逻辑模块
|   |-- gin  gin 服务
|   |-- rootCmd.go 除读取参数外，一般不用动
|   `-- tcpp  tcp服务
|-- common 全局变量 结构体声明
|   |-- conf
|   `-- constants
|-- config.yaml 配置文件
|-- data  数据文件
|   |-- dist 静态网站（适用于仅提供website功能）
|   `-- static 前端网页静态资源
|-- go.mod
|-- go.sum
|-- main.go 一般不用动
|-- mybuild.sh 编译脚本
|-- note.md  一般调试笔记
|-- pkg 自实现的库
|   |-- com 共用库，底层封装
|   |-- db  数据库连接模块（稍高一层的接口封装，下同）
|   |-- httpc
|   |-- klog
|   |-- update
|   |-- utils
|   |-- vendors
|   |-- wait
|   `-- ws
`-- vendor // 第三方依赖库
```

`cmd/gin`为web服务实现的主要目录：
router_xxx.go： 和具体的响应url路由相关的实现


## 编译相关

编译命令：

```
./mybuild.sh
```

运行：
```
./webdemo
```

## 移植相关
工程名称默认为webdemo，如要移植到新工程，将工程所有文件的 webdemo 改名。包括：
go.mod、编译脚本、代码文件

## 问题及解决

