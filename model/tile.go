package model

import (
	"math/rand"
	"time"
)

// TilePool 定义牌池结构体
type TilePool struct {
	Tiles []int // 牌池中所有牌的编号
}

// Shuffle 洗牌
func (pool *TilePool) Shuffle() {
	rand.Seed(time.Now().UnixNano()) // 使用当前时间作为随机数种子
	rand.Shuffle(len(pool.Tiles), func(i, j int) {
		pool.Tiles[i], pool.Tiles[j] = pool.Tiles[j], pool.Tiles[i]
	})
}

// Deal 发牌
func (pool *TilePool) Deal(playerNum int) [][]int {
	// 计算每个玩家应该分到多少张牌
	tileNum := 108 / playerNum
	handTiles := make([][]int, playerNum) // 保存每个玩家的手牌
	for i := 0; i < playerNum; i++ {
		handTiles[i] = make([]int, tileNum)
		for j := 0; j < tileNum; j++ {
			idx := i*tileNum + j
			handTiles[i][j] = pool.Tiles[idx]
		}
	}
	return handTiles
}
