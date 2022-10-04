/*
全局变量

参数控制

*/
package conf

import (
	//"database/sql"
	//"github.com/go-xorm/xorm"

	//"sync"
	"github.com/spf13/viper"
)

// 用于协程等待。
// var WG sync.WaitGroup

var Vendors []string
var ConfDBServer string

// 配置参数
var RunningOS string
var RunningARCH string
var RunMode string
var Args []string

// 全局的配置句柄
var Config *viper.Viper

var DataFileDir string

var AppVersion string

// https
var HttpsEnable bool
var HttpsCertFile string
var HttpsKeyFile string

//////////////////////////////////////////

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
