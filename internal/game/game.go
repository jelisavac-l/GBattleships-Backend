package game

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jelisavac-l/GBattleships/internal/model"
)

type Game struct {
	ID      string
	Player1 *model.Player
	Player2 *model.Player
	Board1  *model.Board
	Board2  *model.Board
	Turn    string
	State   string
}

func CreateGame() *Game {
	return &Game{
		ID:    uuid.New().String(),
		State: "waiting",
	}
}

// to be removed maybe
func (game *Game) JoinGame(player *model.Player) error {
	if game.Player2 != nil {
		return errors.New(game.ID + " game is full!")
	} else {
		game.Player2 = player
	}
	return nil
}

func (game *Game) StartGame() {

	// Request boards

	// Loop moves & check conditions

}
