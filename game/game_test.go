package game

import "testing"

func TestReset(t *testing.T) {
	t.Parallel()
	game := NewGame()
	for i := 0; i < FieldWidth*FieldHeight; i++ {
		err := game.MakeStep(PlayerX, i)
		if err != nil {
			t.Errorf("Can't set value: %v", err)
		}
	}

	game.Reset()

	for i := 0; i < FieldWidth*FieldHeight; i++ {
		if err := game.MakeStep(PlayerX, i); err != nil {
			t.Errorf("Can't set value after reset: %v", err)
		}
	}
}

func TestMakeStep(t *testing.T) {
	t.Parallel()
	game := NewGame()
	if err := game.MakeStep(PlayerX, -1); err == nil {
		t.Errorf("Should be error")
	}

	if err := game.MakeStep(PlayerX, FieldWidth*FieldHeight); err == nil {
		t.Errorf("Should be error")
	}

	for i := 0; i < FieldWidth*FieldHeight; i++ {
		if err := game.MakeStep(PlayerX, i); err != nil {
			t.Errorf("Can't set value: %v", err)
		}
	}
}

func TestCheckWin(t *testing.T) {
	t.Parallel()
	game := NewGame()

	players := [2]PlayerType{PlayerX, PlayerO}
	for i := 0; i < len(players); i++ {
		// Lines
		for line := 0; line < FieldHeight; line++ {
			game.Reset()
			for row := 0; row < FieldWidth; row++ {
				err := game.MakeStep(players[i], row+line*FieldWidth)
				if err != nil {
					t.Errorf("Can't set value: %v", err)
				}
			}

			if ret, pType := game.CheckWin(); ret != true {
				t.Errorf("The victory was not marked!")
			} else if pType != players[i] {
				game.Draw()
				t.Errorf("Wrong victory player type!: %d", pType)
			}
		}

		// Rows
		for row := 0; row < FieldWidth; row++ {
			game.Reset()
			for line := 0; line < FieldHeight; line++ {
				err := game.MakeStep(players[i], row+line*FieldWidth)
				if err != nil {
					t.Errorf("Can't set value: %v", err)
				}
			}

			if ret, pType := game.CheckWin(); ret != true {
				t.Errorf("The victory was not marked!")
			} else if pType != players[i] {
				t.Errorf("Wrong victory player type!")
			}
		}

		// Diagonals
		if FieldWidth != FieldHeight {
			t.Errorf("Can't be diagonale win, field size: %dx%d", FieldWidth, FieldHeight)
		}
		// From left top to right bottom
		game.Reset()
		for j := 0; j < FieldWidth; j++ {
			err := game.MakeStep(players[i], j+j*FieldWidth)
			if err != nil {
				t.Errorf("Can't set value: %v", err)
			}
		}
		if ret, pType := game.CheckWin(); ret != true {
			t.Errorf("The victory was not marked!")
		} else if pType != players[i] {
			t.Errorf("Wrong victory player type!")
		}
		// From right top to left bottom
		game.Reset()
		for j := 0; j < FieldWidth; j++ {
			err := game.MakeStep(players[i], (FieldWidth-1-j)+j*FieldWidth)
			if err != nil {
				t.Errorf("Can't set value: %v", err)
			}
		}
		if ret, pType := game.CheckWin(); ret != true {
			t.Errorf("The victory was not marked!")
		} else if pType != players[i] {
			t.Errorf("Wrong victory player type!")
		}
	}
}
