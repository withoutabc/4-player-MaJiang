package api

const (
	Common float64 = iota
	Create
	Join
	Leave
	ChangeReady
	RoomList
	Start
	LookCard

	Play
	Peng
	Gang
	Hu
)

const (
	Character = 1 // 万
	Bamboo    = 2 // 条
	Dot       = 3 // 筒
)

const (
	MoPai float64 = 100 + iota
	ChuPai
)
