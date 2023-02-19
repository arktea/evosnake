package game

import (
	"fmt"
	"sort"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var colors [][3]uint8 = [][3]uint8{
	{0x00, 0xaa, 0x00}, // Green
	{0xff, 0x33, 0x33}, // Red
	{0x33, 0x99, 0xff}, // Blue
	{0x99, 0x33, 0xff}, // Purple
	{0xff, 0x80, 0x00}, // Orange
}

type UI struct {
	window           *sdl.Window
	surface          *sdl.Surface
	tileSize         int32
	leaderBoardWidth int32
	font             *ttf.Font
}

func newUI(width, height, tileSize int, leaderBoard bool) *UI {
	winWidth := int32(width * tileSize)
	winHeight := int32(height * tileSize)
	var leaderBoardWidth int32

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	if leaderBoard {
		leaderBoardWidth = winWidth / 2
		winWidth += leaderBoardWidth
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

	if err := ttf.Init(); err != nil {
		panic(err)
	}

	font, err := ttf.OpenFont("./assets/RobotoMono.ttf", 11)
	if err != nil {
		panic(err)
	}

	return &UI{
		window: window, 
		surface: surface, 
		tileSize: int32(tileSize), 
		leaderBoardWidth: leaderBoardWidth, 
		font: font,
	}
}

func (ui *UI) close() {
	ui.window.Destroy()
	ttf.Quit()
	sdl.Quit()
}

func (ui *UI) drawElement(e *Element, color [3]uint8) {
	rect := sdl.Rect{X: int32(e.X) * ui.tileSize, Y: int32(e.Y) * ui.tileSize,
		W: int32(ui.tileSize), H: int32(ui.tileSize)}
	colorCode := sdl.MapRGB(ui.surface.Format, color[0], color[1], color[2])
	ui.surface.FillRect(&rect, colorCode)
}

func (ui *UI) drawSnake(s *Snake, color [3]uint8) {
	for _, block := range s.Body {
		ui.drawElement(block, color)
	}
}

func (ui *UI) drawGame(g *Game) {
	bgColor := sdl.MapRGB(ui.surface.Format, 0xff, 0xff, 0xff)
	leaderBoardbgColor := sdl.MapRGB(ui.surface.Format, 0xff, 0xe5, 0xcc)
	w, h := ui.window.GetSize()
	lbRect := sdl.Rect{X: w - int32(ui.leaderBoardWidth), Y: 0, W: int32(ui.leaderBoardWidth), H: h}

	ui.surface.FillRect(nil, bgColor)
	ui.surface.FillRect(&lbRect, leaderBoardbgColor)
	
	for i, snake := range g.Snakes {
		ui.drawSnake(snake, colors[i%len(g.Snakes)])
	}
	for _, food := range g.Foods {
		ui.drawElement(food, [3]uint8{0xff, 0xaa, 00})
	}
	if len(g.Snakes) > 1 {
		ui.drawLeaderBoard(g)
	} else {
		ui.displaySnakeInTitle(g.Snakes[0], g)
	}
	ui.window.UpdateSurface()
}
 
func (ui *UI) drawLBLine(line string, color [3]uint8, i int) {
	text, _ := ui.font.RenderUTF8Blended(line, sdl.Color{R: color[0], G: color[1], B: color[2], A: 255})
	w, _ := ui.window.GetSize()
	text.Blit(nil, ui.surface, &sdl.Rect{X: w - int32(ui.leaderBoardWidth) + 10, Y: 20*int32(i), W: 0, H: 0})
}

func (ui *UI) drawLeaderBoard(g *Game) {
	ui.drawLBLine("Player  Score Deaths Max", [...]uint8{0, 0, 0}, 0)
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
		ui.drawLBLine(snakeScore.str, snakeScore.color, i+1)
	}
}

func (ui *UI) displaySnakeInTitle(s *Snake, g *Game) {
	scoreString := fmt.Sprintf("Evosnake | score:%d / deaths:%d / max:%d", s.score, s.Deaths, s.MaxScore)
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
