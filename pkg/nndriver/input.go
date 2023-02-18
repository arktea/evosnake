package nndriver

import "github.com/taebow/evosnake/pkg/game"


func gameToInput(s *game.Snake, food *game.Element, b *game.Board) []float64 {
	res := make([]float64, 8)
	switch s.Direction {
	case game.Left, game.Right:
		res[0] = float64(s.Direction.X)
	case game.Up, game.Down:
		res[1] = float64(s.Direction.Y)
	}
	head := s.Body[0]
	for i, foodRelPos := range [...]int{
		food.X - head.X,
		food.Y - head.Y,
	}{
		if foodRelPos > 0 {
			res[i+2] = 1
		} else if foodRelPos < 0 {
			res[i+2] = -1
		}
	}

	for i, c := range []game.Coordinates{
		{X: head.X - 1, Y: head.Y},
		{X: head.X + 1, Y: head.Y},
		{X: head.X, Y: head.Y - 1},
		{X: head.X, Y: head.Y + 1},
	}{
		if b.IsDanger(c) {
			res[i+4] = 1
		} else {
			res[i+4] = -1
		}
	}

	return res
}