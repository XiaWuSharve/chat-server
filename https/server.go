package https

import (
	"kama_chat_server/chat/controller"
	"kama_chat_server/jwt"

	"github.com/gin-gonic/gin"
)

var ge = gin.Default()

func SetupRouter() *gin.Engine {
	api := ge.Group("/api")
	api.Use(jwt.Middleware)
	api.GET("/message", controller.GetMessageList)
	api.GET("/message/group", controller.GetGroupMessageList)
	return ge
}

func Start(addr string) error {
	SetupRouter()
	return ge.Run(addr)
}
