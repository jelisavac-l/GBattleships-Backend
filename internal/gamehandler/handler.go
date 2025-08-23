package gamehandler

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jelisavac-l/GBattleships/internal/game"
	"github.com/jelisavac-l/GBattleships/internal/model"
)

var wg sync.WaitGroup

func Run(game *game.Game) {
	game.Player2 = &model.Player{}

	wg.Add(2)
	RegisterHandlerRoutes(game)
	wg.Wait()

	fmt.Println("Game starting...")
	rematch := game.StartGame()
	fmt.Println(rematch)
	// fmt.Println("EXITED FOR LOOP ?!?!?!?")

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func RegisterHandlerRoutes(g *game.Game) {
	http.HandleFunc("/"+g.ID+"/player1", func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Failed to upgrade Player1:", err)
			return
		}
		g.Player1.ID = r.URL.Query().Get("id")
		g.Player1.Username = r.URL.Query().Get("username")
		g.Player1.Conn = conn
		log.Println("Player1 connected")
	})

	http.HandleFunc("/"+g.ID+"/player2", func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Failed to upgrade Player2:", err)
			return
		}
		g.Player2.ID = r.URL.Query().Get("id")
		g.Player2.Username = r.URL.Query().Get("username")
		g.Player2.Conn = conn
		log.Println("Player2 connected")
	})
}
