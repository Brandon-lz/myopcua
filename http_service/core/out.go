// serialize.go is a Go file that contains the code to serialize data to JSON format.

package core

import (
    "net/http"

    "github.com/Brandon-lz/myopcua/utils"
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
)

type apiResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    string `json:"data"`
}

// SuccessHandler 返回成功响应
func SuccessHandler(c *gin.Context, responseData interface{}) {
    if responseData == nil {
        c.JSON(http.StatusOK, apiResponse{
            Code:    http.StatusOK,
            Message: "",
            Data:    "",
        })
    } else {
        c.JSON(http.StatusOK, responseData)
    }
}

func SerializeDataAndValidate[T interface{}](source T, target *T, doSerialize ...bool) T { // target必须为指针类型
    if len(doSerialize) > 0 && doSerialize[0] {
        return SerializeData(source, target)
    } else {
        
        return ValidateSchema(source)
    }
}

func SerializeData[T interface{}](source any, target *T) T { // target必须为指针类型
    return utils.SerializeData(source, target)
}

func ValidateSchema[T interface{}](source T)T{
    if err := ValidateStruct(source); err != nil {
        logrus.Error("数据校验失败:", err)
        panic(NewKnownError(http.StatusInternalServerError, nil, "output数据异常"))
    }
    return source
}