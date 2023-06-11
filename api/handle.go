package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"majiang/model"
	"net/http"
	"reflect"
)

// 定义WebSocket连接
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Room 定义房间结构体
type Room struct {
	ID       string                  // 房间ID
	Owner    *websocket.Conn         // 房主
	Players  map[int]*websocket.Conn // 玩家
	ReadyMap map[int]bool            // 准备状态
	Users    map[int]*model.PlayUser // 用户
	IsFull   bool                    // 是否满员
	IsStart  bool                    // 是否开始
}

// 定义房间列表
var rooms = make(map[string]*Room)

// wsHandler 处理WebSocket连接
func wsHandler(c *gin.Context) {
	// 升级HTTP连接为WebSocket连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// 处理连接
	handleConnection(conn)
}

// 处理连接
func handleConnection(conn *websocket.Conn) {
	// 读取消息
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		// 处理消息
		handleMessage(conn, msg)
	}
	// 关闭连接
	conn.Close()
}

// 处理消息
func handleMessage(conn *websocket.Conn, msg []byte) {
	log.Println("handle message")
	// 解析消息
	var message map[string]interface{}
	err := json.Unmarshal(msg, &message)
	if err != nil {
		log.Println(err)
		return
	}
	// 获取消息类型
	//messageType, ok := message["type"].(int)
	// 根据消息类型调用对应的函数
	switch message["type"] {
	case create:
		log.Println("create")
		log.Println(message["userID"])

		log.Println(reflect.TypeOf(message["userID"]))
		createRoom(conn, 1)
	case join:
		userID, ok := message["userID"].(int)
		if !ok {
			return
		}
		roomID, ok := message["roomID"].(string)
		if !ok {
			return
		}
		joinRoom(conn, roomID, userID)
	case leave:
		userID, ok := message["userID"].(int)
		if !ok {
			return
		}
		roomID, ok := message["roomID"].(string)
		if !ok {
			return
		}
		leaveRoom(conn, roomID, userID)
	case changeReady:
		userID, ok := message["userID"].(int)
		if !ok {
			return
		}
		roomID, ok := message["roomID"].(string)
		if !ok {
			return
		}
		ready, ok := message["ready"].(bool)
		if !ok {
			return
		}
		changeReadyState(conn, roomID, ready, userID)
	}
}
