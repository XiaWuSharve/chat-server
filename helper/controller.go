package helper

import (
	"kama_chat_server/enum"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ControllerJson struct {
	Message string `json:"message"`
	ret     int    `json:"ret"`
	Data    any    `json:"data"`
}

func JsonBack(ctx *gin.Context, message string, ret int, json any) {
	switch ret {
	case enum.RET_BAD_REQUEST:
		if message == "" {
			message = "请求数据有误"
		}
		ctx.JSON(http.StatusBadRequest, &ControllerJson{
			Message: message,
			ret:     ret,
			Data:    json,
		})
	case enum.RET_OK:
		ctx.JSON(http.StatusOK, &ControllerJson{
			Message: message,
			ret:     ret,
			Data:    json,
		})
	case enum.RET_SYSTEM_ERR:
		ctx.JSON(http.StatusOK, &ControllerJson{
			Message: message,
			ret:     ret,
			Data:    json,
		})
	}
}
