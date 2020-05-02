package gin

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)

func response(code int, msg string, data interface{}, ctx *gin.Context)  {
    r := MessageResponse{code, msg, time.Now().UnixNano() / 1e6, data}
    ctx.JSON(http.StatusOK, r)
}

func responseMsg(msg string, data interface{}, ctx *gin.Context)  {
    response(0, msg, data, ctx)
}

func responseOK(data interface{}, ctx *gin.Context)  {
    response(0, "request ok", data, ctx)
}

// 失败为空数据体
func responseFailed(code int, ctx *gin.Context)  {
    response(code, "requeset failed", nil, ctx)
}

func responseFailedMsg(code int, msg string, ctx *gin.Context)  {
    response(code, msg, nil, ctx)
}

