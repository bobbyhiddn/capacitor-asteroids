package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{}

type GameConfig struct {
	WindowWidth  int `json:"window_width"`
	WindowHeight int `json:"window_height"`
}

type Asteroid struct {
	pos_x    float64
	pos_y    float64
	angle    float64
	velocity float64
	size     int
}

type Bullet struct {
	pos_x       float64
	pos_y       float64
	angle       float64
	origin_time time.Time
}

var (
	// Window info
	window_height int = 800
	window_width  int = 1000

	// Game states
	key_states    map[ebiten.Key]int = map[ebiten.Key]int{}
	show_debug    int                = 0
	show_thruster int                = 0

	// Entities
	asteroids []*Asteroid
	bullets   []*Bullet

	// Player info
	player_pos_x        float64 = 0
	player_pos_y        float64 = 0
	player_delta_x      float64 = 0
	player_delta_y      float64 = 0
	player_angle        float64 = 0
	player_velocity_x   float64 = 0
	player_velocity_y   float64 = 0
	player_max_velocity float64 = 5

	// Laser info
	laser_velocity   float64       = 5
	laser_fire_speed time.Duration = 750 // In milliseconds
	laser_lifespan   time.Duration = 250 // In milliseconds
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

	var new_bullets []*Bullet
	for _, b := range bullets {
		b.pos_x += laser_velocity * math.Cos(b.angle) * float64(window_height/160)
		b.pos_y += laser_velocity * math.Sin(b.angle) * float64(window_height/160)
		t := time.Now()
		elapsed := t.Sub(b.origin_time)
		if elapsed < laser_lifespan*time.Millisecond {
			new_bullets = append(new_bullets, b)
		}
		if b.pos_x > float64(window_width) {
			b.pos_x = 0
		}
		if b.pos_y > float64(window_height) {
			b.pos_y = 0
		}
		if b.pos_x < 0 {
			b.pos_x = float64(window_width)
		}
		if b.pos_y < 0 {
			b.pos_y = float64(window_height)
		}
	}

	bullets = new_bullets

	player_pos_x = player_pos_x + player_velocity_x
	player_pos_y = player_pos_y + player_velocity_y

	player_velocity_x = player_velocity_x * 0.99
	player_velocity_y = player_velocity_y * 0.99

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Debug info
	str := `{{.fps}}{{.player_pos_x}}{{.player_pos_y}}{{.player_delta_x}}{{.player_delta_y}}{{.player_angle}}{{.player_velocity_x}}{{.player_velocity_y}}{{.bullet_count}}`
	str = strings.Replace(str, "{{.fps}}", fmt.Sprintf("FPS: %d\n", int(ebiten.CurrentFPS())), -1)
	str = strings.Replace(str, "{{.player_pos_x}}", fmt.Sprintf("player_pos_x: %f\n", player_pos_x), -1)
	str = strings.Replace(str, "{{.player_pos_y}}", fmt.Sprintf("player_pos_y: %f\n", player_pos_y), -1)
	str = strings.Replace(str, "{{.player_delta_x}}", fmt.Sprintf("player_delta_x: %f\n", player_delta_x), -1)
	str = strings.Replace(str, "{{.player_delta_y}}", fmt.Sprintf("player_delta_y: %f\n", player_delta_y), -1)
	str = strings.Replace(str, "{{.player_angle}}", fmt.Sprintf("player_angle: %f\n", player_angle), -1)
	str = strings.Replace(str, "{{.player_velocity_x}}", fmt.Sprintf("player_velocity_x: %f\n", player_velocity_x), -1)
	str = strings.Replace(str, "{{.player_velocity_y}}", fmt.Sprintf("player_velocity_y: %f\n", player_velocity_y), -1)
	str = strings.Replace(str, "{{.bullet_count}}", fmt.Sprintf("bullet_count: %d\n", len(bullets)), -1)

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
	return key_states[key] == 1
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
