package systems

import (
	"github.com/bobbyhiddn/ecs-asteroids/components"
	"github.com/bobbyhiddn/ecs-asteroids/ecs"
	"github.com/bobbyhiddn/ecs-asteroids/game"
)

const (
	padding = 50 // Distance beyond screen edges before destroying entities
)

type MovementSystem struct {
	world *ecs.World
}

func NewMovementSystem(world *ecs.World) *MovementSystem {
	return &MovementSystem{world: world}
}

func (s *MovementSystem) Update(dt float64) {
	positions := s.world.Components["components.Position"]
	velocities := s.world.Components["components.Velocity"]
	rotations := s.world.Components["components.Rotation"]

	for id := range positions {
		pos := positions[id].(components.Position)
		positionUpdated := false

		// Update position based on velocity if entity has one
		if vel, ok := velocities[id].(components.Velocity); ok {
			pos.X += vel.DX * dt
			pos.Y += vel.DY * dt
			positionUpdated = true
		}

		// Update rotation if entity has one
		if rot, ok := rotations[id].(components.Rotation); ok {
			rot.Angle += rot.RotationSpeed * dt
			s.world.AddComponent(id, rot)
		}

		// Handle screen wrapping or off-screen destruction
		if _, isPlayer := s.world.Components["components.Player"][id]; isPlayer {
			// Wrap player position around screen edges
			pos.X = wrapCoordinate(pos.X, game.ScreenWidth)
			pos.Y = wrapCoordinate(pos.Y, game.ScreenHeight)
			positionUpdated = true
		} else {
			// For non-player entities, destroy them if they go too far off screen
			if isOffScreen(pos.X, pos.Y) {
				s.world.DestroyEntity(id)
				continue
			}
		}

		// Save the updated position if it changed
		if positionUpdated {
			s.world.AddComponent(id, pos)
		}
	}
}

func wrapCoordinate(pos float64, max float64) float64 {
	if pos < 0 {
		return pos + max
	}
	if pos >= max {
		return pos - max
	}
	return pos
}

func isOffScreen(x, y float64) bool {
	return x < -padding || x > game.ScreenWidth+padding || y < -padding || y > game.ScreenHeight+padding
}
