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
	Turn    string // what this for
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

func (game *Game) CheckValidBoard(board model.Board) {
	// checking if board is valid before setting game.Board
}

func (game *Game) checkValidMove(x int, y int, boardNo int) { //boardNo being 1 or 2
	// checks if move is valid (if (x,y) cell is iota or ship)
}

func (game *Game) PlayMove(x int, y int) {
	// checkvalidMove

	// board.ShootCell()
}
