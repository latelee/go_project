package gin

import (
	_ "fmt"
	"strconv"

	"webdemo/pkg/com"
	"webdemo/pkg/klog"

	"github.com/kubeedge/beehive/pkg/core"
	beehiveContext "github.com/kubeedge/beehive/pkg/core/context"

	//"net/http"
	"github.com/gin-gonic/gin"

	cli "webdemo/cmd/gin/post_client"
	"webdemo/common/conf"
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

	// 业务初始化
	router := initBusy()

	if router == nil {
		return
	}
	// 运行服务，支持2种方式
	if conf.HttpsEnable {
		if com.IsExist(conf.HttpsCertFile) && com.IsExist(conf.HttpsKeyFile) {
			router.RunTLS(":"+strconv.Itoa(conf.Gin.Port), conf.HttpsCertFile, conf.HttpsKeyFile)
			return
		}

	}
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

func initBusy() *gin.Engine {

	if conf.RunMode == "website" { // 可作为普通的web网站服务器使用
		klog.Println("Running for static web site...")
		return runWebSimple()
	} else if conf.RunMode == "client" { // 可作为普通的web网站服务器使用
		klog.Println("Running for post client...")
		return cli.Client()
	} else { // add more...
		klog.Println("Running for post test...")
		return runAll()
	}

	// return nil
}
