package devServer

import (
    "sync"
    "github.com/latelee/go_project/app/conf"
)

var Config devConfig
var once sync.Once

type devConfig struct {
    conf.DevServer
    // 后可加其它字段
}

func initConfig(opts *conf.DevServer) {
    once.Do(func() {
        Config = devConfig {
            DevServer: *opts,
        }
    })
}