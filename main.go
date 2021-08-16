package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{}

type GameConfig struct {
	WindowWidth  int `json:"window_width"`
	WindowHeight int `json:"window_height"`
}

var (
	window_height int
	window_width  int

	keyStates              = map[ebiten.Key]int{}
	player_pos_x   float64 = 0
	player_pos_y   float64 = 0
	player_delta_x float64 = 0
	player_delta_y float64 = 0
	player_angle   float64 = 0

	show_debug int = 0
)

func (g *Game) Update() error {
	KeyboardHandler()

	// Loop player back to other side of screen
	if player_pos_x > float64(window_width) {
		player_pos_x = 0
	}
	if player_pos_y > float64(window_height) {
		player_pos_y = 0
	}
	if player_pos_x < 0 {
		player_pos_x = float64(window_width)
	}
	if player_pos_y < 0 {
		player_pos_y = float64(window_height)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Show debug info
	str := `{{.player_pos_x}}{{.player_pos_y}}{{.player_delta_x}}{{.player_delta_y}}{{.player_angle}}`
	str = strings.Replace(str, "{{.player_pos_x}}", fmt.Sprintf("player_pos_x:   %f\n", player_pos_x), -1)
	str = strings.Replace(str, "{{.player_pos_y}}", fmt.Sprintf("player_pos_y:   %f\n", player_pos_y), -1)
	str = strings.Replace(str, "{{.player_delta_x}}", fmt.Sprintf("player_delta_x: %f\n", player_delta_x), -1)
	str = strings.Replace(str, "{{.player_delta_y}}", fmt.Sprintf("player_delta_y: %f\n", player_delta_y), -1)
	str = strings.Replace(str, "{{.player_angle}}", fmt.Sprintf("player_angle:   %f\n", player_angle), -1)
	if show_debug == 1 {
		ebitenutil.DebugPrint(screen, str)
		// Show player rotation point
		ebitenutil.DrawRect(screen, player_pos_x, player_pos_y, 1, 1, color.White)
	}

	DrawSpaceship(screen, player_pos_x, player_pos_y, player_angle)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func IsKeyTriggered(key ebiten.Key) bool {
	return keyStates[key] == 1
}

func DrawSpaceship(screen *ebiten.Image, centerX, centerY, angle float64) {
	ship_point_x := centerX + math.Cos(angle) * float64(window_height/36)
	ship_point_y := centerY + math.Sin(angle) * float64(window_height/36)

	rear_line_center_x := centerX - math.Cos(angle) * float64(window_height/114)
	rear_line_center_y := centerY - math.Sin(angle) * float64(window_height/114)

	line_one_x := math.Cos(angle-math.Pi/8) * float64(window_height/20)
	line_one_y := math.Sin(angle-math.Pi/8) * float64(window_height/20)
	line_two_x := math.Cos(angle+math.Pi/8) * float64(window_height/20)
	line_two_y := math.Sin(angle+math.Pi/8) * float64(window_height/20)
	line_three_x := math.Cos(angle-math.Pi/2) * float64(window_height/62)
	line_three_y := math.Sin(angle-math.Pi/2) * float64(window_height/62)
	line_four_x := math.Cos(angle+math.Pi/2) * float64(window_height/62)
	line_four_y := math.Sin(angle+math.Pi/2) * float64(window_height/62)

	ebitenutil.DrawLine(screen, ship_point_x, ship_point_y, ship_point_x-line_one_x, ship_point_y-line_one_y, color.White)
	ebitenutil.DrawLine(screen, ship_point_x, ship_point_y, ship_point_x-line_two_x, ship_point_y-line_two_y, color.White)
	ebitenutil.DrawLine(screen, rear_line_center_x, rear_line_center_y, rear_line_center_x-line_three_x, rear_line_center_y-line_three_y, color.White)
	ebitenutil.DrawLine(screen, rear_line_center_x, rear_line_center_y, rear_line_center_x-line_four_x, rear_line_center_y-line_four_y, color.White)
}

func main() {
	HandleGameConfig()
	
	// Init work
	player_pos_x = float64(window_width)/2
	player_pos_y = float64(window_height)/2
	player_delta_x = math.Cos(player_angle) * float64(window_height/160)
	player_delta_y = math.Sin(player_angle) * float64(window_height/160)

	game := &Game{}

	ebiten.SetWindowSize(window_width, window_height)
	ebiten.SetWindowTitle("Ebiten Asteroids")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
