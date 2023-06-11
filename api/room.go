package api

import (
	"github.com/gorilla/websocket"
	"log"
	"majiang/model"
	"majiang/service"
	"majiang/util"
	"sync"
)

var mutex *sync.Mutex

// createRoom 创建房间
func createRoom(conn *websocket.Conn, UID int64) {
	log.Println("create")
	// 加锁
	//	mutex.Lock()
	//	defer mutex.Unlock()
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
		delete(rooms, roomID)
		return
	}
	//加入玩家链接
	room.Players[int(UID)] = conn
	//设置准备状态
	room.ReadyMap[int(UID)] = false
	//添加用户
	room.Users[int(UID)] = &model.PlayUser{ID: int(UID), Username: user.Username}
	// 添加房间到房间列表中
	rooms[roomID] = room
	// 广播房间信息
	broadcast(create, room, user.Username, true)
}

// joinRoom 加入房间
func joinRoom(conn *websocket.Conn, roomID string, UID int) {
	// 加锁
	mutex.Lock()
	defer mutex.Unlock()
	// 获取房间
	room, ok := rooms[roomID]
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
	user, err := service.SearchUserById(UID)
	if err != nil {
		broadcastInfo(conn, "查找用户信息失败")
		//用户信息查找发生错误
		conn.Close()
		return
	}
	// 添加玩家到房间中
	room.Players[UID] = conn
	//满员需设置
	if len(room.Players) >= 4 {
		room.IsFull = true
	}
	//设置准备状态
	room.ReadyMap[UID] = false
	//添加用户
	room.Users[UID] = &model.PlayUser{ID: UID, Username: user.Username}
	// 广播房间信息
	broadcast(join, room, user.Username, true)
}

// leaveRoom 退出房间
func leaveRoom(conn *websocket.Conn, roomID string, UID int) {
	// 加锁
	mutex.Lock()
	defer mutex.Unlock()
	// 获取房间
	room, ok := rooms[roomID]
	if !ok {
		broadcastInfo(conn, "获取房间失败")
		return
	}

	if room.IsStart == true {
		broadcastInfo(conn, "游戏已开始")
		return
	}

	if room.IsFull == true {
		room.IsFull = false
	}

	if room.ReadyMap[UID] == true {
		broadcastInfo(conn, "请先取消准备")
		return
	}
	// 移除玩家
	delete(room.Players, UID)
	delete(room.Users, UID)
	delete(room.ReadyMap, UID)
	// 如果房主退出房间，则随机一个玩家成为新的房主
	if conn == room.Owner {
		if len(room.Players) > 0 {
			for newOwnerIndex := range room.Players {
				room.Owner = room.Players[newOwnerIndex]
				break
			}
		} else {
			delete(rooms, roomID)
			return
		}
	}
	//查找用户信息
	user, err := service.SearchUserById(UID)
	if err != nil {
		broadcastInfo(conn, "查找用户信息失败")
		return
	}
	conn.Close()
	// 广播房间信息
	broadcast(leave, room, user.Username, true)
}

// changeReadyState 改变准备状态
func changeReadyState(conn *websocket.Conn, roomID string, ready bool, UID int) {
	// 加锁
	mutex.Lock()
	defer mutex.Unlock()
	// 获取房间
	room, ok := rooms[roomID]
	if !ok {
		broadcastInfo(conn, "获取房间失败")
		return
	}
	if room.IsStart == true {
		broadcastInfo(conn, "游戏已开始")
		return
	}
	// 改变准备状态
	if room.ReadyMap[UID] == ready {
		broadcastInfo(conn, "重复更改准备状态")
		return
	}
	room.ReadyMap[UID] = ready
	var count = 0
	//如果全部准备则直接开始
	for _, b := range room.ReadyMap {
		if b != true {
			break
		}
		count++
	}
	//查找用户信息
	user, err := service.SearchUserById(UID)
	if err != nil {
		broadcastInfo(conn, "查找用户信息失败")
		return
	}
	// 广播房间信息
	broadcast(changeReady, room, user.Username, ready)
	if count == 4 {
		room.IsStart = true
		broadcast(common, room, "游戏开始！", true)
	}
}