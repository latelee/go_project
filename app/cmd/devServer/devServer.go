package devServer

import (
    _ "fmt"
    //"com"
    "k8s.io/klog"

    "github.com/kubeedge/beehive/pkg/core"

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

func Register() {
    core.Register(newdevServer(true))
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
}
