package devworker

import (
    _ "fmt"
    //"webdemo/pkg/com"
    "k8s.io/klog"
)

type DTWorker interface {
	Start()
}

type Worker struct {
	ReceiverChan  chan interface{}
	ConfirmChan   chan interface{}
	HeartBeatChan chan interface{}
}

struct DeviceTwin {

}

func (dt *DeviceTwin) RegisterDTModule(name string) {

}
