package websocket

type Pool struct {
	Register   chan *Client
	Broadcast  chan Message
	Unregister chan *Client
	Clients    map[*Client]bool
	Action     chan Action
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		Action:     make(chan Action),
	}
}
