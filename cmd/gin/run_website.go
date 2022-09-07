/*
仅作网站
*/
package gin

import (
	"net/http"
	"webdemo/common/conf"

	"github.com/gin-gonic/gin"
)

// 仅作为web静态服务器
func runWebSimple() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 指定html所在目录
	router.StaticFS("/website", http.Dir(conf.DataFileDir))

	// klog.Println("Server started")
	// router.Run(":" + strconv.Itoa(conf.Gin.Port))

	return router
}
