package gin

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"webdemo/common/conf"
	"webdemo/pkg/com"

	"github.com/gin-gonic/gin"
)

// 发送buffer中的json
func post_from_buffer(url string, caCertFile, certFile, keyFile string, buffer []byte) (respFile, respBody string, err error) {
	respFile = ""
	respBody = ""
	err = nil

	// 发送一个POST请求
	req, err := http.NewRequest("POST", url, bytes.NewReader(buffer))
	if err != nil {
		err = errors.New(fmt.Sprintf("NewRequest failed: %v\n", err.Error()))
	}

	// 设置Header，判断json
	if buffer[0] == '{' {
		req.Header.Set("Content-Type", "application/json")
	}

	var client http.Client

	// 判断是否需要证书
	if url[0:4] == "https" {
		pool := x509.NewCertPool()
		var caCrt []byte
		caCrt, err = ioutil.ReadFile(caCertFile)
		if err != nil {
			fmt.Println("ReadFile err:", err)
			return
		}
		// 解析证书
		pool.AppendCertsFromPEM(caCrt)

		var cliCrt tls.Certificate // 具体的证书加载对象
		cliCrt, err = tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      pool,
				Certificates: []tls.Certificate{cliCrt},
			}, // failed
			// 客户端跳过对证书的校验
			// TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ok
		}

		// 用这种格式设置超时时间为10秒
		client = http.Client{
			Timeout:   10 * time.Second,
			Transport: tr,
		}
	} else {
		client = http.Client{
			Timeout: 10 * time.Second,
		}
	}
	// 执行请求
	resp_inn, err := client.Do(req)
	if err != nil {
		err = errors.New(fmt.Sprintf("send post failed: %v\n", err.Error()))
		return
	}

	defer client.CloseIdleConnections() // 如有相同请求的来，则要调用此句关闭空链接
	defer resp_inn.Body.Close()

	if item, ok := resp_inn.Header["Content-Disposition"]; ok {
		idx := strings.Index(item[0], "filename=")
		idx += len("filename=")
		respFile = item[0][idx:]
		// fmt.Printf("http resp: %v %s\n", respFile[0], respFile[0][idx:])
	}

	// 读取返回内容
	d, err := ioutil.ReadAll(resp_inn.Body)
	if err != nil {
		err = errors.New(fmt.Sprintf("read response failed: %v\n", err.Error()))
		return
	}

	respBody = string(d)
	err = nil

	return
}

// 将buffer以文件形式发送
func post_filedata_from_buffer(url string, caCertFile, certFile, keyFile string, postfilename string, buffer []byte) (respFile, respBody string, err error) {
	respFile = ""
	respBody = ""
	err = nil
	var w io.Writer
	var buff bytes.Buffer
	// 创建一个Writer
	writer := multipart.NewWriter(&buff)

	// file=@xxxx
	w, err = writer.CreateFormFile("file", postfilename)

	// klog.Println("filename: ", postfile)
	if err != nil {
		err = errors.New(fmt.Sprintf("create form failed: %v\n", err.Error()))
		return
	}

	// 把文件内容写入文件中
	w.Write(buffer)
	// com.Sleep(2)
	// 在此处关闭，可能会写完文件，不能用defer，否则文件可能没写完，会提示unexpected EOF
	writer.Close()

	fmt.Println("post file buf: ", string(buff.Bytes()))
	// 发送一个POST请求
	req, err := http.NewRequest("POST", url, &buff)
	if err != nil {
		err = errors.New(fmt.Sprintf("NewRequest failed: %v\n", err.Error()))
	}

	// 设置Header，自动获取，不用手动写
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var client http.Client

	// 判断是否需要证书
	if url[0:4] == "https" {
		pool := x509.NewCertPool()
		var caCrt []byte
		caCrt, err = ioutil.ReadFile(caCertFile)
		if err != nil {
			fmt.Println("ReadFile err:", err)
			return
		}
		// 解析证书
		pool.AppendCertsFromPEM(caCrt)

		var cliCrt tls.Certificate // 具体的证书加载对象
		cliCrt, err = tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      pool,
				Certificates: []tls.Certificate{cliCrt},
			}, // failed
			// 客户端跳过对证书的校验
			// TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ok
		}

		// 用这种格式设置超时时间为10秒
		client = http.Client{
			Timeout:   10 * time.Second,
			Transport: tr,
		}
	} else {
		client = http.Client{
			Timeout: 10 * time.Second,
		}
	}
	// 执行请求
	resp_inn, err := client.Do(req)
	if err != nil {
		err = errors.New(fmt.Sprintf("send post failed: %v\n", err.Error()))
		return
	}

	defer client.CloseIdleConnections() // 如有相同请求的来，则要调用此句关闭空链接
	defer resp_inn.Body.Close()

	if item, ok := resp_inn.Header["Content-Disposition"]; ok {
		idx := strings.Index(item[0], "filename=")
		idx += len("filename=")
		respFile = item[0][idx:]
		// fmt.Printf("http resp: %v %s\n", respFile[0], respFile[0][idx:])
	}

	// 读取返回内容
	d, err := ioutil.ReadAll(resp_inn.Body)
	if err != nil {
		err = errors.New(fmt.Sprintf("read response failed: %v\n", err.Error()))
		return
	}

	respBody = string(d)
	err = nil

	return
}

/*
以文件形式发送请求，curl等效命令：
curl http://127.0.0.1:9000/testing -X POST -F "file=@foo.json"


mypath 为json完整文件路径

和 post_filedata_from_buffer 功能相同
*/
func post_filedata_from_file(url string, caCertFile, certFile, keyFile string, mypath string) (respFile, respBody string, err error) {
	respFile = ""
	respBody = ""
	err = nil

	if !com.IsExist(mypath) {
		err = errors.New(fmt.Sprintf("%v not found", mypath))
		return
	}

	filename := filepath.Base(mypath)
	data, err := ioutil.ReadFile(mypath)
	if err != nil {
		err = errors.New(fmt.Sprintf("ReadFile failed: %v\n", err.Error()))
		return
	}

	return post_filedata_from_buffer(url, caCertFile, certFile, keyFile, filename, data)

}

///////////////////////////
// TODO 改成客户端，命令行带参数
func Client() *gin.Engine {
	fmt.Printf("post client start... %v\n", conf.Args)

	var s conf.BaseRequestqMsg
	var ss conf.MyInfo_t

	ss.Age = 250
	ss.Name = "latelee"

	s.Id = "id_250"
	s.Op = "test"
	s.Timestamp = 1233123211
	s.Data = ss
	jsonBytes, _ := json.Marshal(s)

	fmt.Println("jsonBytes: ", string(jsonBytes))

	var bakfile, body string
	var err error
	// jsonBytes = []byte("hello") // 发原始的

	//bakfile, body, err = post_from_buffer("https://127.0.0.1:9000/testing", "rootCA.crt", "latelee.crt", "latelee.key", jsonBytes)

	// bakfile, body, err = post_from_buffer("http://127.0.0.1:9000/testing", "rootCA.crt", "latelee.crt", "latelee.key", jsonBytes)

	bakfile, body, err = post_filedata_from_file("http://127.0.0.1:9000/testing", "", "", "", "foo.json")

	fmt.Printf("bakfile: [%v] body: [%v] err: [%v]\n", bakfile, body, err)

	return nil
}
