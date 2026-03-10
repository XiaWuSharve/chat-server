package mychat

import (
	"context"
	"kama_chat_server/chat"
	"kama_chat_server/dao"
	"kama_chat_server/dto"
	"kama_chat_server/helper"
	"kama_chat_server/kafka"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func TestKafka(t *testing.T) {
	kafka.Start()
	defer kafka.Close()
	ctx := context.Background()
	err := kafka.SendMessage(ctx, []byte("hello kafka"))
	if err != nil {
		t.Error(err)
	}
	data, err := kafka.ReadMessage(ctx)
	if err != nil {
		t.Error(err)
	}
	if string(data) != "hello kafka" {
		t.Errorf("expected %s, got %s", "hello kafka", string(data))
	}
}

func TestMysql(t *testing.T) {
	m := helper.ChatRequest2Message(&dto.ChatRequest{
		Content: "hello mysql",
	})
	if err := dao.Insert(m); err != nil {
		t.Error(err)
	}
}

func TestChat(t *testing.T) {
	kafka.Start()
	defer kafka.Close()
	chat.Start()
	defer chat.Close()
	// 1. 设置 gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 2. 创建 gin 引擎并定义 WebSocket 路由
	r := gin.Default()
	r.GET("/ws", func(c *gin.Context) {
		// 这里调用你的 NewClient 函数，传入 gin.Context
		// 需要生成一个唯一 id 和合适的缓冲区大小
		id := c.Query("id") // 从请求参数获取 id，或生成唯一 ID
		if id == "" {
			id = "test-client"
		}
		bufSize := 1024
		client, err := chat.NewClient(c, id, bufSize)
		if err != nil {
			// 升级失败时，gin 会处理响应，这里记录错误
			t.Logf("NewClient error: %v", err)
			return
		}
		chat.Register(client)
		// client 已经注册并启动了读写 goroutine，不需要额外操作
		// 注意：c.Writer 已经被升级为 WebSocket 连接，后续由 client 的读写循环接管
		// 此处理函数可以返回，但连接不会关闭，因为 client 持有 conn 并持续读写
	})

	// 3. 启动测试 HTTP 服务器
	server := httptest.NewServer(r)
	defer server.Close()

	// 将 HTTP URL 转换为 WebSocket URL
	wsURL := "ws" + server.URL[4:] + "/ws" // 将 http:// 替换为 ws://

	// 4. 创建两个 WebSocket 客户端连接（模拟前端）
	conn1, _, err := websocket.DefaultDialer.Dial(wsURL+"?id=client1", nil)
	if err != nil {
		t.Fatalf("Dial failed for client1: %v", err)
	}
	defer conn1.Close()

	conn2, _, err := websocket.DefaultDialer.Dial(wsURL+"?id=client2", nil)
	if err != nil {
		t.Fatalf("Dial failed for client2: %v", err)
	}
	defer conn2.Close()

	// 5. 发送消息并验证（示例）
	// 假设客户端发送 ChatRequest，服务端响应 ChatResponse
	// 你需要根据你的 dto 结构调整

	// 发送消息给 client1
	msg := &dto.ChatRequest{
		Content:   "hello chat",
		SendId:    "client1",
		ReceiveId: "client2",
	}
	err = conn1.WriteJSON(msg)
	if err != nil {
		t.Errorf("WriteJSON failed: %v", err)
	}

	// 读取来自 client1 的响应（如果有）
	var resp dto.ChatResponse
	err = conn2.ReadJSON(&resp)
	if err != nil {
		t.Errorf("ReadJSON failed: %v", err)
	}
	if resp.Content != msg.Content {
		t.Errorf("expected: %s, got: %s", msg.Content, resp.Content)
	}
	time.Sleep(5 * time.Second)
}
