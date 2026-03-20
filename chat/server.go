package chat

import (
	"context"
	"kama_chat_server/dao"
	"kama_chat_server/helper"
	"kama_chat_server/zlog"
	"sync"
)

var (
	// 重构为channel + 原生map
	clients     sync.Map
	ctx, cancel = context.WithCancel(context.Background())
)

func Start() {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				req, err := KafkaReceiveMessage(ctx)
				if err != nil {
					zlog.Error("failed to receive kafka message: %v", err)
					continue
				}
				// write to db
				m := helper.ChatRequest2Message(req)
				if err := dao.Insert(m); err != nil {
					zlog.Error("failed to insert into db: %v", err)
					continue
				}
				// send by websocket
				res := helper.Message2ChatResponse(m)
				go func() {
					c, ok := clients.Load(res.ReceiveId)
					if ok {
						c.(*Client).Write(res)
					}
				}()
				// write to redis
				go func() {
					if err := RedisAddPrivateMessage(ctx, res); err != nil {
						zlog.Error("failed to add a private message to redis: %v", err)
						return
					}
				}()
				// TODO: caching reading message and batch store to db
			}
		}
	}()
}

func Close() {
	clients.Range(func(key, value any) bool {
		Unregister(value.(*Client))
		return true
	})
	cancel()
}

func Register(c *Client) {
	clients.Store(c.id, c)
}

func Unregister(c *Client) {
	c.conn.Close()
	clients.Delete(c.id)
}
