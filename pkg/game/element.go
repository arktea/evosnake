package game

type ElementType int

const (
	Block ElementType = iota
	Food
)

type Element struct {
	Coordinates
	color [3]uint8
	elementType ElementType
}

func newElement(c Coordinates, color [3]uint8, t ElementType) *Element {
	return &Element{c, color, t}
}
