#!/bin/sh

cd ../
GO111MODULE=on go build -mod vendor -o app/webdemo app/main.go
cd -