package gamehandler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jelisavac-l/GBattleships/internal/game"
	"github.com/jelisavac-l/GBattleships/internal/model"
)

func Run(game *game.Game) {

	game.Wg.Add(2)
	RegisterHandlerRoutes(game)
	game.Wg.Wait()

	rematch := true
	for rematch {
		log.Println("Game handler " + game.ID + " started...")
		rematch = safeStartGame(game)
	}

	log.Println("Game handler" + game.ID + " stopped.")
	game.Player1.Conn.Close()
	game.Player2.Conn.Close()
}

func safeStartGame(g *game.Game) (rematch bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Game %s crashed: %v", g.ID, r)
			if g.Player1 != nil && g.Player1.Conn != nil {
				game.SendErrorMessage("Game stopped due to error", g.Player1)

			}
			if g.Player2 != nil && g.Player2.Conn != nil {
				game.SendErrorMessage("Game stopped due to error", g.Player2)
			}
			rematch = false
			g.State = "finished"
		}
	}()
	rematch = g.StartGame()
	return rematch

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// Tell the upgrader to allow all incoming origins (including localhost)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RegisterHandlerRoutes(g *game.Game) {
	http.HandleFunc("/"+g.ID+"/player1", func(w http.ResponseWriter, r *http.Request) {
		handlePlayerConnection(&g.Player1, g, w, r)
	})

	http.HandleFunc("/"+g.ID+"/player2", func(w http.ResponseWriter, r *http.Request) {
		handlePlayerConnection(&g.Player2, g, w, r)
	})
}

func handlePlayerConnection(player **model.Player, g *game.Game, w http.ResponseWriter, r *http.Request) {
	// seat taken type beat
	if *player != nil && (*player).Conn != nil {
		http.Error(w, "Player already connected", http.StatusConflict)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade:", err)
		return
	}

	if *player == nil {
		*player = &model.Player{}
	}

	(*player).ID = r.URL.Query().Get("id")
	(*player).Username = r.URL.Query().Get("username")
	(*player).Conn = conn
	log.Printf("Player %s connected", (*player).Username)

	g.Wg.Done()
}
