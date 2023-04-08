#!/bin/bash

VERFILE="VERSION"
Version=`cat $VERFILE`

VER=""

if [ $# = 1 ] ; then
    VER=$1
    Version=""$VER
else
    echo "not set version, using " $Version
    echo "pls confirm(press enter)"
    read
fi

echo "build version" $Version

target=webdemo
suffix=.exe


# 版本和编译时间 TODO：找一个好的方法：

BuildDate=`date +'%Y-%m-%d '`
BuildTime1=`date +'%H:%M:%S'`
BuildTime=$BuildDate$BuildTime1

time GO111MODULE=on go build -mod vendor -ldflags "-X '$target/cmd.BuildTime=${BuildTime}' -X '$target/cmd.Version=${Version}'" -o $target$suffix main.go


