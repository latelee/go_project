package udpp

import (
    _ "fmt"
    "com"
    "k8s.io/klog"
    "github.com/kubeedge/beehive/pkg/core"
    beehiveContext "github.com/kubeedge/beehive/pkg/core/context"

	"net"
)

const (
	UDP_RECV_LEN = 1024
)


type udpServer struct {
    enable bool
    // 后可加其它字段
}

func init() {
    //core.Register(newudpServer(true))
}

func newudpServer(enable bool) *udpServer {
    return &udpServer{
        enable: enable,
    }
}

func Register() {
    core.Register(newudpServer(true))
}

func (a *udpServer) Name() string {
    return "udpServer"
}

func (a *udpServer) Group() string {
    return "udpServer"
}

// Enable indicates whether enable this module
func (a *udpServer) Enable() bool {
    return a.enable
}

func (a *udpServer) Start() {
    UdpServer()
}

// ..
func UdpServer() {
	IpAndPort := ":10086"
	sAddr, err := net.ResolveUDPAddr("udp", IpAndPort)
	if err != nil {
		klog.Infoln("ResolveUDPAddr: ", err)
		return
	}

	klog.Infoln("udp listen: ", IpAndPort)
	ln, err := net.ListenUDP("udp", sAddr)
	if err != nil {
		klog.Errorln("udp listen error: ", err)
		return
	}

	// recieve
	go func() {
		for {
            select {
            case <-beehiveContext.Done():
                klog.Info("Stop udp recv")
                return
            default:
            }
			data := make([]byte, UDP_RECV_LEN)
			n, cAddr, err := ln.ReadFrom(data)
			if err != nil {
				klog.Errorln("udp ReadFrom error: ", err)
				break
			}

			str := string(data[:n])
			klog.Infoln("UDP Received from: ", cAddr.String(), " ", str)
		}
	}()

	// send
	for {
        select {
		case <-beehiveContext.Done():
			klog.Info("Stop upd send")
			return
		default:
		}
        klog.Info("send udp...")
		rAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:1024")
		_, err = ln.WriteToUDP([]byte("hello_all"), rAddr)
		if err != nil {
			klog.Errorln("udp WriteToUDP error: ", err)
			//return
		}

		com.Sleep(10000)
	}
}
