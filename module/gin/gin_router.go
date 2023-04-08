/*
和路由相关，如页面和post请求入口，具体的散见其它文件。
*/

package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//////////////////////////////////////////////
func routerPage(r *gin.Engine) {
	r.GET("/", pageIndex)
	r.GET("/page2", page2)
	r.GET("/testing", pageIndex)
}

func routerPost(r *gin.Engine) {
	routerPostOld(r)

	//
	r.POST("/api/user/login", HandleLogin)
	r.POST("/api/user/logout", HandleLogout)

	r.POST("/testing", HandleTest)
	r.POST("/postfile", HandlePostFile)

}

// 旧的
func routerPostOld(r *gin.Engine) {
	//r.POST("/test", HelloWordPost)
	//r.GET("/test", HelloWordGet)

	// 组
	v1 := r.Group("/device/v1/devlist")
	{
		v1.GET("/", GetAllDevices)
		v1.GET("/:id", GetSingleDevice)
	}

	r.POST("/device/v1/get", DeviceGetHandle)
	r.POST("/device/v1/set", DeviceSetHandle)

	/*
	   // 组
	   v2 := r.Group("/device/ws")
	   {
	       v2.GET("/:id", WSHandler)
	   }
	*/
}

/////////////////////////////////////////////////

func pageIndex(ctx *gin.Context) {
	page := "index.html"

	ctx.HTML(http.StatusOK, page, gin.H{})
}

func page2(ctx *gin.Context) {
	page := "page2.html"

	ctx.HTML(http.StatusOK, page, gin.H{})
}
