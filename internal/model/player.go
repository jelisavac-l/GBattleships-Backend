package model

import "github.com/gorilla/websocket"

type Player struct {
	ID       string
	Username string
	Conn     *websocket.Conn
}
