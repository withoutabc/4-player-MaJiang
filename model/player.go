package model

// Player 玩家信息
type Player struct {
	User    *User `json:"user" form:"user"`         // 玩家对应的用户信息
	Ready   bool  `json:"ready" form:"ready"`       // 玩家是否准备
	Score   int   `json:"score" form:"score"`       // 玩家分数
	IsOwner bool  `json:"is_owner" form:"is_owner"` // 玩家是否是房主
}
