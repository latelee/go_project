/*
与设备打交道的请求
获取/设置工作状态接口

设置抓拍视频功能

获取/设置相机云台接口
获取/设置机械手臂接口
获取工作状态（即LED亮灭时间）
点灯
获取温湿度数据接口  湿度：% 温度 °C
获取北斗数据接口
*/

package gin

import (
    "fmt"
    //"time"
    "io/ioutil"
    "encoding/json"
    //"strconv"
    //"time"

    //"k8s.io/klog"

    "github.com/gin-gonic/gin"
)

// 获取类响应
// http://127.0.0.1:8080/device/v1/get/
func DeviceGetHandle(ctx *gin.Context) {
    data, _ := ioutil.ReadAll(ctx.Request.Body)
    //fmt.Printf("ctx.Request.body: %v\n", string(data))
    var v map[string]interface{}
    err := json.Unmarshal(data, &v)
    if err != nil {
        responseFailedMsg(-1, "json parse error", ctx)
        return
    }
    fmt.Printf("got get msg: %v\n", v)

    if v["op"] == nil {
        fmt.Println("op is nil")
        responseFailedMsg(-1, "op is nil, do nothing", ctx)
        return
    }
    switch v["op"] {
    case "getarm":
        fmt.Println("get arm")
        result := gin.H {
        "left": "on",
        "right":  "off",
        }
        responseOK(result, ctx)

    case "getptz":
        fmt.Println("getptz")
        result := "auto"
        responseOK(result, ctx)
    
    case "gettemphum":
        fmt.Println("get temperature and humidity")
        result := gin.H {
        "temp": 35.6,
        "hum":  65,
        }
        responseOK(result, ctx)
        
    default:
        fmt.Println("op ", v["op"], " not support")
        responseFailedMsg(-1, "op " + v["op"].(string) + " not support", ctx)
    }
    
    
}

// 设置类响应
// http://127.0.0.1:8080/device/v1/set/
func DeviceSetHandle(ctx *gin.Context) {
    data, _ := ioutil.ReadAll(ctx.Request.Body)
    //fmt.Printf("ctx.Request.body: %v\n", string(data))
    var v map[string]interface{}
    err := json.Unmarshal(data, &v)
    if err != nil {
        responseFailedMsg(-1, "json parse error", ctx)
        return
    }
    fmt.Printf("got set msg: %v\n", v)

    if v["op"] == nil {
        fmt.Println("op is nil")
        responseFailedMsg(-1, "op is nil, do nothing", ctx)
        return
    }
    switch v["op"] {
    case "setarm":
        fmt.Println("set arm")
        var vv map[string]interface{}
        // 需转2次
        tmpJson, _ := json.Marshal(v["data"])
        err := json.Unmarshal(tmpJson, &vv)
        if err != nil {
            responseFailedMsg(-1, "json parse error", ctx)
            return
        }
        fmt.Println("left: ", vv["left"], "right: ", vv["right"])
        
        responseOKEmpty(ctx)

    case "setptz":
        fmt.Println("setptz")
        fmt.Println("type: ", v["data"])
        // TODO：发消息到板子上，并判断结果
        responseOKEmpty(ctx)
        
    case "ledon":
        fmt.Println("ledon")
        fmt.Println("type: ", v["data"])
        // TODO：发消息到板子上，并判断结果
        responseOKEmpty(ctx)

    case "ledoff":
        fmt.Println("ledoff")
        fmt.Println("type: ", v["data"])
        // TODO：发消息到板子上，并判断结果
        responseOKEmpty(ctx)

    case "snap":
        fmt.Println("snap")
        // TODO：发消息到板子上，并判断结果
        responseOKEmpty(ctx)

    default:
        fmt.Println("op ", v["op"], " not support")
        responseFailedMsg(-1, "op " + v["op"].(string) + " not support", ctx)
    }
    
    
}