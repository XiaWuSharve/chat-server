package chat

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	id   string
	conn *websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	// 检查连接的Origin头
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewClient(ctx *gin.Context, id string) (*Client, error) {

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade http to websocket: %v", err)
	}

	c := &Client{
		id:   id,
		conn: conn,
	}

	Register(c)

	// go c.Read()
	// go c.Write()

	return c, nil
}

func (c *Client) Read() {

}

func (c *Client) Write(mess any) {

}
