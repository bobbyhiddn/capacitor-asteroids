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

type Bullet struct {
	pos_x, pos_y, angle float64
}

var (
	// Window info
	window_height int = 800
	window_width  int = 1000

	// Game states
	keyStates     map[ebiten.Key]int = map[ebiten.Key]int{}
	show_debug    int                = 0
	show_thruster int                = 0

	// Entiries
	asteroids []*Asteroid
	bullets   []*Bullet

	// Player info
	player_pos_x      float64 = 0
	player_pos_y      float64 = 0
	player_delta_x    float64 = 0
	player_delta_y    float64 = 0
	player_angle      float64 = 0
	player_velocity_x float64 = 0
	player_velocity_y float64 = 0

	// Constants
	bullet_velocity float64 = 2
)

func (g *Game) Update() error {
	KeyboardHandler()

	if IsKeyTriggered(ebiten.KeyP) {
		if show_debug == 0 {
			show_debug = 1
		} else {
			show_debug = 0
		}
	}

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

	for _, b := range bullets {
		b.pos_x += bullet_velocity * math.Cos(b.angle) * float64(window_height/160)
		b.pos_y += bullet_velocity * math.Sin(b.angle) * float64(window_height/160)
	}

	player_pos_x = player_pos_x + player_velocity_x
	player_pos_y = player_pos_y + player_velocity_y

	player_velocity_x = player_velocity_x * 0.99
	player_velocity_y = player_velocity_y * 0.99

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Debug info
	str := `{{.fps}}{{.player_pos_x}}{{.player_pos_y}}{{.player_delta_x}}{{.player_delta_y}}{{.player_angle}}{{.player_velocity_x}}{{.player_velocity_y}}`
	str = strings.Replace(str, "{{.fps}}", fmt.Sprintf("FPS: %d\n", int(ebiten.CurrentFPS())), -1)
	str = strings.Replace(str, "{{.player_pos_x}}", fmt.Sprintf("player_pos_x: %f\n", player_pos_x), -1)
	str = strings.Replace(str, "{{.player_pos_y}}", fmt.Sprintf("player_pos_y: %f\n", player_pos_y), -1)
	str = strings.Replace(str, "{{.player_delta_x}}", fmt.Sprintf("player_delta_x: %f\n", player_delta_x), -1)
	str = strings.Replace(str, "{{.player_delta_y}}", fmt.Sprintf("player_delta_y: %f\n", player_delta_y), -1)
	str = strings.Replace(str, "{{.player_angle}}", fmt.Sprintf("player_angle: %f\n", player_angle), -1)
	str = strings.Replace(str, "{{.player_velocity_x}}", fmt.Sprintf("player_velocity_x: %f\n", player_velocity_x), -1)
	str = strings.Replace(str, "{{.player_velocity_y}}", fmt.Sprintf("player_velocity_y: %f\n", player_velocity_y), -1)

	if show_debug == 1 {
		// Show debug info
		ebitenutil.DebugPrint(screen, str)
		// Show player point of rotation
		ebitenutil.DrawRect(screen, player_pos_x, player_pos_y, 1, 1, color.White)
	}

	if show_thruster == 1 {
		DrawThrusters(screen, player_pos_x, player_pos_y, player_angle)
	}

	DrawSpaceship(screen, player_pos_x, player_pos_y, player_angle)

	for _, a := range asteroids {
		if a.size == 0 {
			DrawAsteroidSmall(screen, a.pos_x, a.pos_y, a.angle)
		} else if a.size == 1 {
			DrawAsteroidMedium(screen, a.pos_x, a.pos_y, a.angle)
		} else if a.size == 2 {
			DrawAsteroidLarge(screen, a.pos_x, a.pos_y, a.angle)
		}
	}

	for _, b := range bullets {
		DrawBullet(screen, b.pos_x, b.pos_y, b.angle)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func IsKeyTriggered(key ebiten.Key) bool {
	return keyStates[key] == 1
}

func main() {
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
