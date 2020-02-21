package modules

import (
	"k8s.io/klog"

	"github.com/kubeedge/beehive/pkg/core"
	c "github.com/kubeedge/beehive/pkg/core/context"
)

//Constant for test module destination group name
const (
	DestinationGroupModule = "destinationgroupmodule"
)

type testModuleDestGroup struct {
	//context *context.Context
    enable bool
}

func init() {
	core.Register(newtestModuleDestGroup(true))
}

func newtestModuleDestGroup(enable bool) *testModuleDestGroup {
    return &testModuleDestGroup{
        enable: enable,
    }
}

func (*testModuleDestGroup) Name() string {
	return DestinationGroupModule
}

func (*testModuleDestGroup) Group() string {
	return DestinationGroup
}

func (a *testModuleDestGroup) Enable() bool {
    return a.enable
}

func (m *testModuleDestGroup) Start() {
	
	message, err := c.Receive(DestinationGroupModule)
	klog.Printf("destination group module receive message:%v error:%v\n", message, err)
	if message.IsSync() {
		resp := message.NewRespByMessage(&message, "10 years old")
		c.SendResp(*resp)
	}
}

