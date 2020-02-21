package server

import (
    _ "fmt"
    //"com"
    "k8s.io/klog"

    "github.com/kubeedge/beehive/pkg/core"
    
    update "github.com/latelee/myproject/app4/cmd/update"
    
    "github.com/gin-gonic/gin"
    "net/http"
)

type server struct {
    enable bool
    // 后可加其它字段
}

func init() {
    //core.Register(newServer(true))
}

func newServer(enable bool) *server {
    return &server{
        enable: enable,
    }
}

func Register() {
    core.Register(newServer(true))
}

func (a *server) Name() string {
    return "server"
}

func (a *server) Group() string {
    return "server"
}

// Enable indicates whether enable this module
func (a *server) Enable() bool {
    return a.enable
}

func (a *server) Start() {
    klog.Infoln("server...")

    router := gin.Default()
	router.POST("/update", UpdateTest)
	router.Run(":4000")
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
