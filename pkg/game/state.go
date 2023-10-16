package game

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"example.com/resistance/pkg/websocket"
)

type Game struct {
	Players []Player
	Rounds  map[int]bool
	Round   int
	Pool    *websocket.Pool
}

func (game *Game) Run() {
	for {
		select {
		case client := <-game.Pool.Register:
			log.Printf("%s joined", client.ID)
			game.Pool.Clients[client] = true
			player := Player{
				ID:     client.ID,
				Name:   client.ID,
				Client: *client,
			}
			msg, ok := game.AddPlayer(player)
			if !ok {
				log.Print("Failed to add player: %s", msg)
				client.Conn.WriteJSON(websocket.Message{Msg: "Can't join. Game already started"})
				break
			}
			for poolClient, _ := range game.Pool.Clients {
				fmt.Println(client)
				poolClient.Conn.WriteJSON(websocket.Message{Msg: fmt.Sprintf("User %s joined", client.ID), User: "system"})
			}
			break
		case client := <-game.Pool.Unregister:
			log.Printf("%s disconnected", client.ID)
			delete(game.Pool.Clients, client)
			for poolClient, _ := range game.Pool.Clients {
				poolClient.Conn.WriteJSON(websocket.Message{Msg: fmt.Sprintf("User %s left", client.ID), User: "system"})
			}
			break
		case message := <-game.Pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client, _ := range game.Pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
			break
		case action := <-game.Pool.Action:
			log.Printf("Action recieved %s from %s", action.Action, action.Client.ID)
			_split := strings.Split(action.Action, " ")
			if _split[0] == "start" {
				game.Start()
			}
			break
		}
	}
}
func NewGame(pool *websocket.Pool) Game {
	return Game{
		Players: make([]Player, 0),
		Rounds:  make(map[int]bool),
		Round:   0,
		Pool:    pool,
	}
}
func (game *Game) AddPlayer(player Player) (string, bool) {
	if game.Round == 0 {
		game.Players = append(game.Players, player)
		return "Added", true
	}
	return "Game already started", false
}

// shuffle shuffles the elements of an array in place
func shuffle(array []Role) {
	for i := range array { //run the loop till the range of array
		j := rand.Intn(i + 1)                   //choose any random number
		array[i], array[j] = array[j], array[i] //swap the random element with current element
	}
}

func (game *Game) Start() bool {
	game.Round = 1
	roles := make([]Role, len(game.Players))

	for idx, _ := range roles {
		if idx < len(roles)/2 {
			roles[idx] = Spy
		} else {
			roles[idx] = Resistance
		}
	}
	shuffle(roles)
	for idx, player := range game.Players {
		player.Role = roles[idx]
		player.Client.Conn.WriteJSON(websocket.Message{Msg: fmt.Sprintf("You are %s", roles[idx])})
	}
	return true
}
