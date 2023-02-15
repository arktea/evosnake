package game

import (
	"time"
	"math/rand"

	// "github.com/veandco/go-sdl2/sdl"
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

func (g *Game) initBrains(drivers []Driver) {
	numDrivers := len(drivers)
	for i, snake := range g.Snakes {
		if i < numDrivers {
			snake.driver = drivers[i]
		}
	}
}

func (g *Game) Update() {
	for _, snake := range g.Snakes {
		snake.UpdateDirection(g)
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

func (g *Game) Run(rounds, frameRate int, gui bool, drivers []Driver) {
	var ui *UI
	if gui {
		ui = NewUI(g.Board.Width, g.Board.Height, 8)
		defer ui.Close()
	}
	g.initBrains(drivers)
	lastRefresh := time.Now()
	running := true
	for running {
		if time.Since(lastRefresh) > time.Second/time.Duration(frameRate) {
			lastRefresh = time.Now()
			g.Update()
			if gui {
				g.Draw(ui)
				g.DisplayState(ui)
			}
		}
		if rounds >= 0 {
			rounds--
			if rounds <= 0 {
				break
			}
		}
		running = manageEvents(drivers[0])
	}
}

func PlayManual(frameRate int) {
	rand.Seed(time.Now().UTC().UnixNano())
	game := NewGame(50, 50, 5, 1, 1)
	game.Run(-1, frameRate, true, []Driver{newKeyboardDriver(3)})
}