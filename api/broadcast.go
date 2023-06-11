package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

// broadcastCard 全局看牌
func broadcastCard(game *Game) {
	for _, player := range game.TurnMap {
		//看手牌
		message := player.Cards
		jsonBytes, err := json.Marshal(message)
		if err != nil {
			log.Println(err)
			return
		}
		if err = player.Conn.WriteMessage(websocket.TextMessage, jsonBytes); err != nil {
			log.Println(err)
			return
		}
		//看碰牌
		message = player.PengCards
		jsonBytes, err = json.Marshal(message)
		if err != nil {
			log.Println(err)
			return
		}
		if err = player.Conn.WriteMessage(websocket.TextMessage, jsonBytes); err != nil {
			log.Println(err)
			return
		}
		//看杠牌
		message = player.GangCards
		jsonBytes, err = json.Marshal(message)
		if err != nil {
			log.Println(err)
			return
		}
		if err = player.Conn.WriteMessage(websocket.TextMessage, jsonBytes); err != nil {
			log.Println(err)
			return
		}
	}
}

// broadcast 全局发送信息
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
		case Start:
			message = "游戏开始"
		}
		// 发送消息
		err := player.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

// broadcastInfo 单独发送信息
func broadcastInfo(conn *websocket.Conn, message string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println(err)
		return
	}

}

// broadcastChuPai 全局出牌广播
func broadcastChuPai(game *Game, suit, point float64) {
	//发送给房间内每一位玩家
	message := fmt.Sprintf("%s打出%d,%d\n请在60秒内选择碰、杠、胡", game.TurnMap[game.currentTurn].Username, int(suit), int(point))
	for _, player := range game.TurnMap {
		// 发送消息
		err := player.Conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

// broadcastMoPai 全局摸牌广播
func broadcastMoPai(game *Game) {
	//发送给房间内每一位玩家
	message := fmt.Sprintf("进入新的回合，%s摸了一张牌", game.TurnMap[game.currentTurn].Username)
	for _, player := range game.TurnMap {
		// 发送消息
		err := player.Conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println(err)
			continue
		}
	}
	//摸完牌看下各自手牌
	broadcastCard(game)
}

func broadcastPeng(game *Game, turn, suit, point float64) {
	//发送给房间内每一位玩家
	message := fmt.Sprintf("%s碰：%d,%d", game.TurnMap[int(turn)].Username, int(suit), int(point))
	for _, player := range game.TurnMap {
		// 发送消息
		err := player.Conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println(err)
			continue
		}
	}

	//碰完牌看下各自手牌
	broadcastCard(game)

}

func broadcastGang(game *Game, suit, point float64) {
	//发送给房间内每一位玩家
	message := fmt.Sprintf("%s杠：%d,%d", game.TurnMap[game.currentTurn].Username, int(suit), int(point))
	for _, player := range game.TurnMap {
		// 发送消息
		err := player.Conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println(err)
			continue
		}
	}
	//杠完牌看下各自手牌
	broadcastCard(game)

}

func broadcastHu(game *Game, turn float64) {
	//发送给房间内每一位玩家
	message := fmt.Sprintf("%s杠胡了", game.TurnMap[int(turn)].Username)
	for _, player := range game.TurnMap {
		// 发送消息
		err := player.Conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println(err)
			continue
		}
	}
	//胡完牌看下各自手牌
	broadcastCard(game)
}
