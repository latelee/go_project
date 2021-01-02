package gin

import (
    _ "fmt"
    "strconv"

//    "webdemo/pkg/com"
    "k8s.io/klog"
    "github.com/kubeedge/beehive/pkg/core"
    beehiveContext "github.com/kubeedge/beehive/pkg/core/context"

    //"net/http"
    "github.com/gin-gonic/gin"

    "webdemo/app/conf"
)

type ginServer struct {
    enable bool
}

func newginServer(enable bool) *ginServer {
    return &ginServer{
        enable: enable,
    }
}

func Register() {
    core.Register(newginServer(conf.Gin.Enable))
}

func (a *ginServer) Name() string {
    return "ginServer"
}

func (a *ginServer) Group() string {
    return "ginServer"
}

// Enable indicates whether enable this module
func (a *ginServer) Enable() bool {
    return a.enable
}

func (a *ginServer) Start() {
    klog.Infoln("ginServer...")
    //go a.doit()
    
    router := gin.Default()
	//router.POST("/test", HelloWordPost)
    //router.GET("/test", HelloWordGet)
    
    // 组
    v1 := router.Group("/device/v1/devlist")
    {
        v1.GET("/", GetAllDevices)
        v1.GET("/:id", GetSingleDevice)
    }

    router.POST("/device/v1/get", DeviceGetHandle)
    router.POST("/device/v1/set", DeviceSetHandle)
        
    /*
    // 组
    v2 := router.Group("/device/ws")
    {
        v2.GET("/:id", WSHandler)
    }
    */

	router.Run(":" + strconv.Itoa(conf.Gin.Port))
}

func (a *ginServer) Cleanup() {

}

// 作用：做监听，退出时清理
// 是否需要在模块入口加处理？ 
// 理论上在beehive中添加回调即可，不用再写
func (a *ginServer) doit() {
    for {
        select {
		case <-beehiveContext.Done():
			klog.Infof("Stop %s", a.Name())
			return
		default:
		}
        //klog.Infoln(".")
        //com.Sleep(10000)
    }
}
