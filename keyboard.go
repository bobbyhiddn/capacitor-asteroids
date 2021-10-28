package main

import (
	"fmt"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	bullet_timer time.Time = time.Now()
)

func KeyboardHandler() {
	// Forward
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		player_velocity_x = player_velocity_x + player_delta_x*0.02
		player_velocity_y = player_velocity_y + player_delta_y*0.02
		show_thruster = 1
	} else {
		show_thruster = 0
	}

	// Back
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		player_velocity_x = player_velocity_x - player_delta_x*0.02
		player_velocity_y = player_velocity_y - player_delta_y*0.02
		show_thruster = 0
	}

	// Turn left
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		player_angle -= 0.10
		// Reset
		if player_angle <= 0 {
			player_angle = 6.283
		}
		player_delta_x = math.Cos(player_angle) * float64(window_height/160)
		player_delta_y = math.Sin(player_angle) * float64(window_height/160)
	}

	// Turn right
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		player_angle += 0.10
		// Reset
		if player_angle >= 6.283 {
			player_angle = 0
		}
		player_delta_x = math.Cos(player_angle) * float64(window_height/160)
		player_delta_y = math.Sin(player_angle) * float64(window_height/160)
	}

	// Shoot
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		t := time.Now()
		elapsed := t.Sub(bullet_timer)
		fmt.Println(elapsed)
		if elapsed > 500*time.Millisecond {
			bullets = append(bullets, &Bullet{
				pos_x: player_pos_x + math.Cos(player_angle)*float64(window_height/36),
				pos_y: player_pos_y + math.Sin(player_angle)*float64(window_height/36),
				angle: player_angle,
			})
		} else {
			return
		}
		bullet_timer = time.Now()
	}

	if ebiten.IsKeyPressed(ebiten.KeyP) {
		key_states[ebiten.KeyP]++
	} else {
		key_states[ebiten.KeyP] = 0
	}
}
