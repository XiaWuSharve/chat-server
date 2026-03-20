package controller

import (
	"kama_chat_server/chat"
	"kama_chat_server/dto"
	"kama_chat_server/enum"
	"kama_chat_server/helper"

	"github.com/gin-gonic/gin"
)

func GetMessageList(c *gin.Context) {
	var req dto.GetMessageListReqDto
	if err := c.ShouldBindQuery(&req); err != nil {
		helper.JsonBack(c, "", enum.RET_BAD_REQUEST, nil)
		return
	}
	msg, rsp, ret := chat.GetMessageList(c, req.UserOneId, req.UserTwoId)
	helper.JsonBack(c, msg, ret, rsp)
}
