package controller

import (
	"fmt"
	"kama_chat_server/chat"
	"kama_chat_server/helper"

	"github.com/gin-gonic/gin"
)

func WssLogin(c *gin.Context) {
	userId, _ := c.Get("user_id")
	_, err := chat.NewClient(c, userId.(string), 1024)
	if err != nil {
		helper.JsonBack(c, "无法建立通讯", -1, fmt.Errorf("failed to create websocket client: %v", err))
		return
	}
	helper.JsonBack(c, "建立通讯成功", 0, nil)
}

func WssLogout(c *gin.Context) {
	userId, _ := c.Get("user_id")
	chat.Unregister(userId.(string))
	helper.JsonBack(c, "断开通讯成功", 0, nil)
}
