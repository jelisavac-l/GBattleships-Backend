package routes

import (
	"encoding/json"
	"net/http"

	"github.com/jelisavac-l/GBattleships/internal/game"
	"github.com/jelisavac-l/GBattleships/internal/gamehandler"
	"github.com/jelisavac-l/GBattleships/internal/model"
)

type newGameRequest struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
type newGameResponse struct {
	GID string `json:"gid"`
}
type availableGamesResponse struct {
	GameList []string `json:"gamelist"` // gameID:PlayerUsername is element
}

var availableGames = []*game.Game{}
var gameIdCounter = 0

func RegisterServerRoutes() {
	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getAvailableGames(w, r)
		case http.MethodPost:
			// Parse request data
			var data newGameRequest
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// Read player data and call createGame
			player := model.Player{ID: data.ID, Username: data.Username}
			createGame(w, r, &player)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

// POST /game
func createGame(w http.ResponseWriter, r *http.Request, player *model.Player) {
	gamePtr := game.CreateGame(*player, gameIdCounter)
	gameIdCounter++

	availableGames = append(availableGames, gamePtr)

	go gamehandler.Run(gamePtr)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newGameResponse{GID: gamePtr.ID})
}

// GET /game
func getAvailableGames(w http.ResponseWriter, r *http.Request) {
	gameList := []string{}
	for _, gamePtr := range availableGames {
		gameList = append(gameList, gamePtr.ID+":"+gamePtr.Player1.Username)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(availableGamesResponse{GameList: gameList})
}
