package main

import (
	"fmt"
	"log"
	"net/http"

	"example.com/resistance/pkg/game"
	"example.com/resistance/pkg/websocket"
	"github.com/gin-gonic/gin"
)

func ws(c *gin.Context) {

}

func main() {
	router := gin.Default()
	pool := websocket.NewPool()
	gameState := game.NewGame(pool)
	go gameState.Run()
	router.GET("/ws", func(c *gin.Context) {
		conn, err := websocket.Upgrade(c)
		if err != nil {
			log.Fatalf("%s", err)
		}
		id := c.Query("id")
		if id == "" {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Did not find id param"))
			return
		}
		client := &websocket.Client{Conn: conn, Pool: pool, ID: id}
		pool.Register <- client
		client.Read()
	})

	router.Run(":8080")
}
