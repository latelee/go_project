/*
仅作网站
*/
package gin

import (
	"net/http"
	"strconv"
	"webdemo/common/conf"
	"webdemo/pkg/com"

	"github.com/gin-gonic/gin"
)

func (this *GinServer_t) RunStaticWebSite() {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	dir := conf.Config.GetString("modules.gin.website")
	// 指定html所在目录
	router.StaticFS("/", http.Dir(dir))

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

// 仅作为web静态服务器
func runWebSimple() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	dir := conf.Config.GetString("modules.gin.website")
	// 指定html所在目录
	router.StaticFS("/", http.Dir(dir))

	return router
}
