package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"kama_chat_server/dto"
	"kama_chat_server/zlog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	id      string
	conn    *websocket.Conn
	send    chan *dto.ChatResponse
	receive chan *dto.ChatRequest
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
		send:    make(chan *dto.ChatResponse, bufSize),
		receive: make(chan *dto.ChatRequest, bufSize),
	}

	Register(c)

	go func() {
		if err := c.handleWrite(ctx); err != nil {
			zlog.Error(err)
			Unregister(c)
		}
	}()

	go func() {
		if err := c.handleRead(ctx); err != nil {
			zlog.Error(err)
			Unregister(c)
		}
	}()

	return c, nil
}

func (c *Client) handleRead(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, data, err := c.conn.ReadMessage()
			if err != nil {
				return fmt.Errorf("failed to read from websocket: %v", err)
			}
			var res dto.ChatRequest
			if err := json.Unmarshal(data, &res); err != nil {
				zlog.Error(fmt.Errorf("failed to parse data: %v", err))
				continue
			}
			KafkaSendMessage(ctx, &res)
		}
	}
}

func (c *Client) handleWrite(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case m := <-c.send:
			data, err := json.Marshal(m)
			if err != nil {
				return fmt.Errorf("failed to marshal message: %v", err)
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return fmt.Errorf("failed to write message to websocket: %v", err)
			}
		}
	}
}

func (c *Client) Write(mess *dto.ChatResponse) {
	c.send <- mess
}
