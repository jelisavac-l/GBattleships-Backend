package game

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/jelisavac-l/GBattleships/internal/model"
)

type Game struct {
	ID      string
	Player1 *model.Player
	Player2 *model.Player
	Board1  *model.Board
	Board2  *model.Board
	Turn    bool // true for player1 	false for player2
	State   string
}

func CreateGame() *Game {
	return &Game{
		ID:    uuid.New().String(),
		State: "waiting",
	}
}

func (game *Game) StartGame() bool {
	var wg sync.WaitGroup
	wg.Add(2)
	go game.getBoard(game.Player1, &wg)
	go game.getBoard(game.Player2, &wg)
	wg.Wait()

	game.State = "playing"
	game.tellGameStarted()

	var hit bool
	var err error
	var previousX, previousY int
	for game.State != "finished" {
		x, y := game.getMove(game.Turn, previousX, previousY)
		if game.checkValidMove(x, y, !game.Turn) {
			hit, err = game.PlayMove(x, y)
			if err != nil {
				// handle it ??
			}
			game.sendHitOrMiss(hit)
		} else {
			// tell invalid move
			continue
		}
		game.State = game.checkState()
		game.Turn = !game.Turn
	}

	rematch := game.tellResultsAskRematch()
	return rematch
}

func (game *Game) checkState() string {
	if game.Turn { // player1 turn
		if game.Board2.Hits == 17 {
			return "finished"
		}
	} else if !game.Turn { // player2 turn
		if game.Board1.Hits == 17 {
			return "finished"
		}
	}
	return "playing"

}

func (game *Game) checkValidBoard(board model.Board) bool {
	dim := len(board.Cells)
	if dim == 0 {
		return false
	}

	visited := make([][]bool, dim)
	for i := range visited {
		visited[i] = make([]bool, len(board.Cells[i]))
	}

	required := map[int]int{
		5: 1,
		4: 1,
		3: 2,
		2: 1,
	}
	found := map[int]int{}

	var dfs func(x, y int) int
	dfs = func(x, y int) int {
		if x < 0 || x >= dim || y < 0 || y >= dim {
			return 0
		}
		if visited[x][y] || board.Cells[x][y] != model.Ship {
			return 0
		}
		visited[x][y] = true

		length := 1

		length += dfs(x+1, y)
		length += dfs(x-1, y)
		length += dfs(x, y+1)
		length += dfs(x, y-1)

		return length
	}

	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			if !visited[i][j] && board.Cells[i][j] == model.Ship {
				shipLen := dfs(i, j)
				found[shipLen]++
			}
		}
	}

	for length, reqCount := range required {
		if found[length] != reqCount {
			return false
		}
	}

	for length, count := range found {
		if required[length] == 0 && count > 0 {
			return false
		}
	}

	return true
}

func (game *Game) checkValidMove(x int, y int, boardNo bool) bool {
	if x > game.Board1.Size || y > game.Board1.Size {
		return false
	}
	switch boardNo {
	case true:
		if game.Board1.Cells[x][y] == model.Hit || game.Board1.Cells[x][y] == model.Miss {
			return false
		} else {
			return true
		}
	case false:
		if game.Board2.Cells[x][y] == model.Hit || game.Board2.Cells[x][y] == model.Miss {
			return false
		} else {
			return true
		}
	default:
		return false
	}
}

func (game *Game) PlayMove(x int, y int) (bool, error) {
	if !game.checkValidMove(x, y, game.Turn) {
		return false, fmt.Errorf("move invalid")
	}
	var ret bool
	var err error
	if game.Turn { // player1 turn
		ret, err = game.Board1.ShootCell(x, y)
	} else if !game.Turn { // player2 turn
		ret, err = game.Board2.ShootCell(x, y)
	} else {
		return false, fmt.Errorf("game.Turn somehow not 1 nor 2")
	}
	return ret, err
}

func (game *Game) getBoard(player1 *model.Player, wg *sync.WaitGroup) {
	defer wg.Done()
	// game.checkValidBoard()
	panic("unimplemented")
}

func (game *Game) tellGameStarted() {
	// tells both players game has started, and tells them which is 1st to play (prob not neeeded)
	panic("unimplemented")
}

func (game *Game) getMove(hit bool, x int, y int) (int, int) {
	// tells result of opps move and asks for next move
	// depends on game.Turn
	panic("unimplemented")
}

func (game *Game) sendHitOrMiss(hit bool) {
	// sends hit to player1 or player 2 depending on game.Turn
	panic("unimplemented")
}

func (game *Game) tellResultsAskRematch() bool {
	// tell game finished
	// wait for response for rematch (probably wait 120 seconds before automatic no rematch)
	panic("unimplemented")
}
