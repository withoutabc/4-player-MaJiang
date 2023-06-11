package model

import (
	"github.com/gorilla/websocket"
)

type ClientPool struct {
	Clients    map[*websocket.Conn]bool
	Broadcast  chan []byte
	AddClient  chan *websocket.Conn
	DelClient  chan *websocket.Conn
	PlayerList []string
	ReadyList  []bool
}
