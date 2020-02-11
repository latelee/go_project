package modules

import (
	"fmt"

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
	
	message, err := c.Receive(DestinationModule)
	fmt.Printf("destination module receive message:%v error:%v\n", message, err)
	message, err = c.Receive(DestinationModule)
	fmt.Printf("destination module receive message:%v error:%v\n", message, err)
	resp := message.NewRespByMessage(&message, "fine")
	if message.IsSync() {
		c.SendResp(*resp)
	}

	message, err = c.Receive(DestinationModule)
	fmt.Printf("destination module receive message:%v error:%v\n", message, err)
	if message.IsSync() {
		resp = message.NewRespByMessage(&message, "fine")
		c.SendResp(*resp)
	}

	//message, err = c.Receive(DestinationModule)
	//fmt.Printf("destination module receive message:%v error:%v\n", message, err)
	//if message.IsSync() {
	//	resp = message.NewRespByMessage(&message, "20 years old")
	//	c.SendResp(*resp)
	//}
}

