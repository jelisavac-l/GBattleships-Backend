package gamehandler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jelisavac-l/GBattleships/internal/game"
	"github.com/jelisavac-l/GBattleships/internal/model"
)

func Run(game *game.Game) {
	game.Player1 = &model.Player{}
	game.Player2 = &model.Player{}
	RegisterHandlerRoutes(*game)

	// game.StartGame()
	// testing block by evil gpt
	{
		conn := game.Player1.Conn
		go func() {
			defer conn.Close()
			for {
				mt, msg, err := conn.ReadMessage()
				if err != nil {
					log.Println("Player1 read error:", err)
					return
				}
				log.Printf("Player1 sent: %s", string(msg))

				// Echo back the message
				err = conn.WriteMessage(mt, msg)
				if err != nil {
					log.Println("Player1 write error:", err)
					return
				}
			}
		}()
		conn2 := game.Player2.Conn
		go func() {
			defer conn2.Close()
			for {
				mt, msg, err := conn2.ReadMessage()
				if err != nil {
					log.Println("Player1 read error:", err)
					return
				}
				log.Printf("Player1 sent: %s", string(msg))

				// Echo back the message
				err = conn2.WriteMessage(mt, msg)
				if err != nil {
					log.Println("Player1 write error:", err)
					return
				}
			}
		}()
	}

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func RegisterHandlerRoutes(g game.Game) {
	http.HandleFunc("/"+g.ID+"/player1", func(w http.ResponseWriter, r *http.Request) {
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

	log.Println("Starting server on: 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
