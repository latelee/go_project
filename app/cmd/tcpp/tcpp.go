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

func Register(opts *conf.TcpServer) {
    initConfig(opts)
    core.Register(newtcpServer(opts.Enable))
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
	IpAndPort := "0.0.0.0:" + strconv.Itoa(Config.Port)
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
	klog.Infof("New TCP Connect from [%s:%s] ...", ip, port)

	for {
        select {
		case <-beehiveContext.Done():
			klog.Info("Stop tcp handle")
			return
		default:
		}
		buf := make([]byte, TCP_RECV_LEN)
		n, err := conn.Read(buf)
		if err != nil {
			if err.Error() != "EOF" {
				klog.Errorln(err)
			}
			break
		}
		klog.Infof("TCP Received from [%s] buf: %v %v", RemoteAddr, hex.Dump(buf[:n]), string(buf))

        // strings.Compare(string(buf), "world")
		if buf[0] == 0x68 {
            backbuf := "hello_back11111111111111111111111111111111111"
            klog.Info("send ", backbuf)
            conn.Write([]byte(backbuf))
		} else if  buf[0] == 0x77 {
            backbuf := "world_back22222222222222222222222222222222222222"
            klog.Info("send1 ", backbuf)
            conn.Write([]byte(backbuf))
        }
	}
}
