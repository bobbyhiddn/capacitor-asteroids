package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func KeyboardHandler() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		player_pos_x = player_pos_x + player_delta_x
		player_pos_y = player_pos_y + player_delta_y
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		player_pos_x = player_pos_x - player_delta_x
		player_pos_y = player_pos_y - player_delta_y
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		player_angle -= 0.10
		// Reset
		if player_angle <= 0 {
			player_angle = 6.283
		}
		player_delta_x = math.Cos(player_angle) * float64(window_height/160)
		player_delta_y = math.Sin(player_angle) * float64(window_height/160)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		player_angle += 0.10
		// Reset
		if player_angle >= 6.283 {
			player_angle = 0
		}
		player_delta_x = math.Cos(player_angle) * float64(window_height/160)
		player_delta_y = math.Sin(player_angle) * float64(window_height/160)
	}

	if ebiten.IsKeyPressed(ebiten.KeyP) {
		keyStates[ebiten.KeyP]++
	} else {
		keyStates[ebiten.KeyP] = 0
	}

	if IsKeyTriggered(ebiten.KeyP) {
		if show_debug == 0 {
			show_debug = 1
		} else {
			show_debug = 0
		}
	}
}
