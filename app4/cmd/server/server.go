package server

import (
    _ "fmt"
    "com"
    "k8s.io/klog"

    "github.com/kubeedge/beehive/pkg/core"
)

type server struct {
    enable bool
    // 后可加其它字段
}

func newServer(enable bool) *server {
    return &server{
        enable: enable,
    }
}

func Register() {
    core.Register(newServer(true))
}

func (a *server) Name() string {
    return "server"
}

func (a *server) Group() string {
    return "server"
}

// Enable indicates whether enable this module
func (a *server) Enable() bool {
    return a.enable
}

func (a *server) Start() {
    klog.Infoln("server...")
    go doit()
}

func doit() {
    for {
        klog.Infoln(".")
        com.Sleep(10000)
    }
}
