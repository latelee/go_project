package gin

import (
    _ "fmt"
    "strconv"

//    "github.com/latelee/go_project/pkg/com"
    "k8s.io/klog"
    "github.com/kubeedge/beehive/pkg/core"
    beehiveContext "github.com/kubeedge/beehive/pkg/core/context"

    "net/http"
    "github.com/gin-gonic/gin"
    
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
    //go a.doit()
    
    router := gin.Default()
	router.POST("/test", HelloWordPost)
    router.GET("/test", HelloWordGet)
    router.POST("/update", UpdateTest)
    
    // 组
    v1 := router.Group("/api/v1/userinfo")
    {
        v1.GET("/", FetchAllUsers)
        v1.GET("/:id", FetchSingleUser)
    }

    //InitWS()
    // 组
    v2 := router.Group("/device/ws")
    {
        v2.GET("/:id", WSHandler)
    }
    //router.Run(":4000")
	router.Run(":" + strconv.Itoa(Config.Port))
}

func (a *ginServer) Cleanup() {
    delListAll()
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

func HelloWordPost(c *gin.Context) {
	c.String(http.StatusOK, "hello world post")
}

func HelloWordGet(c *gin.Context) {
    c.String(http.StatusOK, "hello world get")
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
