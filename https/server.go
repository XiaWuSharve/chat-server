package https

import (
	"kama_chat_server/chat/controller"

	"github.com/gin-gonic/gin"
)

var ge = gin.Default()

func SetupRouter() *gin.Engine {
	ge.GET("/message", controller.GetMessageList)
	return ge
}

func Start(addr string) error {
	SetupRouter()
	return ge.Run(addr)
}
