package models

import "golang.org/x/net/websocket"

type Player struct {
	ID			string
	Username 	string
	Conn 		websocket.Conn
}
