package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
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
		bullets = append(bullets, &Bullet{
			pos_x: player_pos_x + math.Cos(player_angle)*float64(window_height/36),
			pos_y: player_pos_y + math.Sin(player_angle)*float64(window_height/36),
			angle: player_angle,
		})
	}

	if ebiten.IsKeyPressed(ebiten.KeyP) {
		keyStates[ebiten.KeyP]++
	} else {
		keyStates[ebiten.KeyP] = 0
	}
}
