package model

import "fmt"

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
		Size:  dim,
	}
}

func (board *Board) SetCells(cells [][]CellState) error {
	if len(cells) != board.Size {
		return fmt.Errorf("size of input cells not matching board size")
	}
	board.Cells = cells
	return nil
}

func (board *Board) ShootCell(x int, y int) (bool, error) {
	// true if hit, false if miss
	switch board.Cells[x][y] {
	case Empty:
		board.Cells[x][y] = Miss
		return false, nil
	case Ship:
		board.Cells[x][y] = Hit
		return true, nil
	}
	return false, fmt.Errorf("ShootCell: cell is not empty nor ship")
}
