/*
全局变量

参数控制

*/
package conf

import (
	//"database/sql"
	//"github.com/go-xorm/xorm"

	//"sync"
)

// 用于协程等待。
// var WG sync.WaitGroup

var Vendors []string
var ConfDBServer string

//////////////////////////////////////////

type Gin_t struct {
    Enable bool `json:"enable,omitempty"`
    Port   int `json:"port,omitempty"`
}

var Gin Gin_t

type UpdServer struct {
    Enable bool `json:"enable,omitempty"`
    Port   int `json:"port,omitempty"`
}

type TcpServer_t struct {
    Enable bool `json:"enable,omitempty"`
    Port   int `json:"port,omitempty"`
}

var TcpServer TcpServer_t


type DevServer struct {
    Enable   bool `json:"enable,omitempty"`
    Name     string `json:"name,omitempty"`
    Protocol string `json:"protocol,omitempty"`
	Port     int `json:"port,omitempty"`
}