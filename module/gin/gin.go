package gin

import (
	_ "fmt"
	"strconv"
	"webdemo/common/conf"
	"webdemo/pkg/com"
	"webdemo/pkg/klog"

	//"net/http"
	"github.com/gin-gonic/gin"
)

type GinServer_t struct {
	webtype int
	*gin.Engine
}

func NewGinServer() *GinServer_t {
	return &GinServer_t{}
}

func (this *GinServer_t) Start(webtype int) {
	this.webtype = webtype
	// 业务初始化
	router := this.initBusy()

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

func (this *GinServer_t) initBusy() *gin.Engine {

	if this.webtype == WEB_WEBSITE { // 可作为普通的web网站服务器使用
		klog.Println("Running for static web site...")
		return runWebSimple()
	} else if this.webtype == WEB_CLINET { // 客户端，目前没有用到
		klog.Println("Running for post client...")
		// return cli.Client()
	} else { // 全部，含前后端
		klog.Println("Running for post test...")
		return runAll()
	}

	return nil
}
