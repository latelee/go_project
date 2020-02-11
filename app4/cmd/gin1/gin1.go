package gin1

import (
    _ "fmt"
    "com"
    "k8s.io/klog"
    "github.com/kubeedge/beehive/pkg/core"
    
    "github.com/gin-gonic/gin"
    "net/http"
)

type ginTest1 struct {
    enable bool
    // 后可加其它字段
}

func init() {
    core.Register(newginTest1(true))
}

func newginTest1(enable bool) *ginTest1 {
    return &ginTest1{
        enable: enable,
    }
}

func Register() {
    core.Register(newginTest1(true))
}

func (a *ginTest1) Name() string {
    return "ginTest1"
}

func (a *ginTest1) Group() string {
    return "ginTest1"
}

// Enable indicates whether enable this module
func (a *ginTest1) Enable() bool {
    return a.enable
}

func (a *ginTest1) Start() {
    klog.Infoln("ginTest1...")
    //go doit()
    
    router := gin.Default()
	router.POST("/test", HelloWordPost)
	router.Run(":4000")
}

func HelloWordPost (c *gin.Context) {
	c.String(http.StatusOK, "hello world post")
}

func doit() {
    for {
        klog.Infoln(".")
        com.Sleep(1000)
    }
}
