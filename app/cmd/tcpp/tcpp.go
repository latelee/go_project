package tcpp

import (
	_ "fmt"
	"strconv"
	_ "webdemo/pkg/com"
	"webdemo/app/conf"

    "k8s.io/klog"
    "github.com/kubeedge/beehive/pkg/core"
    beehiveContext "github.com/kubeedge/beehive/pkg/core/context"

	"net"
	"strings"
    //"encoding/hex"
)

const (
	TCP_RECV_LEN = 1024
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
    core.Register(newtcpServer(conf.TcpServer.Enable))
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

func (a *tcpServer) Cleanup() {
}

// TODO：添加断开的处理

func TcpServer() {
	IpAndPort := "0.0.0.0:" + strconv.Itoa(conf.TcpServer.Port)
	klog.Info("tcp listen on: ", IpAndPort)
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
	klog.Infof("New TCP Connect from [%s:%s]", ip, port)

	for {
        select {
		case <-beehiveContext.Done():
			klog.Info("Stop tcp handle")
            conn.Close()
			return
		default:
		}
        // 1.接收
		buf := make([]byte, TCP_RECV_LEN)
		n, err := conn.Read(buf)
		if err != nil {
            klog.Errorln(err)
			if err.Error() == "EOF" {
				conn.Close()
			}
			break
		}
		//klog.Infof("TCP Received from [%s] buf: %v", RemoteAddr, hex.Dump(buf[:n]))
        // 2.处理
        backbuf, backlen := handle(buf, n)
        
        // 3.返回
        if backlen > 0 {
            conn.Write(backbuf)
        }
	}
}
