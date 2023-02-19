package game

type ElementType int

const (
	Block ElementType = iota
	Food
)

type Element struct {
	Coordinates
	elementType ElementType
}

func newElement(c Coordinates, t ElementType) *Element {
	return &Element{c, t}
}
