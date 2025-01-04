package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/samuel-pratt/ebiten-asteroids/components"
	"github.com/samuel-pratt/ebiten-asteroids/ecs"
	"github.com/samuel-pratt/ebiten-asteroids/game"
	"math/rand"
)

type InputSystem struct {
	world *ecs.World
}

func NewInputSystem(world *ecs.World) *InputSystem {
	return &InputSystem{world: world}
}

func (s *InputSystem) Update(dt float64) {
	players := s.world.Components["components.Player"]
	inputs := s.world.Components["components.Input"]

	for id, playerInterface := range players {
		player := playerInterface.(components.Player)
		input := inputs[id].(components.Input)

		// Check for game over restart
		if player.IsGameOver {
			// Check if any key is pressed
			if len(inpututil.AppendPressedKeys(nil)) > 0 {
				// Reset player
				player.Lives = 3
				player.Score = 0
				player.IsGameOver = false
				s.world.AddComponent(id, player)

				// Reset position
				if pos, ok := s.world.Components["components.Position"][id].(components.Position); ok {
					pos.X = float64(game.ScreenWidth / 2)
					pos.Y = float64(game.ScreenHeight / 2)
					s.world.AddComponent(id, pos)
				}

				// Reset velocity
				if vel, ok := s.world.Components["components.Velocity"][id].(components.Velocity); ok {
					vel.DX = 0
					vel.DY = 0
					s.world.AddComponent(id, vel)
				}

				// Reset rotation
				if rot, ok := s.world.Components["components.Rotation"][id].(components.Rotation); ok {
					rot.Angle = 0
					s.world.AddComponent(id, rot)
				}

				// Make ship visible again
				if renderable, ok := s.world.Components["components.Renderable"][id].(components.Renderable); ok {
					renderable.Visible = true
					s.world.AddComponent(id, renderable)
				}

				// Clear all asteroids
				for asteroidID := range s.world.Components["components.Asteroid"] {
					s.world.DestroyEntity(asteroidID)
				}

				// Create initial asteroids
				for i := 0; i < 4; i++ {
					game.CreateAsteroid(s.world, rand.Intn(3))
				}

				continue
			}
		}

		// Normal input handling
		input.Rotate = 0
		if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			input.Rotate = -1
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			input.Rotate = 1
		}

		input.Forward = ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW)
		input.Shoot = inpututil.IsKeyJustPressed(ebiten.KeySpace)

		s.world.AddComponent(id, input)
	}
}
