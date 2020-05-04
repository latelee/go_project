/*
本文返回json格式
获取设备信息接口
*/

package gin

import (
    // "fmt"
    "strconv"
    //"time"

    //"k8s.io/klog"

    "github.com/gin-gonic/gin"
)

// 设备信息
type DeviceInfo struct {
    Id      string
    Ip      string
    Mac     string
    Location string
    Version string
}

// 获取所有设备列表
// http://127.0.0.1:8080/device/v1/devlist/
func GetAllDevices(ctx *gin.Context) {

    devInfo := make([]DeviceInfo, 3)
    
    for i := 0; i < len(devInfo); i++ {
        devInfo[i].Id = "test00" + strconv.Itoa(i+1)
        devInfo[i].Ip = "192.168.0." + strconv.Itoa(i+1)
        devInfo[i].Mac = "08:00:27:81:48:b" + strconv.Itoa(i+1)
        devInfo[i].Location = "Shenzhen"
        devInfo[i].Version = "v0.1"
    }
/*
    result := gin.H {
        "result": devInfo,
        "count":  3,
    }
*/
    responseOK(devInfo, ctx)
}

// 获取指定设备列表
// http://127.0.0.1:8080/device/v1/devlist/1
// http://127.0.0.1:8080/device/v1/devlist/250
func GetSingleDevice(ctx *gin.Context) {
    id := ctx.Param("id")

    var devInfo DeviceInfo
    
    // 理论上是查询，此处从简，直接赋值，根据ID区别
    devInfo.Id = id
    devInfo.Ip = "192.168.0.1"
    devInfo.Mac = "08:00:27:81:48:ba"
    devInfo.Location = "Shenzhen"
    devInfo.Version = "v0.1"

    //  测试不存在的设备返回信息
    if id == "250" {
        //responseFailed(-1, ctx)
        responseFailedMsg(-1, "No such device", ctx)
        return
    }
/*
    result := gin.H {
        "result": devInfo,
        "count":  1,
    }
*/

    responseOK(devInfo, ctx)
}