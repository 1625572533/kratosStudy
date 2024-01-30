package ginrouter

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Data struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}, code int32, msg string) {
	c.Status(http.StatusOK)
	resp := &Data{
		Code: code,
		Data: data,
		Msg:  msg,
	}
	c.JSON(http.StatusOK, resp)
	return
}

func AbortWithError(c *gin.Context, httpCode int, err error) {
	if err != nil {
		c.AbortWithStatus(httpCode)
		c.JSON(httpCode, &Data{
			Code: int32(httpCode),
			Msg:  err.Error(),
			Data: nil,
		})
	}
}
