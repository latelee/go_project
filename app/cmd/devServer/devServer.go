package devServer

import (
    _ "fmt"
    //"github.com/latelee/go_project/pkg/com"
    "k8s.io/klog"

    "github.com/kubeedge/beehive/pkg/core"

    "github.com/latelee/go_project/app/conf"
    "github.com/latelee/go_project/app/pkg/vendors/NanWang"

)

type devServer struct {
    enable bool
    // 后可加其它字段
}

func init() {
    //core.Register(newdevServer(true))
}

func newdevServer(enable bool) *devServer {
    return &devServer{
        enable: enable,
    }
}

func Register(opts *conf.DevServer) {
    initConfig(opts)
    core.Register(newdevServer(opts.Enable))
}

func (a *devServer) Name() string {
    return "devServer"
}

func (a *devServer) Group() string {
    return "devServer"
}

// Enable indicates whether enable this module
func (a *devServer) Enable() bool {
    return a.enable
}

func (a *devServer) Start() {
    klog.Infoln("devServer start...")
    klog.Println(Config.Port, Config.Name, Config.Protocol)
    nanwang.Parser()
}
