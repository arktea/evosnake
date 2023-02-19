package game

import (
	"fmt"
	"sort"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var colors [][3]uint8 = [][3]uint8{
	{0x00, 0xaa, 0x00}, // Green
	{0xee, 0x00, 0x00}, // Red
	{0x00, 0x00, 0xaa}, // Blue
}

type UI struct {
	window           *sdl.Window
	surface          *sdl.Surface
	tileSize         int
	leaderBoardWidth int
	font             *ttf.Font
}

func NewUI(width, height, tileSize int, leaderBoard bool) *UI {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	winWidth := int32(width * tileSize)
	winHeight := int32(height * tileSize)
	var leaderBoardWidth int32

	if leaderBoard {
		println(winWidth / 2)
		leaderBoardWidth = winWidth / 2
		winWidth += leaderBoardWidth
	}
	if err := ttf.Init(); err != nil {
		panic(err)
	}

	font, err := ttf.OpenFont("./assets/RobotoMono.ttf", 11)
	if err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("Evosnake", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	return &UI{window: window, surface: surface, tileSize: tileSize, leaderBoardWidth: int(leaderBoardWidth), font: font}
}

func (ui *UI) Close() {
	ui.window.Destroy()
	ttf.Quit()
	sdl.Quit()
}

func (e *Element) draw(ui *UI, color [3]uint8) {
	rect := sdl.Rect{X: int32(e.X * ui.tileSize), Y: int32(e.Y * ui.tileSize),
		W: int32(ui.tileSize), H: int32(ui.tileSize)}
	colorCode := sdl.MapRGB(ui.surface.Format, color[0], color[1], color[2])
	ui.surface.FillRect(&rect, colorCode)
}

func (s Snake) draw(ui *UI, color [3]uint8) {
	for _, block := range s.Body {
		block.draw(ui, color)
	}
}

func (g Game) Draw(ui *UI) {
	bgColor := sdl.MapRGB(ui.surface.Format, 0xff, 0xff, 0xff)
	leaderBoardbgColor := sdl.MapRGB(ui.surface.Format, 0xff, 0xe5, 0xcc)
	ui.surface.FillRect(nil, bgColor)
	w, h := ui.window.GetSize()
	rect := sdl.Rect{X: w - int32(ui.leaderBoardWidth), Y: 0, W: int32(ui.leaderBoardWidth), H: h}
	ui.surface.FillRect(&rect, leaderBoardbgColor)
	for i, snake := range g.Snakes {
		snake.draw(ui, colors[i%len(g.Snakes)])
	}
	for _, food := range g.Foods {
		food.draw(ui, [3]uint8{0xff, 0xaa, 00})
	}
	if len(g.Snakes) > 1 {
		g.DisplayLeaderBoard(ui)
	} else {
		g.DisplayState(ui)
	}
	ui.window.UpdateSurface()
}

func (g Game) DisplayLeaderBoard(ui *UI) {
	var err error
	var text *sdl.Surface
	head := "Player  Score Deaths Max"
	if text, err = ui.font.RenderUTF8Blended(head, sdl.Color{R: 0, G: 0, B: 0, A: 255}); err != nil {
		panic(err)
	}
	w, _ := ui.window.GetSize()
	if err = text.Blit(nil, ui.surface, &sdl.Rect{X: w - int32(ui.leaderBoardWidth) + 10, Y: 0, W: 0, H: 0}); err != nil {
		panic(err)
	}
	snakesScores := make([]struct{str string; color [3]uint8; score int}, len(g.Snakes))
	for i, snake := range g.Snakes {
		s, d, m := snake.score, snake.Deaths, snake.MaxScore
		snakesScores[i].str = fmt.Sprintf("Snake %-2d  %-6d%-6d%-5d", i, s, d, m)
		snakesScores[i].color = colors[i%len(g.Snakes)]
		snakesScores[i].score = 10*m -(d*d)
	}
	sort.SliceStable(snakesScores, func (i, j int) bool {
		return snakesScores[i].score > snakesScores[j].score
	})
	for i, snakeScore := range snakesScores {
		if text, err = ui.font.RenderUTF8Blended(snakeScore.str, sdl.Color{R: snakeScore.color[0], G: snakeScore.color[1], B: snakeScore.color[2], A: 255}); err != nil {
			panic(err)
		}
		w, _ := ui.window.GetSize()
		if err = text.Blit(nil, ui.surface, &sdl.Rect{X: w - int32(ui.leaderBoardWidth) + 10, Y: 20*int32(i+1), W: 0, H: 0}); err != nil {
			panic(err)
		}
	}
}

func (g Game) DisplayState(ui *UI) {
	snake := g.Snakes[0]
	scoreString := fmt.Sprintf("Evosnake | score:%d / deaths:%d / max:%d", snake.score, snake.Deaths, snake.MaxScore)
	ui.window.SetTitle(scoreString)
}

func manageEvents(driver Driver) bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.KeyboardEvent:
			driver, ok := driver.(*KeyboardDriver)
			if ok && t.State == sdl.PRESSED {
				driver.scanKey(t.Keysym.Sym)
			}
		}
	}
	return true
}
