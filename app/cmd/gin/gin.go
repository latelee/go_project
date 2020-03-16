package gin

import (
    _ "fmt"
    "strconv"

    "github.com/latelee/go_project/pkg/com"
    "k8s.io/klog"
    "github.com/kubeedge/beehive/pkg/core"
    beehiveContext "github.com/kubeedge/beehive/pkg/core/context"

    "github.com/gin-gonic/gin"
    "net/http"
    
    "github.com/latelee/go_project/app/pkg/update"
    "github.com/latelee/go_project/app/conf"
)

type ginServer struct {
    enable bool
}

func init() {
    //core.Register(newginServer(true))
}

func newginServer(enable bool) *ginServer {
    return &ginServer{
        enable: enable,
    }
}

func Register(opts *conf.Gin) {
    initConfig(opts)
    core.Register(newginServer(opts.Enable))
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
    go doit()
    
    router := gin.Default()
	router.POST("/test", HelloWordPost)
    router.POST("/update", UpdateTest)
    //router.Run(":4000")
	router.Run(":" + strconv.Itoa(Config.Port))
}

func HelloWordPost (c *gin.Context) {
	c.String(http.StatusOK, "hello world post")
}


func UpdateTest(c *gin.Context) {
    var status int;
    var backInfo string;

    klog.Printf("got self update.\n");

    err := update.EnterUpgradeApp();
    if err == false {
        status = http.StatusInternalServerError;
        backInfo = "process failed";
    } else {
        status = http.StatusOK;
        backInfo = "process ok";
    }
    c.String(status, backInfo);
}

func doit() {
    for {
        select {
		case <-beehiveContext.Done():
			klog.Info("Stop gin1")
			return
		default:
		}
        klog.Infoln(".")
        com.Sleep(10000)
    }
}
