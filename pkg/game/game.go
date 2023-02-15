package game

import (
	"time"
	"math/rand"
)

type Game struct {
	initSnakeSize int
	Board *Board
	Snakes []*Snake
	Foods []*Element
}

func NewGame(width, height, initSnakeSize, numSnakes, numFoods int) *Game {
	game := &Game{initSnakeSize: initSnakeSize}
	game.Board = newBoard(width, height)
	game.initSnakes(numSnakes)
	game.initFoods(numFoods)
	return game
}

func (g *Game) initSnakes(n int) {
	g.Snakes = make([]*Snake, n)
	for i := range g.Snakes {
		g.Snakes[i] = g.Board.newSnake(g.initSnakeSize)
	}
}

func (g *Game) initFoods(n int) {
	g.Foods = make([]*Element, n)
	for i := range g.Foods {
		g.Foods[i] = g.Board.newFood()
	}
}

func (g *Game) update(directions ...Direction) {
	for i, snake := range g.Snakes {
		if i < len(directions) {
			snake.UpdateDirection(directions[i])
		}
		if collide, elem := snake.nextMoveCollide(g.Board); collide {
			if elem == nil || elem.elementType == Block {
				snake.die(g.Board)
				g.Board.respawnSnake(snake, g.initSnakeSize)
			} else if elem.elementType == Food {
				snake.eat(g.Board)
				snake.move(g.Board)
				g.Board.respawnFood(elem)
			}
		} else {
			snake.move(g.Board)
		}
	}
}

func (g *Game) getDirections(drivers []Driver) []Direction {
	directions := make([]Direction, len(drivers))
	for i, driver := range(drivers) {
		directions[i] = driver.GetDirection(g.Snakes[i], g)
	}
	return directions
}

func (g *Game) Run(rounds, frameRate int, gui bool, drivers ...Driver) {
	var ui *UI
	if gui {
		ui = NewUI(g.Board.Width, g.Board.Height, 8)
		defer ui.Close()
	}
	lastRefresh := time.Now()
	running := true
	for running {
		if rounds >= 0 {
			rounds--
			if rounds < 0 {
				break
			}
		}
		if time.Since(lastRefresh) > time.Second/time.Duration(frameRate) {
			g.update(g.getDirections(drivers)...)
			if gui {
				g.Draw(ui)
				g.DisplayState(ui)
			}
			lastRefresh = time.Now()
		}
		running = manageEvents(drivers[0])
	}
}

func RunMulti(games []*Game, rounds int, multiDriver MultiDriver) {
	for i := 0; i<rounds; i++ {
		directions := multiDriver.GetDirections(games)
		for i, g := range games {
			g.update(directions[i]...)
		}
	}
}

func PlayManual(frameRate int) {
	rand.Seed(time.Now().UTC().UnixNano())
	game := NewGame(50, 50, 5, 1, 1)
	game.Run(-1, frameRate, true, newKeyboardDriver(3))
}