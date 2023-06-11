package api

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"majiang/model"
	"majiang/service"
	"majiang/util"
	"sync"
)

var mutex sync.Mutex

// createRoom 创建房间
func createRoom(conn *websocket.Conn, UID float64) {
	log.Println("create")
	// 加锁
	mutex.Lock()
	defer mutex.Unlock()
	// 生成房间ID
	roomID := util.GenerateRoomId()
	// 创建房间
	room := &Room{
		ID:       roomID,
		Owner:    conn,
		Players:  make(map[int]*websocket.Conn),
		ReadyMap: make(map[int]bool),
		Users:    make(map[int]*model.PlayUser),
		IsFull:   false,
		IsStart:  false,
	}
	//查找用户信息
	user, err := service.SearchUserById(int(UID))
	if err != nil {
		broadcastInfo(conn, "查找用户信息失败")
		conn.Close()
		//从房间列表中移除该房间
		delete(Rooms, roomID)
		return
	}
	//加入玩家链接
	room.Players[int(UID)] = conn
	//设置准备状态
	room.ReadyMap[int(UID)] = false
	//添加用户
	room.Users[int(UID)] = &model.PlayUser{ID: int(UID), Username: user.Username}
	// 添加房间到房间列表中
	Rooms[roomID] = room
	log.Println(Rooms)
	// 广播房间信息
	broadcast(Create, room, user.Username, true)
}

// joinRoom 加入房间
func joinRoom(conn *websocket.Conn, roomID string, UID float64) {
	// 加锁
	mutex.Lock()
	defer mutex.Unlock()
	// 获取房间
	room, ok := Rooms[roomID]
	if !ok {
		broadcastInfo(conn, "获取房间失败")
		return
	}
	if room.IsFull == true {
		// 房间已满
		broadcastInfo(conn, "房间已满员")
		conn.Close()
		return
	}
	//查找用户信息
	user, err := service.SearchUserById(int(UID))
	if err != nil {
		broadcastInfo(conn, "查找用户信息失败")
		//用户信息查找发生错误
		conn.Close()
		return
	}
	// 添加玩家到房间中
	room.Players[int(UID)] = conn
	//满员需设置
	if len(room.Players) >= 4 {
		room.IsFull = true
	}
	//设置准备状态
	room.ReadyMap[int(UID)] = false
	//添加用户
	room.Users[int(UID)] = &model.PlayUser{ID: int(UID), Username: user.Username}
	// 广播房间信息
	broadcast(Join, room, user.Username, true)
	//if room.IsFull == true {
	//	broadcast(Start, room, user.Username, true)
	//	StartGame(room)
	//}
}

// leaveRoom 退出房间
func leaveRoom(conn *websocket.Conn, roomID string, UID float64) {
	log.Println(Rooms)
	log.Println("leave")

	// 获取房间
	room, ok := Rooms[roomID]
	if !ok {
		broadcastInfo(conn, "获取房间失败")
		return
	}
	// 加锁
	mutex.Lock()
	defer mutex.Unlock()
	if room.IsStart == true {
		broadcastInfo(conn, "游戏已开始")
		return
	}

	if room.IsFull == true {
		room.IsFull = false
	}

	if room.ReadyMap[int(UID)] == true {
		broadcastInfo(conn, "请先取消准备")
		return
	}
	// 移除玩家
	delete(room.Players, int(UID))
	delete(room.Users, int(UID))
	delete(room.ReadyMap, int(UID))
	// 如果房主退出房间，则随机一个玩家成为新的房主
	if conn == room.Owner {
		if len(room.Players) > 0 {
			for newOwnerIndex := range room.Players {
				room.Owner = room.Players[newOwnerIndex]
				break
			}
		} else {
			delete(Rooms, roomID)
			return
		}
	}
	//查找用户信息
	user, err := service.SearchUserById(int(UID))
	if err != nil {
		broadcastInfo(conn, "查找用户信息失败")
		return
	}
	conn.Close()
	// 广播房间信息
	broadcast(Leave, room, user.Username, true)
}

// changeReadyState 改变准备状态
func changeReadyState(conn *websocket.Conn, roomID string, ready bool, UID float64) {
	// 加锁
	mutex.Lock()
	defer mutex.Unlock()
	// 获取房间
	room, ok := Rooms[roomID]
	if !ok {
		broadcastInfo(conn, "获取房间失败")
		return
	}
	if room.IsStart == true {
		broadcastInfo(conn, "游戏已开始")
		return
	}
	// 改变准备状态
	if room.ReadyMap[int(UID)] == ready {
		broadcastInfo(conn, "重复更改准备状态")
		return
	}
	room.ReadyMap[int(UID)] = ready
	var count = 0
	//如果全部准备则直接开始
	for _, b := range room.ReadyMap {
		if b != true {
			break
		}
		count++
	}
	//查找用户信息
	user, err := service.SearchUserById(int(UID))
	if err != nil {
		broadcastInfo(conn, "查找用户信息失败")
		return
	}
	// 广播房间信息
	broadcast(ChangeReady, room, user.Username, ready)
	if count == 4 {
		room.IsStart = true
		broadcast(Common, room, "游戏开始！", true)
	}
}

func GetRoomList(conn *websocket.Conn) {
	var boolMap = map[bool]string{
		true:  "已开始",
		false: "未开始",
	}
	var sentences string
	for roomID, room := range Rooms {
		sentence := fmt.Sprintf("房间号:%s 状态:%s 人数:%d\n", roomID, boolMap[room.IsStart], len(room.Players))
		sentences += sentence
	}
	broadcastInfo(conn, sentences)
}

func chat(conn *websocket.Conn, UID float64, roomID string, sentence string) {
	// 加锁
	mutex.Lock()
	defer mutex.Unlock()
	//查找用户信息
	user, err := service.SearchUserById(int(UID))
	if err != nil {
		broadcastInfo(conn, "查找用户信息失败")
		return
	}
	// 获取房间
	room, ok := Rooms[roomID]
	if !ok {
		broadcastInfo(conn, "获取房间失败")
		return
	}
	sentence = fmt.Sprintf("%s:%s", user.Username, sentence)
	broadcast(Common, room, sentence, true)
}
