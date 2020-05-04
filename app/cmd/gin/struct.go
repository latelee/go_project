/*
本模块用到的结构体
*/

package gin

// 通用数据传输结构体
type BaseMessage struct {
    Id string `json:"id"`
    Op string `json:"op"`
    Timestamp int64 `json:"timestamp"`
    Data interface{} `json:"data"`
}

// 通用响应结构体
// 如用于http获取信息返回
type MessageResponse struct {
    Code int `json:"code"`
    Msg string `json:"msg"`
    Data interface{} `json:"data"`
}

// http客户端请求的
