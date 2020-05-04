package game

import (
	"errors"
	"fmt"
)

/**
 * Constants
 */

const FieldWidth = 3
const FieldHeight = 3

/**
 * Types
 */

type cellType rune

const (
	cellEmpty cellType = ' '
	cellX     cellType = 'X'
	cellO     cellType = 'O'
)

type PlayerType int

const (
	PlayerNone PlayerType = iota
	PlayerX
	PlayerO
)

type Game struct {
	cells [FieldWidth * FieldHeight]cellType
}

/**
 * Methods
 */

func New() *Game {
	game := Game{}
	game.Reset()
	return &game
}

func (game *Game) Reset() {
	for i := 0; i < FieldWidth*FieldHeight; i++ {
		game.cells[i] = cellEmpty
	}
}

func (game Game) Draw() {
	for i := 0; i < FieldWidth*2+1; i++ {
		fmt.Printf("-")
	}
	fmt.Println()

	for line := 0; line < FieldHeight; line++ {
		for row := 0; row < FieldWidth; row++ {
			fmt.Printf("|%c", game.cells[line*FieldHeight+row])
		}
		fmt.Println("|")

		for i := 0; i < FieldWidth*2+1; i++ {
			fmt.Printf("-")
		}
		fmt.Println()
	}
}

func (game *Game) MakeStep(player PlayerType, pos int) error {
	if pos < 0 || pos >= (FieldWidth*FieldHeight) {
		return errors.New("Invalid position value")
	}

	if game.cells[pos] != cellEmpty {
		return errors.New("Field is already busy")
	}

	switch player {
	case PlayerX:
		game.cells[pos] = cellX
	case PlayerO:
		game.cells[pos] = cellO
	default:
		return errors.New("Wrong player type")
	}

	return nil
}

func getPlayerType(cell cellType) PlayerType {
	switch cell {
	case cellX:
		return PlayerX
	case cellO:
		return PlayerO
	}
	return PlayerNone
}

func (game Game) CheckWin() (bool, PlayerType) {
	var flag bool
	// Lines
	for line := 0; line < FieldHeight; line++ {
		if game.cells[0+line*FieldWidth] == cellEmpty {
			continue
		}

		flag = true
		for row := 0; row < FieldWidth-1; row++ {
			if game.cells[row+line*FieldWidth] != game.cells[(row+1)+line*FieldWidth] {
				flag = false
				break
			}
		}

		if !flag {
			continue
		}

		return true, getPlayerType(game.cells[0+line*FieldWidth])
	}

	// Rows
	for row := 0; row < FieldWidth; row++ {
		if game.cells[row+0*FieldWidth] == cellEmpty {
			continue
		}

		flag = true
		for line := 0; line < FieldHeight-1; line++ {
			if game.cells[row+line*FieldWidth] != game.cells[row+(line+1)*FieldWidth] {
				flag = false
				break
			}
		}

		if !flag {
			continue
		}

		return true, getPlayerType(game.cells[row+0*FieldWidth])
	}

	// Diagonals
	if FieldWidth == FieldHeight {
		// From top left to bottom right
		if game.cells[0] != cellEmpty {
			flag = true
			for i := 0; i < FieldWidth-1; i++ {
				if game.cells[i+i*FieldWidth] != game.cells[(i+1)+(i+1)*FieldWidth] {
					flag = false
					break
				}
			}
		}

		if flag {
			return true, getPlayerType(game.cells[0])
		}

		// From top right to bottom left
		if game.cells[FieldWidth-1] != cellEmpty {
			flag = true
			for i := 0; i < FieldWidth-1; i++ {
				if game.cells[(FieldWidth-1-i)+i*FieldWidth] != game.cells[(FieldWidth-2-i)+(i+1)*FieldWidth] {
					flag = false
					break
				}
			}
		}

		if flag {
			return true, getPlayerType(game.cells[FieldWidth-1])
		}
	}

	// Draw
	flag = true
	for _, val := range game.cells {
		if val == cellEmpty {
			flag = false
			break
		}
	}
	if flag {
		return true, PlayerNone
	}

	return false, PlayerNone
}

func (game Game) GameToSlice() [][]string {
	array := make([][]string, 0)

	for line := 0; line < FieldHeight; line++ {
		tmp := make([]string, 0)
		for row := 0; row < FieldWidth; row++ {
			switch game.cells[line * FieldWidth + row] {
			case cellEmpty:
				tmp = append(tmp, " ")
			case cellX:
				tmp = append(tmp, "X")
			case cellO:
				tmp = append(tmp, "O")
			}
		}
		array = append(array, tmp)
	}
	return array
}
