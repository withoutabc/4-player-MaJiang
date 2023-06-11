package model

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Room struct {
	ID         int
	MaxPlayers int
	NumPlayers int
	IsFull     bool
	IsPlaying  bool
	Players    []*websocket.Conn
	Mutex      *sync.Mutex `form:"mutex"`
}
