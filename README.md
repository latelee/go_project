

## 背景
web服务集大者。带测试用的代码（单独文件）

### 设计思路

命令行使用 cobra 框架。  
系统参数使用 viper，可读取，可保存。  
web 服务使用 gin 框架。    
前端使用 bootstrap，但和golang紧密结合，利用template设计页面。  

## 功能

- [x] 默认整合post接口及内部网站
- [ ] post各种测试示例 如发送的json在文件中，直接发送post
- [ ] post客户端
- [ ] 转发测试
- [ ] https
- [ ] golang 调用C++（实为C）
- [ ] 文件打包成golang代码
- [ ] 参数保存
- [ ] WS
- [ ] 页面（优先级低）
- [ ] 

## 目录说明
data
static：静态资源、动态库  


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

