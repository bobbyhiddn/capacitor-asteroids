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

type Asteroid struct {
	pos_x, pos_y, angle, velocity float64
	size                          int
}

var (
	window_height int
	window_width  int

	keyStates map[ebiten.Key]int = map[ebiten.Key]int{}
	asteroids []*Asteroid

	player_pos_x   float64 = 0
	player_pos_y   float64 = 0
	player_delta_x float64 = 0
	player_delta_y float64 = 0
	player_angle   float64 = 0

	show_debug    int = 0
	show_thruster int = 0
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

	for _, a := range asteroids {
		a.pos_x += a.velocity * math.Cos(a.angle) * float64(window_height/160)
		a.pos_y += a.velocity * math.Sin(a.angle) * float64(window_height/160)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Debug info
	str := `{{.player_pos_x}}{{.player_pos_y}}{{.player_delta_x}}{{.player_delta_y}}{{.player_angle}}`
	str = strings.Replace(str, "{{.player_pos_x}}", fmt.Sprintf("player_pos_x:   %f\n", player_pos_x), -1)
	str = strings.Replace(str, "{{.player_pos_y}}", fmt.Sprintf("player_pos_y:   %f\n", player_pos_y), -1)
	str = strings.Replace(str, "{{.player_delta_x}}", fmt.Sprintf("player_delta_x: %f\n", player_delta_x), -1)
	str = strings.Replace(str, "{{.player_delta_y}}", fmt.Sprintf("player_delta_y: %f\n", player_delta_y), -1)
	str = strings.Replace(str, "{{.player_angle}}", fmt.Sprintf("player_angle:   %f\n", player_angle), -1)
	if show_debug == 1 {
		// Show debug info
		ebitenutil.DebugPrint(screen, str)
		// Show player point of rotation
		ebitenutil.DrawRect(screen, player_pos_x, player_pos_y, 1, 1, color.White)
	}

	DrawSpaceship(screen, player_pos_x, player_pos_y, player_angle)

	if show_thruster == 1 {
		DrawThrusters(screen, player_pos_x, player_pos_y, player_angle)
	}

	for _, a := range asteroids {
		if a.size == 0 {
			DrawAsteroidSmall(screen, a.pos_x, a.pos_y, a.angle)
		} else if a.size == 1 {
			DrawAsteroidMedium(screen, a.pos_x, a.pos_y, a.angle)
		} else if a.size == 2 {
			DrawAsteroidLarge(screen, a.pos_x, a.pos_y, a.angle)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func IsKeyTriggered(key ebiten.Key) bool {
	return keyStates[key] == 1
}

func main() {
	HandleGameConfig()

	// Init work
	player_pos_x = float64(window_width) / 2
	player_pos_y = float64(window_height) / 2
	player_delta_x = math.Cos(player_angle) * float64(window_height/160)
	player_delta_y = math.Sin(player_angle) * float64(window_height/160)

	// Test asteroid
	asteroids = append(asteroids, &Asteroid{
		pos_x:    30,
		pos_y:    30,
		angle:    1,
		velocity: 0.5,
		size:     2,
	})

	game := &Game{}

	ebiten.SetWindowSize(window_width, window_height)
	ebiten.SetWindowTitle("Ebiten Asteroids")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
