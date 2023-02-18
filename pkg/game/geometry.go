package game

import (
	"math"
	"math/rand"
)

type Coordinates struct {
	X, Y int
}

type Direction Coordinates

type Rect struct {
	x, y, w, h int
}

var (
	Up = Direction{X: 0, Y: -1}
	Down = Direction{X: 0, Y: 1}
	Left = Direction{X: -1, Y: 0}
	Right = Direction{X: 1, Y: 0}
)

func randDirection() Direction {
	directions := [...]Direction{Up, Down, Left, Right}
	rand.Shuffle(len(directions), func (i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})
	return directions[0]
}

func (c Coordinates) add(other Coordinates) Coordinates {
	return Coordinates{X: c.X + other.X, Y: c.Y + other.Y}
}

func (c Coordinates) sub(other Coordinates) Coordinates {
	return Coordinates{X: c.X - other.X, Y: c.Y - other.Y}
}

type directionMethod func (Coordinates, Direction) Coordinates

func (c Coordinates) forward(d Direction) Coordinates {
	return c.add(Coordinates(d))
}

func (c Coordinates) backward(d Direction) Coordinates {
	return c.sub(Coordinates(d))
}

func (c Coordinates) left(d Direction) Coordinates {
	return Coordinates{X: c.X + d.Y, Y: c.Y - d.X}
}

func (c Coordinates) right(d Direction) Coordinates {
	return Coordinates{X: c.X - d.Y, Y: c.Y + d.Y}
}

func calcRect(coordinates []Coordinates) Rect {
	var sx, sy []int
	for _, c := range coordinates {
		sx = append(sx, c.X)
		sy = append(sy, c.Y)
	}
	minMax := func (s []int) (min, max int) {
		min, max = math.MaxInt, math.MinInt
		for _, v := range s {
			if v < min {
				min = v
			} else if v > max {
				max = v
			}
		}
		return
	}
	xMin, xMax := minMax(sx)
	yMin, yMax := minMax(sy)
	return Rect{xMin, yMin, xMax-xMin+1, yMax-yMin+1}
}

func (r Rect) getCoordinates() []Coordinates {
	c := make([]Coordinates, r.w*r.h)
	for i := 0; i<r.w; i++ {
		for j := 0; j<r.h; j++ {
			c[i*r.h+j] = Coordinates{X: r.x+i, Y: r.y+j}
		}
	}
	return c
}