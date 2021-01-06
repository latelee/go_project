#!/bin/sh

# 版本和编译时间 TODO：找一个好的方法：
Version="v1.0"
BuildTime=`date +'%Y-%m-%d %H:%M:%S'`

 
cd ../
GO111MODULE=on go build -mod vendor -ldflags "-X 'webdemo/app/cmd.BuildTime=${BuildTime}' -X 'webdemo/app/cmd.Version=${Version}'" -o app/webdemo app/main.go
cd -