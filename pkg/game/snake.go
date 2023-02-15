package game

type Snake struct {
	body      []*Element
	direction Direction
	color     [3]uint8
	Deaths    int
	MaxScore  int
	score     int
}

func newSnake(c Coordinates, size int, direction Direction) *Snake {
	color := [...]uint8{0x00, 0xaa, 0x00}
	snake := &Snake{direction: direction, color: color}
	snake.initBody(c, size)
	return snake
}

func (s *Snake) respawn(c Coordinates, size int, direction Direction) {
	s.direction = direction
	s.score = 0
	s.initBody(c, size)
}

func (s *Snake) getSpace(forward, backward, left, right int) Rect {
	sideFunc := func(l int, start Coordinates, f directionMethod) (res Coordinates) {
		res = start
		for i := 0; i < l; i++ {
			res = f(res, s.direction)
		}
		return
	}
	sides := []Coordinates{
		sideFunc(forward, s.body[0].Coordinates, Coordinates.forward),
		sideFunc(left, s.body[1].Coordinates, Coordinates.left),
		sideFunc(right, s.body[1].Coordinates, Coordinates.right),
		sideFunc(backward, s.body[len(s.body)-1].Coordinates, Coordinates.backward),
	}
	return calcRect(sides)
}

func (s *Snake) hasSpace(b *Board) bool {
	var forward int
	switch s.direction {
	case Up, Down:
		forward = b.Height / 2
	case Left, Right:
		forward = b.Width / 2
	}
	rect := s.getSpace(forward, 2, 2, 2)
	collide, _ := b.collide(rect.getCoordinates()...)
	return !collide
}

func (s *Snake) initBody(c Coordinates, size int) {
	s.body = []*Element{s.newBodyElement(c)}
	for i := 0; i < size-1; i++ {
		s.grow()
	}
}

func (s *Snake) grow() {
	var dir Direction
	tail := s.body[len(s.body)-1]
	if len(s.body) == 1 {
		dir = s.direction
	} else {
		tail2 := s.body[len(s.body)-2]
		dir = Direction(tail2.sub(tail.Coordinates))
	}
	c := tail.backward(dir)
	s.body = append(s.body, s.newBodyElement(c))
}

func (s *Snake) newBodyElement(c Coordinates) *Element {
	return newElement(c, s.color, Block)
}

func (s *Snake) move(b *Board) {
	c := s.body[0].forward(s.direction)
	newHead := s.newBodyElement(c)
	b.set(newHead)
	b.free(s.body[len(s.body)-1])
	s.body = append([]*Element{newHead}, s.body[:len(s.body)-1]...)
}

func (s *Snake) UpdateDirection(d Direction) {
	checkCoordinates := Coordinates(s.direction).add(Coordinates(d))
	if checkCoordinates.x != 0 && checkCoordinates.y != 0 {
		s.direction = d
	}
}

func (s *Snake) die(b *Board) {
	s.Deaths++
	s.score = 0
	b.free(s.body...)
}

func (s *Snake) eat(b *Board) {
	s.grow()
	s.score++
	if s.score > s.MaxScore {
		s.MaxScore = s.score
	}
	b.set(s.body[len(s.body)-1])
}

func (s *Snake) nextMoveCollide(board *Board) (bool, *Element) {
	nextCoordinates := s.body[0].Coordinates.forward(s.direction)
	return board.collide(nextCoordinates)
}

func (s *Snake) See(food *Element, b *Board) []float64 {
	res := make([]float64, 8)
	switch s.direction {
	case Left, Right:
		res[0] = float64(s.direction.x)
	case Up, Down:
		res[1] = float64(s.direction.y)
	}
	foodX := food.x - s.body[0].x
	foodY := food.y - s.body[0].y
	if foodX > 0 {
		res[2] = 1
	} else if foodX < 0 {
		res[2] = -1
	} else {
		res[2] = 0
	}

	if foodY > 0 {
		res[3] = 1
	} else if foodY < 0 {
		res[3] = -1
	} else {
		res[3] = 0
	}

	coordLeft := Coordinates{s.body[0].x - 1, s.body[0].y}
	coordRight := Coordinates{s.body[0].x + 1, s.body[0].y}
	coordUp := Coordinates{s.body[0].x, s.body[0].y - 1}
	coordDown := Coordinates{s.body[0].x, s.body[0].y + 1}

	if coordLeft.danger(b) {
		res[4] = 1
	} else {
		res[4] = -1
	}
	if coordRight.danger(b) {
		res[5] = 1
	} else {
		res[5] = -1
	}
	if coordUp.danger(b) {
		res[6] = 1
	} else {
		res[6] = -1
	}
	if coordDown.danger(b) {
		res[7] = 1
	} else {
		res[7] = -1
	}

	return res
}

func (c Coordinates) danger(b *Board) bool {
	if collide, _ := b.collide(c); collide {
		elem := b.get(c)
		if elem == nil || elem.elementType == Block {
			return true
		}
	}
	return false
}
