package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game implements ebiten.Game interface.
type Game struct{}

var (
	//keyStates = map[ebiten.Key]int{}
	player_pos_x      float64 = 0
	player_pos_y      float64 = 0
	player_delta_x    float64 = 0
	player_delta_y    float64 = 0
	player_angle      float64 = 0
)

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	keyboard_handling()
	// Write your game's logical update.
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {

	// TODO: move info to toggled debug menu
	fmt.Printf("player_pos_x: %f\n", player_pos_x)
	fmt.Printf("player_pos_y: %f\n", player_pos_y)
	fmt.Printf("player_delta_x: %f\n", player_delta_x)
	fmt.Printf("player_delta_y: %f\n", player_delta_y)
	fmt.Printf("player_angle: %f\n", player_angle)

	ebitenutil.DrawRect(screen, float64(player_pos_x), float64(player_pos_y), 12, 12, color.White)

}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	game := &Game{}
	player_delta_x = math.Cos(player_angle) * 2
	player_delta_y = math.Sin(player_angle) * 2
	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Your game's title")
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
			log.Fatal(err)
	}
}