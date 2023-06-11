package api

import (
	"github.com/gorilla/websocket"
)

var Games = make(map[*Room]*Game)

// Chupai 出牌
func Chupai(conn *websocket.Conn, roomID string, userID, suit, point float64) {
	// 获取房间
	room, ok := Rooms[roomID]
	if !ok {
		broadcastInfo(conn, "获取房间失败")
		return
	}
	//获取游戏
	game, ok := Games[room]
	if !ok {
		broadcastInfo(conn, "获取游戏失败")
		return
	}
	//手牌区消除这张牌
	for i, card := range game.TurnMap[game.currentTurn].Cards {
		if int(suit) == card.Suit && int(point) == card.Point {
			cards := game.TurnMap[game.currentTurn].Cards
			cards = append(cards[:i], cards[i+1:]...)
		}
	}
	//把出的牌放到弃牌堆
	game.DiscardPile = append(game.DiscardPile, Card{
		Suit:  int(suit),
		Point: int(point),
	})
	//全局广播
	broadcastChuPai(game, suit, point)
	game.Ticker <- struct{}{} //开始60秒计时
}

// DoPeng 碰
func DoPeng(conn *websocket.Conn, roomID string, turn, suit, point float64) {
	// 获取房间
	room, ok := Rooms[roomID]
	if !ok {
		broadcastInfo(conn, "获取房间失败")
		return
	}
	//获取游戏
	game, ok := Games[room]
	if !ok {
		broadcastInfo(conn, "获取游戏失败")
		return
	}
	//判断不能自己碰自己
	if int(turn) == game.currentTurn {
		broadcastInfo(conn, "不能碰自己哦")
	}
	//判断碰的牌在他手上有没有2个
	targetCard := Card{Point: int(point), Suit: int(suit)}
	var count = 0
	var discard []int
	pengCards := game.TurnMap[int(turn)].PengCards
	//遍历，看看有几个
	for i, card := range game.TurnMap[int(turn)].Cards {
		if card == targetCard {
			count++
		}
		discard = append(discard, i)
	}
	//大于2就从手牌放2个到碰牌列表
	if count >= 2 {
		for _, cardIndex := range discard {
			pengCards = append(pengCards, game.TurnMap[int(turn)].Cards[cardIndex])
		}
		//放一个刚才摸的
		pengCards = append(pengCards, targetCard)
		game.TurnMap[int(turn)].PengCards = pengCards
		//手牌区删除这2张
		var count1 = 0
		var res []Card
		for i := 0; i < len(game.TurnMap[int(turn)].Cards); i++ {
			if count1 < 2 && i == discard[count] {
				count1++
				continue // 跳过要删除的索引
			}
			res = append(res, game.TurnMap[int(turn)].Cards[i]) // 将其他位置的元素添加到新的切片中
			if count1 == 2 {                                    // 计数器达到2，结束循环
				break
			}
		}
		game.TurnMap[int(turn)].Cards = res
	} else {
		broadcastInfo(conn, "这张牌不能碰哦")
		return
	}
	//全局广播
	broadcastPeng(game, turn, suit, point)
	//让这个人出牌
	//改变轮次
	game.currentTurn = int(turn)
}

// DoGang 杠
func DoGang(conn *websocket.Conn, roomID string, turn, suit, point float64) {
	// 获取房间
	room, ok := Rooms[roomID]
	if !ok {
		broadcastInfo(conn, "获取房间失败")
		return
	}
	//获取游戏
	game, ok := Games[room]
	if !ok {
		broadcastInfo(conn, "获取游戏失败")
		return
	}
	//判断不能别人杠自己
	if int(turn) != game.currentTurn {
		broadcastInfo(conn, "不能杠别人哦")
	}
	//判断碰的牌在他手上有没有3个
	targetCard := Card{Point: int(point), Suit: int(suit)}
	var count = 0
	var discard []int
	gangCards := game.TurnMap[int(turn)].GangCards
	//遍历，看看有几个
	for i, card := range game.TurnMap[int(turn)].Cards {
		if card == targetCard {
			count++
		}
		discard = append(discard, i)
	}
	//大于3就从手牌放3个到杠牌列表
	if count >= 3 {
		for _, cardIndex := range discard {
			gangCards = append(gangCards, game.TurnMap[int(turn)].Cards[cardIndex])
		}
		//放一个刚才摸的
		gangCards = append(gangCards, targetCard)
		game.TurnMap[int(turn)].PengCards = gangCards
		//手牌区删除这2张
		var count1 = 0
		var res []Card
		for i := 0; i < len(game.TurnMap[int(turn)].Cards); i++ {
			if count1 < 3 && i == discard[count] {
				count1++
				continue // 跳过要删除的索引
			}
			res = append(res, game.TurnMap[int(turn)].Cards[i]) // 将其他位置的元素添加到新的切片中
			if count1 == 3 {                                    // 计数器达到 3，结束循环
				break
			}
		}
		game.TurnMap[int(turn)].Cards = res
	} else {
		broadcastInfo(conn, "这张牌不能杠哦")
		return
	}
	//全局广播
	broadcastGang(game, suit, point)
	//让这个人摸牌
	//TODO:给他再摸一次牌
}

// DoHu 胡
func DoHu(conn *websocket.Conn, roomID string, turn, suit, point float64) {
	// 获取房间
	room, ok := Rooms[roomID]
	if !ok {
		broadcastInfo(conn, "获取房间失败")
		return
	}
	//获取游戏
	game, ok := Games[room]
	if !ok {
		broadcastInfo(conn, "获取游戏失败")
		return
	}
	var newCards []Card
	//组合所有牌
	cards := game.TurnMap[int(turn)].Cards
	pendCards := game.TurnMap[int(turn)].PengCards
	gangCards := game.TurnMap[int(turn)].GangCards
	newCard := Card{Suit: int(suit), Point: int(point)}
	newCards = append(append(newCards, cards...), pendCards...)
	newCards = append(append(newCards, gangCards...), newCard)
	if JudgeHu(newCards) {
		broadcastHu(game, turn)
	} else {
		broadcastInfo(conn, "不能胡")
	}
}

// JudgeHu 判断能不能胡
func JudgeHu(card []Card) bool {
	//TODO:判断
	return false
}
