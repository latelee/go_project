#!/bin/sh

# 版本和编译时间 TODO：找一个好的方法：
Version="v1.0"
BuildTime=`date +'%Y-%m-%d %H:%M:%S'`


#-ldflags "-X 'github.com/latelee/dbtool/cmd.BuildTime=${BuildTime}' -X 'github.com/latelee/dbtool/cmd.Version=${Version}'"

 
cd ../
GO111MODULE=on go build -mod vendor -o app/webdemo app/main.go
cd -