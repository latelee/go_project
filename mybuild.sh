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

proname=webdemo
target=$proname-win.exe

BUILD_TYPE="win"

ARCH_NAME1=`uname -m`
# check for platform
if [ $ARCH_NAME1 = "x86_64" ]; then
SYS=`uname -s`
if [ $SYS = "Linux" ]; then
BUILD_TYPE="linux"
target=$proname
fi
elif [ $ARCH_NAME1 = "aarch64" ]; then
BUILD_TYPE="arm"
target=$proname-arm
elif [ $ARCH_NAME1 = "loongarch64" ]; then
BUILD_TYPE="mips"
target=$proname-mips
fi
echo "build for platform:" $BUILD_TYPE " version:" $Version " output file:" $target

# 版本和编译时间 TODO：找一个好的方法：

BuildDate=`date +'%Y-%m-%d '`
BuildTime1=`date +'%H:%M:%S'`
BuildTime=$BuildDate$BuildTime1

# -x for details
time GO111MODULE=on go build -mod vendor -ldflags "-X '$proname/cmd.BuildTime=${BuildTime}' -X '$proname/cmd.Version=${Version}'" -o $target main.go || exit 1

if [ $BUILD_TYPE != "win" ]; then
echo "need copy" $target "to bin"
mkdir -p bin
cp $target bin/$target.$BuildDate
fi
