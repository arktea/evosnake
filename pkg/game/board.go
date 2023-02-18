package game

import "math/rand"

type Board struct {
	Width, Height int
	board [][]*Element
}

func newBoard(width, height int) *Board {
	board := make([][]*Element, width)
	for i := range board {
		board[i] = make([]*Element, height)
	}
	return &Board{board: board, Width: width, Height: height}
}

func (b *Board) isValid(c Coordinates) bool {
	return c.X >= 0 && c.X < b.Width && c.Y >= 0 && c.Y < b.Height
}

func (b *Board) set(elems ...*Element) {
	for _, e := range elems {
		if b.isValid(e.Coordinates) {
			b.board[e.X][e.Y] = e
		}
	}
}

func (b *Board) get(c Coordinates) *Element {
	if b.isValid(c) {
		return b.board[c.X][c.Y]
	}
	return nil
}

func (b *Board) free(elems ...*Element) {
	for _, e := range elems {
		if b.isValid(e.Coordinates) {
			b.board[e.X][e.Y] = nil
		}
	}
}

func (b *Board) Collide(coordinates ...Coordinates) (bool, *Element) {
	for _, c := range coordinates {
		elem := b.get(c)
		if elem != nil {
			return true, elem
		} else if !b.isValid(c) {
			return true, nil
		}
	}
	return false, nil
}

func (b *Board) randCoordinates() Coordinates {
	for {
		c := Coordinates{rand.Int() % b.Width, rand.Int() % b.Height}
		if collide, _ := b.Collide(c); !collide {
			return c
		}
	}
}

func (b *Board) newFood() *Element {
	elem := newElement(b.randCoordinates(), [3]uint8{0xff, 0xaa, 00}, Food)
	b.set(elem)
	return elem
}

func (b *Board) respawnFood(food *Element) {
	food.Coordinates = b.randCoordinates()
	b.set(food)
}

func (b *Board) newSnake(size int) *Snake {
	for {
		s := newSnake(b.randCoordinates(), size, randDirection())
		if s.hasSpace(b) {
			b.set(s.Body...)
			return s
		} 
	}
}

func (b *Board) respawnSnake(s *Snake, size int) {
	for {
		s.respawn(b.randCoordinates(), size, randDirection())
		if s.hasSpace(b) {
			b.set(s.Body...)
			break
		}
	}
}

func (b *Board) IsDanger(c Coordinates) bool {
	if collide, _ := b.Collide(c); collide {
		elem := b.get(c)
		if elem == nil || elem.elementType == Block {
			return true
		}
	}
	return false
}