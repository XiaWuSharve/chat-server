package chat

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	id      string
	conn    *websocket.Conn
	send    chan any
	receive chan any
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	// 检查连接的Origin头
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewClient(ctx *gin.Context, id string, bufSize int) (*Client, error) {

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade http to websocket: %v", err)
	}

	c := &Client{
		id:      id,
		conn:    conn,
		send:    make(chan any, bufSize),
		receive: make(chan any, bufSize),
	}

	Register(c)

	// go c.Read()
	// go c.Write()

	return c, nil
}

func (c *Client) handleRead() {

}

func (c *Client) handleWrite() {
	for {
		select {
		// case m := <-c.send:

		}
	}
}

func (c *Client) Read() {
	// mess := <-c.receive

}

func (c *Client) Write(mess any) {
	c.send <- mess
}
