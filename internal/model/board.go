package model

type CellState int

const (
	Empty CellState = iota
	Ship
	Hit
	Miss
)

type Board struct {
	Size  int
	Cells [][]CellState
}

func NewBoard(dim int) *Board {
	return &Board{
		Cells: make([][]CellState, dim),
		Size:  dim * dim,
	}
}

func (board *Board) SetCells(cells [][]CellState) {
	if len(cells) != board.Size {
		// error stuff
		return
	}
	board.Cells = cells
}

func (board *Board) ShootCell(x int, y int) {
	// check cell if ship or empty, set to hit or miss, return hit or miss
}
