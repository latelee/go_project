package modules

import (
	//"klog"
    "com"
    "k8s.io/klog"
	"github.com/kubeedge/beehive/pkg/core"
	c "github.com/kubeedge/beehive/pkg/core/context"
)

//Constants for module name and group
const (
	DestinationModule = "destinationmodule"
	DestinationGroup  = "destinationgroup"
)

type testModuleDest struct {
	//context *context.Context
    enable bool
}

func init() {
	core.Register(newtestModuleDest(true))
}

func newtestModuleDest(enable bool) *testModuleDest {
    return &testModuleDest{
        enable: enable,
    }
}

func (*testModuleDest) Name() string {
	return DestinationModule
}

func (*testModuleDest) Group() string {
	return DestinationGroup
}

func (a *testModuleDest) Enable() bool {
    return a.enable
}

func (m *testModuleDest) Start() {
	
    for {
        select {
        case <-c.Done():
            klog.Info("Stop test recv")
            return
        default:
        message, err := c.Receive(DestinationModule)
        //klog.Printf("destination module receive message:%v error:%v\n", message, err)
        if err != nil {
            continue
        }
        msg := message.GetContent()
        srcMdl := message.GetSource()
        if msg == "test1" {
            go func() {
                klog.Printf("got test1 from %v ++++++++++++", srcMdl)
                com.Sleep(1000)
                klog.Printf("got again111 %v %v", message.GetSource(), message.GetContent)
            }()
        } else if msg == "test2" {
            go func() {
                klog.Printf("got test2 from %v --------------", srcMdl)
                com.Sleep(1000)
                klog.Printf("got again222 %v %v", message.GetSource(), message.GetContent)
            }()
        } else {
            klog.Printf("got default %s %s...", msg, message.GetContent())
            resp := message.NewRespByMessage(&message, "fine")
            if message.IsSync() {
                c.SendResp(*resp)
            }
        }


    }
}

