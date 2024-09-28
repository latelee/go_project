/*

 */

package gin

import (
	"net/http"
	"strconv"
	"webdemo/common/conf"
	"webdemo/pkg/com"

	"github.com/gin-gonic/gin"
)

func runAll() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	staticdir := conf.DataFileDir
	// 指定html所在目录
	router.LoadHTMLGlob(staticdir + "/static/html/*")
	router.StaticFS("/css", http.Dir(staticdir+"/static/css"))
	router.StaticFS("/js", http.Dir(staticdir+"/static/js"))
	router.StaticFS("/img", http.Dir(staticdir+"/static/img"))
	// 是否加上ico的？
	router.StaticFile("/favicon.ico", staticdir+"/static/img/favicon.ico")

	router.StaticFS("/website", http.Dir("./data/dist")) // TODO 目录是否由外部传入？

	routerPage(router)
	routerPost(router)

	return router
}

func (this *GinServer_t) RunAll() {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// dir := conf.Config.GetString("modules.gin.website")
	// // 指定html所在目录
	// router.StaticFS("/", http.Dir(dir))

	staticdir := conf.DataFileDir

	// 指定html所在目录
	router.LoadHTMLGlob(staticdir + "/static/html/*")
	router.StaticFS("/css", http.Dir(staticdir+"/static/css"))
	router.StaticFS("/js", http.Dir(staticdir+"/static/js"))
	router.StaticFS("/img", http.Dir(staticdir+"/static/img"))
	// 是否加上ico的？
	router.StaticFile("/favicon.ico", staticdir+"/static/img/favicon.ico")

	router.StaticFS("/website", http.Dir("./data/dist")) // TODO 目录是否由外部传入？

	routerPage(router)
	routerPost(router)

	// 运行服务，支持2种方式
	if conf.HttpsEnable {
		if com.IsExist(conf.HttpsCertFile) && com.IsExist(conf.HttpsKeyFile) {
			router.RunTLS(":"+strconv.Itoa(conf.Gin.Port), conf.HttpsCertFile, conf.HttpsKeyFile)
			return
		}

	}
	router.Run(":" + strconv.Itoa(conf.Gin.Port))
}
