package game

type Snake struct {
	Body      []*Element
	Direction Direction
	color     [3]uint8
	Deaths    int
	MaxScore  int
	score     int
}

func newSnake(c Coordinates, size int, direction Direction) *Snake {
	color := [...]uint8{0x00, 0xaa, 0x00}
	snake := &Snake{Direction: direction, color: color}
	snake.initBody(c, size)
	return snake
}

func (s *Snake) respawn(c Coordinates, size int, direction Direction) {
	s.Direction = direction
	s.score = 0
	s.initBody(c, size)
}

func (s *Snake) getSpace(forward, backward, left, right int) Rect {
	sideFunc := func(l int, start Coordinates, f directionMethod) (res Coordinates) {
		res = start
		for i := 0; i < l; i++ {
			res = f(res, s.Direction)
		}
		return
	}
	sides := []Coordinates{
		sideFunc(forward, s.Body[0].Coordinates, Coordinates.forward),
		sideFunc(left, s.Body[1].Coordinates, Coordinates.left),
		sideFunc(right, s.Body[1].Coordinates, Coordinates.right),
		sideFunc(backward, s.Body[len(s.Body)-1].Coordinates, Coordinates.backward),
	}
	return calcRect(sides)
}

func (s *Snake) hasSpace(b *Board) bool {
	var forward int
	switch s.Direction {
	case Up, Down:
		forward = b.Height / 2
	case Left, Right:
		forward = b.Width / 2
	}
	rect := s.getSpace(forward, 2, 2, 2)
	collide, _ := b.Collide(rect.getCoordinates()...)
	return !collide
}

func (s *Snake) initBody(c Coordinates, size int) {
	s.Body = []*Element{s.newBodyElement(c)}
	for i := 0; i < size-1; i++ {
		s.grow()
	}
}

func (s *Snake) grow() {
	var dir Direction
	tail := s.Body[len(s.Body)-1]
	if len(s.Body) == 1 {
		dir = s.Direction
	} else {
		tail2 := s.Body[len(s.Body)-2]
		dir = Direction(tail2.sub(tail.Coordinates))
	}
	c := tail.backward(dir)
	s.Body = append(s.Body, s.newBodyElement(c))
}

func (s *Snake) newBodyElement(c Coordinates) *Element {
	return newElement(c, s.color, Block)
}

func (s *Snake) move(b *Board) {
	c := s.Body[0].forward(s.Direction)
	newHead := s.newBodyElement(c)
	b.set(newHead)
	b.free(s.Body[len(s.Body)-1])
	s.Body = append([]*Element{newHead}, s.Body[:len(s.Body)-1]...)
}

func (s *Snake) UpdateDirection(d Direction) {
	checkCoordinates := Coordinates(s.Direction).add(Coordinates(d))
	if checkCoordinates.X != 0 && checkCoordinates.Y != 0 {
		s.Direction = d
	}
}

func (s *Snake) die(b *Board) {
	s.Deaths++
	s.score = 0
	b.free(s.Body...)
}

func (s *Snake) eat(b *Board) {
	s.grow()
	s.score++
	if s.score > s.MaxScore {
		s.MaxScore = s.score
	}
	b.set(s.Body[len(s.Body)-1])
}

func (s *Snake) nextMoveCollide(board *Board) (bool, *Element) {
	nextCoordinates := s.Body[0].Coordinates.forward(s.Direction)
	return board.Collide(nextCoordinates)
}

