package tcpp

import (
    _ "fmt"
    _ "com"
    "k8s.io/klog"
    "github.com/kubeedge/beehive/pkg/core"

	"net"
	"strings"
    "encoding/hex"
)

const (
	TCP_RECV_LEN        = 512
)

type tcpServer struct {
    enable bool
    // 后可加其它字段
}

func newtcpServer(enable bool) *tcpServer {
    return &tcpServer{
        enable: enable,
    }
}

func Register() {
    core.Register(newtcpServer(true))
}

func (a *tcpServer) Name() string {
    return "tcpServer"
}

func (a *tcpServer) Group() string {
    return "tcpServer"
}

// Enable indicates whether enable this module
func (a *tcpServer) Enable() bool {
    return a.enable
}

func (a *tcpServer) Start() {
    TcpServer()
}

// TODO：添加断开的处理

func TcpServer() {
	IpAndPort := "0.0.0.0:8000"
	klog.Info("tcp listen: ", IpAndPort)
	ln, err := net.Listen("tcp", IpAndPort)
	if err != nil {
		klog.Errorln("tcp listen error: ", err)
		return
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			klog.Errorln(err)
			continue
		}

		go handleConnection(conn)
	}
}

// rev handle
func handleConnection(conn net.Conn) {
	defer conn.Close()

	RemoteAddr := conn.RemoteAddr().String()
	ip := (strings.Split(RemoteAddr, ":"))[0]
    port := (strings.Split(RemoteAddr, ":"))[1]
	klog.Infof("New TCP Connect from [%s:%s] ......", ip, port)

	for {
		buf := make([]byte, TCP_RECV_LEN)
		n, err := conn.Read(buf)
		if err != nil {
			if err.Error() != "EOF" {
				klog.Errorln(err)
			}
			break
		}
		klog.Infof("TCP Received from [%s] buf: %v", RemoteAddr, hex.Dump(buf[:n]))
	}
}
