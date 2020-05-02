package tcpp

import (
    "sync"
    "webdemo/app/conf"
)

var Config tcpConfig
var once sync.Once

type tcpConfig struct {
    conf.TcpServer
    // 后可加其它字段
}

func initConfig(opts *conf.TcpServer) {
    once.Do(func() {
        Config = tcpConfig {
            TcpServer: *opts,
        }
    })
}