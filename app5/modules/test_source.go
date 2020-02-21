package modules

import (
	"k8s.io/klog"
	"time"

	"github.com/kubeedge/beehive/pkg/core"
	c "github.com/kubeedge/beehive/pkg/core/context"
	"github.com/kubeedge/beehive/pkg/core/model"
)

//Constants for module source and group
const (
	SourceModule = "sourcemodule"
	SourceGroup  = "sourcegroup"
)

type testModuleSource struct {
	//context *context.Context
    enable bool
}

func init() {
	core.Register(newtestModuleSource(true))
}

func newtestModuleSource(enable bool) *testModuleSource {
    return &testModuleSource{
        enable: enable,
    }
}

func (*testModuleSource) Name() string {
	return SourceModule
}

func (*testModuleSource) Group() string {
	return SourceGroup
}

func (a *testModuleSource) Enable() bool {
    return a.enable
}

func (m *testModuleSource) Send(msg string) {
    message := model.NewMessage("").SetRoute(SourceModule, "").
		SetResourceOperation("test", model.InsertOperation).FillBody(msg)
	c.Send(DestinationModule, *message)
}

func (m *testModuleSource) Start() {
    m.Send("test2")
    m.Send("test1")
    
	message := model.NewMessage("").SetRoute(SourceModule, "").
		SetResourceOperation("test", model.InsertOperation).FillBody("hello")
	c.Send(DestinationModule, *message)

	message = model.NewMessage("").SetRoute(SourceModule, "").
		SetResourceOperation("test", model.UpdateOperation).FillBody("how are you")
	resp, err := c.SendSync(DestinationModule, *message, 3*time.Second)

	if err != nil {
		klog.Printf("failed to send sync message, error:%v\n", err)
	} else {
		klog.Printf("get resp: %v\n", resp)
	}

	message = model.NewMessage("").SetRoute(SourceModule, DestinationGroup).
		SetResourceOperation("test", model.DeleteOperation).FillBody("fine")
	c.SendToGroup(DestinationGroup, *message)
}
