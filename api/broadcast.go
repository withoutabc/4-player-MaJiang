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
		case join:
			message = fmt.Sprintf("%s进入了房间", username)
		case create:
			message = fmt.Sprintf("%s创建了房间", username)
		case leave:
			message = fmt.Sprintf("%s离开了房间", username)
		case changeReady:
			if b == true {
				message = fmt.Sprintf("%s已准备", username)
			} else {
				message = fmt.Sprintf("%s已取消准备", username)
			}
		case common:
			message = username
			// 发送消息
			err := player.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println(err)
				continue
			}

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
