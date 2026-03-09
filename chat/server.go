package chat

import (
	"context"
	"kama_chat_server/dao"
	"kama_chat_server/helper"
	"kama_chat_server/zlog"
	"sync"
)

var (
	mu      sync.RWMutex
	clients map[string]*Client
)

func Start() {
	go func() {
		ctx := context.Background()
		for {
			req, err := ReceiveMessage(ctx)
			if err != nil {
				zlog.Fatal(err)
			}
			// write to db
			m := helper.KafkaRequest2Message(req)
			if err := dao.Insert(m); err != nil {
				zlog.Fatal(err)
			}
			go func() {
				mu.RLock()
				defer mu.RUnlock()
				if c, ok := clients[req.ReceiveId]; ok {
					c.Write(req.Content)
				}
			}()
			// write to redis
			go func ()  {
				
			}
			// TODO: caching reading message and batch store to db

		}
	}()
}

func Register(c *Client) {
	mu.Lock()
	defer mu.Unlock()
	clients[c.id] = c
}
