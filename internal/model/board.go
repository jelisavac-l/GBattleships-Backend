package model

type CellState int

const (
	Empty CellState = iota
	Ship
	Hit
)

type Board struct {
	Size  int
	Cells [][]CellState
}

func NewBoard(cells [][]CellState) *Board {
	return &Board{
		Size:  len(cells),
		Cells: cells,
	}
}
