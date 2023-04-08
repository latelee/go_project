/*
用户登录相关
*/

package gin

import (

	//"time"

	//"webdemo/pkg/klog"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 用户信息
// TODO 理论上，用结构体对应请求json，是最快的方便，后续再优化
type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// ctx.Request.Method == "GET"
func HandleLogin(ctx *gin.Context) {
	// json格式：
	data, _ := ioutil.ReadAll(ctx.Request.Body)
	fmt.Printf("ctx.Request.body: %v\n", string(data))
	// 再解析data
	var v map[string]interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		responseFailedMsg(-1, "json parse error", ctx)
		return
	}
	fmt.Printf("got get msg: %v\n", v)

	userName := v["userName"]
	passWord := v["passWord"]
	//  先固定！！！
	if (userName != "admin" || passWord != "admin") &&
		(userName != "test" || passWord != "test") &&
		(userName != "jky" || passWord != "jky@2021") &&
		(userName != "latelee" || passWord != "latelee@2021") {
		fmt.Println("user info error")
		responseFailedMsg(-1, "user info error", ctx)
		return
	}

	level := "test"
	level_code := 3
	if userName == "latelee" {
		level = "verysuper_admin"
		level_code = 0
	} else if userName == "jky" {
		level = "super_admin"
		level_code = 1
	} else if userName == "admin" {
		level = "admin"
		level_code = 2
	}

	// // form形式 数据在url中出现
	// fmt.Printf("c.Request.Method: %v\n", ctx.Request.Method)
	// fmt.Printf("c.Request.ContentType: %v\n", ctx.ContentType())
	// ctx.Request.ParseForm()
	// for k, v := range ctx.Request.PostForm {
	// 	fmt.Printf("k:%v\n", k)
	// 	fmt.Printf("v:%v\n", v)
	// }
	// // 不知为何，用PostForm获取不到
	// username := ctx.PostForm("userName")
	// password := ctx.PostForm("passWord")
	// fmt.Printf("got param: %+v\n %s %s \n", ctx.Request.Body, username, password)

	// ctx.PureJSON
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"code": 0,
			"msg":  "ok",
			"data": gin.H{
				"username":   userName,
				"token":      "5b21ca65059261c92b8e996669129783",
				"level":      level,
				"level_code": level_code,
			},
		},
	)
}

func HandleLogout(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"code": 0,
			"msg":  "ok",
			"data": gin.H{
				"result": "ok",
			},
		},
	)
}
