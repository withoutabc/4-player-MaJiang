package api

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

func broadcast(broadType float64, room *Room, username string, b bool) {
	var message string
	for _, player := range room.Players {
		switch broadType {
		case Join:
			message = fmt.Sprintf("%s进入了房间", username)
		case Create:
			message = fmt.Sprintf("房间号：%s\n%s创建了房间", room.ID, username)
		case Leave:
			message = fmt.Sprintf("%s离开了房间", username)
		case ChangeReady:
			if b == true {
				message = fmt.Sprintf("%s已准备", username)
			} else {
				message = fmt.Sprintf("%s已取消准备", username)
			}
		case Common:
			message = username
		}
		// 发送消息
		err := player.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func broadcastInfo(conn *websocket.Conn, message string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println(err)
		return
	}
}
