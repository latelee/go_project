package update

import (
    _ "fmt"
    "com"
    "k8s.io/klog"

    "github.com/kubeedge/beehive/pkg/core"
)

type update struct {
    enable bool
    // 后可加其它字段
}

func newUpdate(enable bool) *update {
    return &update{
        enable: enable,
    }
}

func Register() {
    core.Register(newUpdate(true))
}

func (a *update) Name() string {
    return "update"
}

func (a *update) Group() string {
    return "update"
}

// Enable indicates whether enable this module
func (a *update) Enable() bool {
    return a.enable
}

func (a *update) Start() {
    klog.Infoln("update...")
    //go doit()
}

func doit() {
    for {
        klog.Infoln(".")
        com.Sleep(1000)
    }
}
