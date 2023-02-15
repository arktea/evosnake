package game

import (
	"math"
	"math/rand"
)

type Coordinates struct {
	x, y int
}

type Direction Coordinates

type Rect struct {
	x, y, w, h int
}

var (
	Up = Direction{x: 0, y: -1}
	Down = Direction{x: 0, y: 1}
	Left = Direction{x: -1, y: 0}
	Right = Direction{x: 1, y: 0}
)

func randDirection() Direction {
	directions := [...]Direction{Up, Down, Left, Right}
	rand.Shuffle(len(directions), func (i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})
	return directions[0]
}

func (c Coordinates) add(other Coordinates) Coordinates {
	return Coordinates{x: c.x + other.x, y: c.y + other.y}
}

func (c Coordinates) sub(other Coordinates) Coordinates {
	return Coordinates{x: c.x - other.x, y: c.y - other.y}
}

type directionMethod func (Coordinates, Direction) Coordinates

func (c Coordinates) forward(d Direction) Coordinates {
	return c.add(Coordinates(d))
}

func (c Coordinates) backward(d Direction) Coordinates {
	return c.sub(Coordinates(d))
}

func (c Coordinates) left(d Direction) Coordinates {
	return Coordinates{x: c.x + d.y, y: c.y - d.x}
}

func (c Coordinates) right(d Direction) Coordinates {
	return Coordinates{x: c.x - d.y, y: c.y + d.x}
}

func calcRect(coordinates []Coordinates) Rect {
	var sx, sy []int
	for _, c := range coordinates {
		sx = append(sx, c.x)
		sy = append(sy, c.y)
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
			c[i*r.h+j] = Coordinates{x: r.x+i, y: r.y+j}
		}
	}
	return c
}