package game

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jelisavac-l/GBattleships/internal/model"
)

type Game struct {
	ID      string
	Player1 *model.Player
	Player2 *model.Player
	Board1  *model.Board
	Board2  *model.Board
	Turn    int
	State   string
}

func CreateGame() *Game {
	return &Game{
		ID:    uuid.New().String(),
		State: "waiting",
	}
}

func (game *Game) StartGame() {

	// Request boards

	// Loop moves & check conditions

}

func (game *Game) CheckValidBoard(board model.Board) bool {
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

func (game *Game) checkValidMove(x int, y int, boardNo int) bool { //boardNo being 1 or 2
	if x > game.Board1.Size || y > game.Board1.Size {
		return false
	}
	switch boardNo {
	case 1:
		if game.Board1.Cells[x][y] == model.Hit || game.Board1.Cells[x][y] == model.Miss {
			return false
		} else {
			return true
		}
	case 2:
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
	if game.Turn == 1 {
		ret, err = game.Board1.ShootCell(x, y)
	} else if game.Turn == 2 {
		ret, err = game.Board2.ShootCell(x, y)
	} else {
		return false, fmt.Errorf("game.Turn somehow not 1 nor 2")
	}
	return ret, err
}
