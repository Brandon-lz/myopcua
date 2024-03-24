package httpservice

import (
	"earth/http_service/core"

	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
)


type HttpResponse struct {
    Message string      `json:"msg"`
    Code  int         `json:"code"`
    Data    interface{} `json:"data"`
}

func ErrorHandler(c *gin.Context, err any) {
    var httpResponse HttpResponse
    switch v := err.(type) {
    case *core.KnownError: // 自定义异常
        httpResponse = HttpResponse{Code: v.Code, Data: v.Data, Message: v.Msg}
    default: // 系统异常
        goErr := errors.Wrap(err, 2)
        httpResponse = HttpResponse{Message: "Internal server error", Code: 500, Data: goErr.Error()}
    }
    c.AbortWithStatusJSON(500, httpResponse) // 有些公司用200
}

