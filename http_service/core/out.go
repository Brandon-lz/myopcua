// serialize.go is a Go file that contains the code to serialize data to JSON format.

package core

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)


type apiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
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
	} else{
		c.JSON(http.StatusOK, responseData)
	}
	
}



func SerializeDataAndValidate[T interface{}](source interface{}, target *T) (T)  {       // target必须为指针类型
	jsonData, err := json.Marshal(source)
	if err != nil {
		logrus.Error("JSON序列化失败:", err)
		panic(NewKnownError(http.StatusInternalServerError, nil,"output数据异常"))
	}

	err = json.Unmarshal(jsonData, target)
	if err != nil {
		logrus.Error("JSON反序列化失败:", err)
		panic(NewKnownError(http.StatusInternalServerError, nil,"output数据异常"))
	}

	if err := ValidateStruct(target); err != nil {
		logrus.Error("数据校验失败:", err)
		panic(NewKnownError(http.StatusInternalServerError, nil,"output数据异常"))
    }

	return *target
}

