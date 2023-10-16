package game

import "example.com/resistance/pkg/websocket"

type Role string

type Player struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Role   Role   `json:"role"`
	Client websocket.Client
}

const (
	Resistance Role = "Resistance"
	Spy        Role = "Spy"
)
