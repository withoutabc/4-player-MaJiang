package api

import (
	"github.com/gorilla/websocket"
	"math/rand"
	"sync"
	"time"
)

type Card struct {
	Suit  int //花色 1-3
	Point int //点数 1-9
}

type Player struct {
	turn      int    //玩家轮次
	Username  string //用户名
	Cards     []Card //手牌
	PengCards []Card //碰牌列表
	GangCards []Card //杠牌列表
	Conn      *websocket.Conn
}

type Game struct {
	ID          string          //房间id
	TurnMap     map[int]*Player //随机分配轮次
	Banker      int             //庄家(随机)
	Wall        []Card          //牌墙
	DiscardPile []Card          //弃牌区
	currentTurn int             //当前轮次
	Mutex       sync.Mutex
}

// InitWall 初始化牌墙
func InitWall() []Card {
	var wall []Card
	for i := 0; i < 4; i++ {
		for j := 1; j <= 9; j++ {
			wall = append(wall, Card{Suit: Character, Point: j})
			wall = append(wall, Card{Suit: Bamboo, Point: j})
			wall = append(wall, Card{Suit: Dot, Point: j})
		}
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(wall), func(i, j int) { wall[i], wall[j] = wall[j], wall[i] })
	return wall
}

func NewGame(room *Room) *Game {
	game := &Game{
		ID:          room.ID,
		TurnMap:     make(map[int]*Player),
		Wall:        InitWall(),
		currentTurn: 1,
	}

	players := make([]*Player, 0)
	for userID, user := range room.Users {
		player := &Player{
			Username: user.Username,
			Conn:     room.Players[userID],
			turn:     -1,
			Cards:    make([]Card, 0),
		}
		players = append(players, player)
	}

	// 随机分配轮次
	rand.Seed(time.Now().UnixNano())
	turns := rand.Perm(4) // 生成一个随机轮次切片
	for i, player := range players {
		player.turn = turns[i] + 1 //分配随机轮次
		game.TurnMap[player.turn] = player
	}
	//
	game.Banker = 1
	// 发牌
	for i := 0; i < 13; i++ {
		for _, player := range players {
			player.Cards = append(player.Cards, game.Wall[0]) // 给玩家发牌
			game.Wall = game.Wall[1:]                         // 牌墙去掉已经发出的牌
		}
	}
	for _, player := range players {
		if player.turn == game.Banker {
			player.Cards = append(player.Cards, game.Wall[0]) // 庄家多发一张牌
			game.Wall = game.Wall[1:]
			break
		}
	}
	return game
}
