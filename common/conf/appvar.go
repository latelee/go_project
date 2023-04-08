/*
全局变量

参数控制

*/
package conf

import (
	//"database/sql"
	//"github.com/go-xorm/xorm"

	//"sync"
	"time"

	"github.com/spf13/viper"
)

// 用于协程等待。
// var WG sync.WaitGroup

var Vendors []string
var ConfDBServer string

var RunMode string
var Args []string

// 当前运行的操作系统 目前只处理 windows或linux
// 注：windows 系统无法计算路径
var RunningOS string
var RunningArch string

var AppVersion string
var AppVersionInfo string

var CfgFile string = "./config.yaml"

var HostName string
var StartTime time.Time

var LogDir string = ""

// 全局的配置句柄
var Config *viper.Viper

var DataFileDir string

var CurDir string

// https
var HttpsEnable bool
var HttpsCertFile string
var HttpsKeyFile string

////////////////////////
const DEFUALT_PORT int = 9000

//////////////////////////////////////////

// 命令列表，包括名称，帮助信息
type UserCmdFunc struct {
	Name      string
	ShortHelp string
	// LongHelp string
	Func func(args []string)

	NeedDb bool
}

type Gin_t struct {
	Enable bool `json:"enable,omitempty"`
	Port   int  `json:"port,omitempty"`
}

var Gin Gin_t

type UpdServer struct {
	Enable bool `json:"enable,omitempty"`
	Port   int  `json:"port,omitempty"`
}

type TcpServer_t struct {
	Enable bool `json:"enable,omitempty"`
	Port   int  `json:"port,omitempty"`
}

var TcpServer TcpServer_t

type DevServer struct {
	Enable   bool   `json:"enable,omitempty"`
	Name     string `json:"name,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Port     int    `json:"port,omitempty"`
}

///////////

// 封装一层，实际数据在Data中，通过Op区分
type BaseRequestMsg struct {
	Id        string      `json:"id"`
	Op        string      `json:"op"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type BaseRespondMsg struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 测试
type MyInfo_t struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
	Code int    `json:"code"`
	Age  int    `json:"age"`
}

///////////////////////////////////////////
