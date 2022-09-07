/*

 */

package gin

import (
	"net/http"
	"webdemo/common/conf"

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
