package game

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type UI struct {
	window *sdl.Window
	surface *sdl.Surface
	tileSize int
}

func NewUI(width, height, tileSize int) *UI {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	
	window, err := sdl.CreateWindow("evosnake", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(width * tileSize), int32(height * tileSize), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	
	return &UI{window, surface, tileSize}
}

func (ui *UI) Close() {
	ui.window.Destroy()
	sdl.Quit()
}

func (e *Element) draw(ui *UI) {
	rect := sdl.Rect{X: int32(e.X * ui.tileSize), Y: int32(e.Y * ui.tileSize), 
		W: int32(ui.tileSize), H: int32(ui.tileSize)}
	colorCode := sdl.MapRGB(ui.surface.Format, e.color[0], e.color[1], e.color[2])
	ui.surface.FillRect(&rect, colorCode)
}

func (s Snake) draw(ui *UI) {
	for _, block := range s.Body {
		block.draw(ui)
	}
}

func (g Game) Draw(ui *UI) {
	bgColor := sdl.MapRGB(ui.surface.Format, 0xee, 0xee, 0xee)
	ui.surface.FillRect(nil, bgColor)
	for _, snake := range g.Snakes {
		snake.draw(ui)
	}
	for _, food := range g.Foods {
		food.draw(ui)	
	}
	ui.window.UpdateSurface()
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