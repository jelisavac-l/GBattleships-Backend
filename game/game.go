package game

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jelisavac-l/GBattleships/models"
)

type Game struct {
	ID			string
	Player1		*models.Player
	Player2		*models.Player
	Board1		*models.Board
	Board2		*models.Board
	Turn		string
	State 		string
}

func CreateGame(player *models.Player) *Game {
	return &Game{
		ID: uuid.New().String(),
		Player1: player,
		State: "waiting",
	}
}

func (game *Game) JoinGame(player *models.Player) error {
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


