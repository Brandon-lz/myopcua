package httpservice

import (
	"earth/http_service/core"

	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
)


type HttpResponse struct {
    Message string      `json:"msg"`
    Status  int         `json:"code"`
    Data    interface{} `json:"data"`
}

func ErrorHandler(c *gin.Context, err any) {
    var httpResponse HttpResponse
    switch v := err.(type) {
    case *core.KnownError: // 自定义异常
        httpResponse = HttpResponse{Status: v.Code, Data: v.Data, Message: v.Msg}
    default: // 系统异常
        goErr := errors.Wrap(err, 2)
        httpResponse = HttpResponse{Message: "Internal server error", Status: 500, Data: goErr.Error()}
    }
    c.AbortWithStatusJSON(500, httpResponse) // 有些公司用200
}

