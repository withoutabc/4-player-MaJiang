package client

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
)

var Conn *websocket.Conn

// InitWsSocket 连接WebSocket
func InitWsSocket() {
	u := url.URL{Scheme: "ws", Host: "localhost:2022", Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	Conn = conn
}

// WaitResponse 等待接收响应
func WaitResponse() {
	for {
		_, message, err := Conn.ReadMessage()
		if err != nil {
			fmt.Println("连接中断")
			return
		}
		fmt.Println(string(message))
	}
}

// SendMessage 发送消息
func SendMessage(messageType int, data map[string]interface{}) error {
	// 构造消息
	message := map[string]interface{}{
		"type": messageType,
	}
	for k, v := range data {
		message[k] = v
	}
	payload, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	// 发送消息
	err = Conn.WriteMessage(websocket.TextMessage, payload)
	if err != nil {
		fmt.Println("write:", err)
		return err
	}
	return nil
}
