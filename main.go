package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{}

var (
	//keyStates = map[ebiten.Key]int{}
	// TODO: replace with a config/game.json config file
	window_height		 int     = 480
	window_width		 int     = 640

	player_pos_x      float64 = float64(window_width/4) // divided by 4 because layout is divided by 2
	player_pos_y      float64 = float64(window_height/4)
	player_delta_x    float64 = 0
	player_delta_y    float64 = 0
	player_angle      float64 = 0
)

func (g *Game) Update() error {
	keyboard_handling()

	// Loop player back to other side of screen
	if (player_pos_x > float64(window_width/2)) {
		player_pos_x = 0
	}
	if (player_pos_y > float64(window_height/2)) {
		player_pos_y = 0
	}
	if (player_pos_x < 0) {
		player_pos_x = float64(window_width/2)
	}
	if (player_pos_y < 0) {
		player_pos_y = float64(window_height/2)
	}
	
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// TODO: move info to toggled debug menu
	fmt.Printf("player_pos_x: %f\n", player_pos_x)
	fmt.Printf("player_pos_y: %f\n", player_pos_y)
	fmt.Printf("player_delta_x: %f\n", player_delta_x)
	fmt.Printf("player_delta_y: %f\n", player_delta_y)
	fmt.Printf("player_angle: %f\n", player_angle)

	// TODO: replace line with player, player being a collection of lines in shape of ship
	ebitenutil.DrawLine(screen, float64(player_pos_x), float64(player_pos_y), float64(player_pos_x) + player_delta_x * 4, float64(player_pos_y) +player_delta_y * 4, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth/2, outsideHeight/2
}

func main() {
	// Init work
	player_delta_x = math.Cos(player_angle) * 2
	player_delta_y = math.Sin(player_angle) * 2

	game := &Game{}

	ebiten.SetWindowSize(window_width, window_height)
	ebiten.SetWindowTitle("Ebiten Asteroids")
	
	if err := ebiten.RunGame(game); err != nil {
			log.Fatal(err)
	}
}