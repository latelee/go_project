package gin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"webdemo/common/conf"

	"github.com/gin-gonic/gin"
)

//////////////////////////////////////////
/*
一般性的post请求处理

curl http://127.0.0.1:9000/testing -X POST -H "Content-Type:application/json" -d  '{"id":"test_001", "op":"etc", "timestamp":12342134341234, "data":{"name":"foo", "addr":"bar", "code":450481, "age":100}}'

// 发送空 或单纯字符
curl http://127.0.0.1:9000/testing -X POST -d "abc"

// 发送文件
curl http://127.0.0.1:9000/testing -X POST -F "file=@foo.json"
*/
func HandleTest(ctx *gin.Context) {
	var s conf.BaseRequestqMsg
	var ss conf.MyInfo_t
	var err error

	var resdata gin.H

	// fmt.Println("Content-Type: ", ctx.Request.Header.Get("Content-Type"))

	contentType := ctx.Request.Header.Get("Content-Type")
	// 判断参数
	if strings.Contains(contentType, "multipart/form-data") { // 文件形式
		// 2种方式都可，但 ctx.Request.FormFile 可以得到文件句柄，可直接拷贝 指定的关键字为 file
		//file, err := ctx.FormFile("file")
		file, header, err := ctx.Request.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code": "-1",
				"msg":  err.Error(),
				"data": gin.H{}})
			return
		}

		//fmt.Printf("Request: %+v\n", ctx.Request);
		//fmt.Printf("Formfile: %+v | %+v |||  %v %v\n", file, header, err, reflect.TypeOf(file));

		// 拿到文件和长度
		var jsonfilename string = header.Filename
		mysize := header.Size
		fmt.Printf("filename: %s size: %d\n", jsonfilename, mysize)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code": "-1",
				"msg":  err.Error(),
				"data": gin.H{}})
			return
		}

		// 处理json文件
		jsonbuf := make([]byte, mysize)
		_, err = file.Read(jsonbuf)

		// fmt.Println("json file content: ", string(jsonbuf))

		// tmpDataJson, _ := json.Marshal(jsonbuf)
		err = json.Unmarshal(jsonbuf, &s)
		if err != nil {
			resdata = gin.H{
				"code": "0",
				"msg":  "parse json failed " + err.Error(),
				"data": gin.H{},
			}
		} else {
			resdata = gin.H{
				"code": "0",
				"msg":  "ok got multipart json",
				"data": gin.H{},
			}
		}

	} else if strings.Contains(contentType, "application/json") { // 纯json
		fmt.Println("got json request")
		if err = ctx.ShouldBindJSON(&s); err != nil {
			fmt.Printf("bind json failed: %v\n", err)
			return
		}

		// 组装返回json
		resdata = gin.H{
			"code": "0",
			"msg":  "ok got only json",
			"data": gin.H{},
		}
	} else { // 其它
		param, _ := ctx.GetRawData()
		fmt.Println("rawdata: ", param)
		// 组装返回json
		resdata = gin.H{
			"code": "-1",
			"msg":  "only support json",
			"data": gin.H{
				"yourcontext": string(param),
			},
		}
	}

	// 再解析Data字段
	tmpDataJson, _ := json.Marshal(s.Data)
	json.Unmarshal(tmpDataJson, &ss)

	fmt.Printf("got json: [%#v]\n", s)
	fmt.Printf("got json data: [%#v]\n", ss)

	ctx.JSON(http.StatusOK, resdata)
}

/*
发送json文件
curl http://127.0.0.1:9000/testing -X POST -F "file=@foo.json"

*/
func HandlePostFile(ctx *gin.Context) {

}
