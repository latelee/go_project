package gin

import (
    _ "fmt"
    "strconv"

    "com"
    "k8s.io/klog"
    "github.com/kubeedge/beehive/pkg/core"
    beehiveContext "github.com/kubeedge/beehive/pkg/core/context"

    "github.com/gin-gonic/gin"
    "net/http"
    
    "github.com/latelee/myproject/app/pkg/update"
    "github.com/latelee/myproject/app/conf"
)

type ginServer struct {
    enable bool
    conf.Gin
    // 后可加其它字段
}

func init() {
    //core.Register(newginServer(true))
}

func newginServer(opts *conf.Gin) *ginServer {
    return &ginServer{
        enable: opts.Enable,
        Gin:    *opts,
    }
}

func Register(opts *conf.Gin) {
    core.Register(newginServer(opts))
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
	router.Run(":" + strconv.Itoa(a.Gin.Port))
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
