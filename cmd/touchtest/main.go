package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

// Simulate mobile environment for testing
func isTouchPrimaryInput() bool {
	return true // This would normally be false on desktop
}

type InputManager struct {
	pressed     map[ebiten.Key]struct{}
	prevPressed map[ebiten.Key]struct{}
	touchMode   bool
}

func NewInputManager() *InputManager {
	return &InputManager{
		pressed:     make(map[ebiten.Key]struct{}),
		prevPressed: make(map[ebiten.Key]struct{}),
		touchMode:   isTouchPrimaryInput(), // Initialize based on platform
	}
}

func (i *InputManager) Update() {
	// Save previous state
	i.prevPressed = make(map[ebiten.Key]struct{})
	for k := range i.pressed {
		i.prevPressed[k] = struct{}{}
	}
	i.pressed = make(map[ebiten.Key]struct{})

	// Handle keyboard input
	keys := []ebiten.Key{
		ebiten.KeyEnter,
		ebiten.KeySpace,
		ebiten.KeyLeft,
		ebiten.KeyDown,
		ebiten.KeyRight,
	}
	for _, k := range keys {
		if ebiten.IsKeyPressed(k) {
			i.pressed[k] = struct{}{}
		}
	}

	// Process mouse for testing on desktop
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y >= screenHeight-64 {
			switch {
			case screenWidth <= x:
				// Nothing assigned to far right
			case screenWidth*3/4 <= x:
				// Action buttons
				i.pressed[ebiten.KeyEnter] = struct{}{}
				i.pressed[ebiten.KeySpace] = struct{}{}
			case screenWidth*2/4 <= x:
				i.pressed[ebiten.KeyDown] = struct{}{}
			case screenWidth/4 <= x:
				i.pressed[ebiten.KeyRight] = struct{}{}
			default:
				i.pressed[ebiten.KeyLeft] = struct{}{}
			}
		}
	}

	// Process touch input
	for _, t := range ebiten.TouchIDs() {
		x, y := ebiten.TouchPosition(t)
		if y >= screenHeight-64 {
			switch {
			case screenWidth <= x:
				// Nothing assigned to far right
			case screenWidth*3/4 <= x:
				// Action buttons
				i.pressed[ebiten.KeyEnter] = struct{}{}
				i.pressed[ebiten.KeySpace] = struct{}{}
			case screenWidth*2/4 <= x:
				i.pressed[ebiten.KeyDown] = struct{}{}
			case screenWidth/4 <= x:
				i.pressed[ebiten.KeyRight] = struct{}{}
			default:
				i.pressed[ebiten.KeyLeft] = struct{}{}
			}
		}
	}

	// Enable touch mode if touches/mouse detected or if we're on mobile
	if len(ebiten.TouchIDs()) > 0 || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) || isTouchPrimaryInput() {
		i.touchMode = true
	}
}

func (i *InputManager) IsKeyPressed(key ebiten.Key) bool {
	_, ok := i.pressed[key]
	return ok
}

func (i *InputManager) IsKeyJustPressed(key ebiten.Key) bool {
	_, ok := i.pressed[key]
	if !ok {
		return false
	}
	_, ok = i.prevPressed[key]
	return !ok
}

func (i *InputManager) IsTouchEnabled() bool {
	return i.touchMode || isTouchPrimaryInput()
}

type Game struct {
	input *InputManager
}

func NewGame() *Game {
	return &Game{
		input: NewInputManager(),
	}
}

func (g *Game) Update() error {
	g.input.Update()

	// Log when keys/touches are just pressed
	if g.input.IsKeyJustPressed(ebiten.KeyLeft) {
		log.Println("Left pressed")
	}
	if g.input.IsKeyJustPressed(ebiten.KeyRight) {
		log.Println("Right pressed")
	}
	if g.input.IsKeyJustPressed(ebiten.KeyDown) {
		log.Println("Down pressed")
	}
	if g.input.IsKeyJustPressed(ebiten.KeySpace) || g.input.IsKeyJustPressed(ebiten.KeyEnter) {
		log.Println("Action pressed")
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw info text
	msg := fmt.Sprintf("Touch/Keyboard Test\n\nTouch enabled: %v\n\nControls: \n - Arrow keys or touch areas\n - Space/Enter for action",
		g.input.IsTouchEnabled())
	ebitenutil.DebugPrint(screen, msg)

	// Draw touch control areas if touch is enabled
	if g.input.IsTouchEnabled() {
		buttonH := 64.0
		buttonY := float64(screenHeight - 64)

		// Left movement
		ebitenutil.DrawRect(screen, 0, buttonY, float64(screenWidth/4), buttonH, color.RGBA{R: 255, A: 64})
		// Right movement
		ebitenutil.DrawRect(screen, float64(screenWidth/4), buttonY, float64(screenWidth/4), buttonH, color.RGBA{G: 255, A: 64})
		// Down button
		ebitenutil.DrawRect(screen, float64(screenWidth*2/4), buttonY, float64(screenWidth/4), buttonH, color.RGBA{B: 255, A: 64})
		// Action button
		ebitenutil.DrawRect(screen, float64(screenWidth*3/4), buttonY, float64(screenWidth/4), buttonH, color.RGBA{R: 255, G: 255, A: 64})

		// Show currently pressed buttons
		y := buttonY + 20
		if g.input.IsKeyPressed(ebiten.KeyLeft) {
			ebitenutil.DebugPrintAt(screen, "LEFT", 10, int(y))
		}
		if g.input.IsKeyPressed(ebiten.KeyRight) {
			ebitenutil.DebugPrintAt(screen, "RIGHT", screenWidth/4+10, int(y))
		}
		if g.input.IsKeyPressed(ebiten.KeyDown) {
			ebitenutil.DebugPrintAt(screen, "DOWN", screenWidth*2/4+10, int(y))
		}
		if g.input.IsKeyPressed(ebiten.KeyEnter) || g.input.IsKeyPressed(ebiten.KeySpace) {
			ebitenutil.DebugPrintAt(screen, "ACTION", screenWidth*3/4+10, int(y))
		}
	}

	// Draw active touch points
	for _, id := range ebiten.TouchIDs() {
		x, y := ebiten.TouchPosition(id)
		ebitenutil.DrawRect(screen, float64(x)-5, float64(y)-5, 10, 10, color.RGBA{R: 255, B: 255, A: 255})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Input Test")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
