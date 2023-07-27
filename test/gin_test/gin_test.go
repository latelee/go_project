/*
结构体
{
    "id": "test_001",
    "op": "etc",
    "timestamp": 12342134341234,
    "data": {
        "name": "foo",
        "addr": "bar",
        "code": 450481,
        "age": 100
    }
}

curl http://127.0.0.1:9000/foo -X POST -H "Content-Type:application/json" -d  '{"id":"test_001", "op":"etc", "timestamp":12342134341234, "data":{"name":"foo", "addr":"bar", "code":450481, "age":100}}'

curl http://127.0.0.1:9000/bar -X POST -H "Content-Type:application/json" -d  '{"id":"test_001", "op":"etc", "timestamp":12342134341234, "data":{"name":"foo", "addr":"bar", "code":450481, "age":100}}'

*/

package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

var g_port string = "9000"

type MyRequest_t struct {
	Id        string    `json:"id"`
	Op        string    `json:"op"`
	Timestamp int       `json:"timestamp"`
	Data      ReqData_t `json:"data"`
}

type ReqData_t struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
	Code int    `json:"code"`
	Age  int    `json:"age"`
}

func routerPost(r *gin.Engine) {
	r.POST("/foo", HandleGinShouldBindJSON)
	r.POST("/bar", HandleGinUnmarshal)
}

func initGin() {
	fmt.Println("run gin")
	router := gin.New()
	routerPost(router)

	router.Run(":" + g_port)
}

func HandleGinShouldBindJSON(ctx *gin.Context) {
	var request MyRequest_t
	var err error
	ctxType := ctx.Request.Header.Get("Content-Type")
	if strings.Contains(ctxType, "application/json") { // 纯 json
		// 先获取总的json
		if err = ctx.ShouldBindJSON(&request); err != nil {
			fmt.Printf("ShouldBindJSON failed: %v\n", err)
			return
		}

		fmt.Printf("ShouldBindJSON: request: #%v\n", request)
	} else {
		fmt.Println("非json")
		return
	}
}

func HandleGinUnmarshal(ctx *gin.Context) {
	var request MyRequest_t
	var err error
	var reqbuffer []byte
	ctxType := ctx.Request.Header.Get("Content-Type")
	if strings.Contains(ctxType, "application/json") { // 纯 json

		reqbuffer, err = ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			fmt.Printf("ReadAll body failed: %v\n", err)
			return
		}
		err = json.Unmarshal(reqbuffer, &request)
		if err != nil {
			fmt.Printf("Unmarshal to request failed: %v\n", err)
			return
		}
		fmt.Printf("Unmarshal request: #%v\n", request)
	} else {
		fmt.Println("非json")
		return
	}
}

func TestGin(t *testing.T) {
	fmt.Println("test of gin")

	initGin()

	// for {
	// 	time.Sleep(time.Duration(1) * time.Second)
	// }
}
