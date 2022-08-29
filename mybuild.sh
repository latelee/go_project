#!/bin/bash

# 版本和编译时间 TODO：找一个好的方法：
Version="v1.0"
BuildTime=`date +'%Y-%m-%d %H:%M:%S'`

GO111MODULE=on go build -mod vendor -ldflags "-X 'webdemo/cmd.BuildTime=${BuildTime}' -X 'webdemo/cmd.Version=${Version}'" -o webdemo main.go
