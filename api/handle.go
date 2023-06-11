package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"majiang/model"
	"net/http"
	"sync"
	"time"
)

var Ticker chan struct{}

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
	mutex    sync.Mutex
}

// wsHandler 处理WebSocket连接
func wsHandler(c *gin.Context) {
	// 升级HTTP连接为WebSocket连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		log.Println("1")
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
			log.Println("2")
			break
		}
		// 处理消息
		handleMessage(conn, msg)
	}
	// 关闭连接
	conn.Close()
}

// handleMessage 处理消息
func handleMessage(conn *websocket.Conn, msg []byte) {
	//把每一个满员的房间开启
	go StartRoom()
	log.Println("handle message")
	// 解析消息
	var message map[string]interface{}
	err := json.Unmarshal(msg, &message)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(message["type"])
	// 获取消息类型
	//messageType, ok := message["type"].(int)
	// 根据消息类型调用对应的函数
	switch message["type"] {
	case Create:
		userID, ok := message["userID"].(float64)
		if !ok {
			log.Println("ok")
			return
		}
		createRoom(conn, userID)
	case Join:
		userID, ok := message["userID"].(float64)
		if !ok {
			return
		}
		roomID, ok := message["roomID"].(string)
		if !ok {
			return
		}
		joinRoom(conn, roomID, userID)
	case Leave:
		userID, ok := message["userID"].(float64)
		if !ok {
			log.Println("ds")
			return
		}
		roomID, ok := message["roomID"].(string)
		if !ok {
			log.Println("das")
			return
		}
		leaveRoom(conn, roomID, userID)
	case ChangeReady:
		userID, ok := message["userID"].(float64)
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
	case RoomList:
		GetRoomList(conn)
	case Common:
		userID, ok := message["userID"].(float64)
		if !ok {
			log.Println("n")
			return
		}
		roomID, ok := message["roomID"].(string)
		if !ok {
			log.Println("a")
			return
		}
		sentence, ok := message["sentence"].(string)
		if !ok {
			log.Println("m")
			return
		}
		chat(conn, userID, roomID, sentence)
	case ChuPai:
		userID, ok := message["userID"].(float64)
		if !ok {
			log.Println("n")
			return
		}
		roomID, ok := message["roomID"].(string)
		if !ok {
			log.Println("a")
			return
		}
		suit, ok := message["suit"].(float64)
		if !ok {
			log.Println("n")
			return
		}
		point, ok := message["point"].(float64)
		if !ok {
			log.Println("n")
			return
		}
		Chupai(conn, roomID, userID, suit, point)
	}

}

// StartRoom 初始化游戏
func StartRoom() {
	for _, room := range Rooms {
		if room.IsStart == true {
			game := NewGame(room)
			Games[room] = game
			//打印手牌
			log.Println(game.TurnMap[1].Cards)
			log.Println(game.TurnMap[2].Cards)
			log.Println(game.TurnMap[3].Cards)
			log.Println(game.TurnMap[4].Cards)
			broadcastCard(game) //广播各自手牌
			//开启协程，每一轮出牌过60秒后自动切换到下一个人轮次
			go Mopai(game)
		}
	}
}

// Mopai 循环轮次变更摸牌的处理
func Mopai(game *Game) {
	for {
		//判断一下还有没有牌了
		if len(game.Wall) == 0 {
			broadcast(Common, nil, "没牌了，游戏结束", true)
			return
		}
		if game.Count == 1 { //第一轮不摸直接出
			//提示玩家出牌
			broadcastInfo(game.TurnMap[game.currentTurn].Conn, "庄家请出一张牌")
		} else {
			//给当前轮次者发牌
			player := game.TurnMap[game.currentTurn]
			card := game.Wall[0]
			player.Cards = append(player.Cards, card) // 给玩家发牌
			game.Wall = game.Wall[1:]                 // 牌墙去掉已经发出的牌
			//全局广播
			broadcastMoPai(game)
			//提示玩家出牌
			broadcastInfo(game.TurnMap[game.currentTurn].Conn, "请出一张牌")
			//等待玩家出完牌
		}
		<-Ticker
		//等待60秒
		time.Sleep(time.Minute)
		//轮次变更
		game.currentTurn++
		if game.currentTurn > 4 {
			game.Count++
			game.currentTurn = 1
		}
	}
}
