package game

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Driver interface {
	GetDirection(*Snake, *Game) Direction
}

type MultiDriver interface {
	GetDirections([]*Game) [][]Direction
}

type KeyboardDriver struct {
	buffer []Direction
	times []time.Time
	maxBufferSize int
}

func newKeyboardDriver(maxBufferSize int) *KeyboardDriver {
	return &KeyboardDriver{maxBufferSize: maxBufferSize}
}

func (driver *KeyboardDriver) GetDirection(*Snake, *Game) (d Direction) {
	if len(driver.buffer) > 0 {
		d, driver.buffer = driver.buffer[0], driver.buffer[1:]
		driver.times = driver.times[1:]
	}
	return
}

func (driver *KeyboardDriver) scanKey(k sdl.Keycode) {
	// purge
	if len(driver.buffer) != 0 && time.Since(driver.times[0]) > time.Second {
		driver.buffer = driver.buffer[1:]
	}

	// scan
	if len(driver.buffer) < driver.maxBufferSize {
		var direction Direction
		switch k {
		case sdl.K_UP:
			direction = Up
		case sdl.K_DOWN:
			direction = Down
		case sdl.K_LEFT:
			direction = Left
		case sdl.K_RIGHT:
			direction = Right
		}
		if direction != (Direction{}) {
			driver.buffer = append(driver.buffer, direction)
			driver.times = append(driver.times, time.Now())
		}
	}
}