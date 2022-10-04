package gin

import (
	"encoding/json"
	"fmt"
	"strings"
	"webdemo/common/conf"

	"github.com/gin-gonic/gin"
)

//////////////////////////////////////////
/*
一般性的post请求处理
curl http://127.0.0.1:9000/testing -X POST -H "Content-Type:application/json" -d  '{"id":"test_001", "op":"etc", "timestamp":12342134341234, "data":{"name":"foo", "addr":"bar", "res_code":450481, "age":100}}'

// 发送空 或单纯字符
curl http://127.0.0.1:9000/testing -X POST -d "abc"

// 发送文件
curl http://127.0.0.1:9000/testing -X POST -F "file=@foo.json"
*/
func HandleTest(ctx *gin.Context) {
	var s conf.BaseRequestMsg
	var ss conf.MyInfo_t
	var err error

	var res_code int
	var res_msg string
	res_data := make(map[string]interface{})

	var tmpDataJson []byte
	// fmt.Println("Content-Type: ", ctx.Request.Header.Get("Content-Type"))

	contentType := ctx.Request.Header.Get("Content-Type")
	// 判断参数
	if strings.Contains(contentType, "multipart/form-data") { // 文件形式
		fmt.Println("got file json request")
		// 2种方式都可，但 ctx.Request.FormFile 可以得到文件句柄，可直接拷贝 指定的关键字为 file
		//file, err := ctx.FormFile("file")
		file, header, err := ctx.Request.FormFile("file")
		if err != nil {
			// 使用封装的函数返回
			res_code = -1
			res_msg = err.Error()
			// 可直接用ctx.JSON组装并返回，下同
			// ctx.JSON(http.StatusBadRequest, gin.H{
			// 	"res_code": "-1",
			// 	"res_msg":  err.Error(),
			// 	"data": gin.H{}})
			// return

			goto end
		}

		//fmt.Printf("Request: %+v\n", ctx.Request);
		//fmt.Printf("Formfile: %+v | %+v |||  %v %v\n", file, header, err, reflect.TypeOf(file));

		// 拿到文件和长度
		var jsonfilename string = header.Filename
		mysize := header.Size
		fmt.Printf("filename: %s size: %d\n", jsonfilename, mysize)

		if err != nil {
			res_code = -1
			res_msg = err.Error()
			goto end
		}

		// 处理json文件
		jsonbuf := make([]byte, mysize)
		_, err = file.Read(jsonbuf)

		// fmt.Println("json file content: ", string(jsonbuf))

		// tmpDataJson, _ := json.Marshal(jsonbuf)
		err = json.Unmarshal(jsonbuf, &s)
		if err != nil {
			res_code = -1
			res_msg = "parse json failed " + err.Error()
			goto end
		}

		res_code = 0
		res_msg = "ok got multipart json"

		// 特殊，以文件形式请求，返回文件形式
		back_file := fmt.Sprintf("res_%v", jsonfilename)
		ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("form-data; name=\"file\"; filename=\"%s\"", back_file))

		// 用现成的文件
		// ctx.File(back_file)
		// 内部实际调用http.ServeFile
		// http.ServeFile(ctx.Writer, ctx.Request, back_file)

		//return

	} else if strings.Contains(contentType, "application/json") { // 纯json
		fmt.Println("got json request")
		if err = ctx.ShouldBindJSON(&s); err != nil {
			fmt.Printf("ShouldBindJSON failed: %v\n", err)
			res_code = -1
			res_msg = "ShouldBindJSON failed " + err.Error()
		}

		res_code = 0
		res_msg = "ok got only json"
	} else { // 其它
		param, _ := ctx.GetRawData()
		fmt.Println("rawdata: ", param)

		res_code = -1
		res_msg = "only support json"
		res_data["yourcontext"] = string(param)
		goto end
	}

	// 再解析Data字段
	tmpDataJson, _ = json.Marshal(s.Data)
	json.Unmarshal(tmpDataJson, &ss)

	fmt.Printf("got json: [%#v]\n", s)
	fmt.Printf("got json data: [%#v]\n", ss)

end:
	response(res_code, res_msg, res_data, ctx)
	// ctx.JSON(http.StatusOK, resdata)
}

/*
发送json文件
curl http://127.0.0.1:9000/testing -X POST -F "file=@foo.json"

*/
func HandlePostFile(ctx *gin.Context) {

}
