package gin

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func response(code int, msg string, data interface{}, ctx *gin.Context)  {
    r := MessageResponse{code, msg, data}
    ctx.JSON(http.StatusOK, r)
}

func responseMsg(msg string, data interface{}, ctx *gin.Context)  {
    response(0, msg, data, ctx)
}

func responseOK(data interface{}, ctx *gin.Context)  {
    response(0, "ok", data, ctx)
}

func responseOKEmpty(ctx *gin.Context)  {
    response(0, "ok", "", ctx)
}

// 失败为空数据体
func responseFailed(code int, ctx *gin.Context)  {
    response(code, "failed", "", ctx)
}

func responseFailedMsg(code int, msg string, ctx *gin.Context)  {
    response(code, msg, "", ctx)
}

