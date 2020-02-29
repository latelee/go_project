package gin

import (
    "sync"
    "github.com/latelee/myproject/app/conf"
)

var Config ginConfig
var once sync.Once

type ginConfig struct {
    conf.Gin
    // 后可加其它字段
}

func initConfig(opts *conf.Gin) {
    once.Do(func() {
        Config = ginConfig {
            Gin: *opts,
        }
    })
}