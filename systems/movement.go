package systems

import (
	"github.com/bobbyhiddn/ecs-asteroids/components"
	"github.com/bobbyhiddn/ecs-asteroids/ecs"
	"github.com/bobbyhiddn/ecs-asteroids/game"
)

type MovementSystem struct {
	world  *ecs.World
	screen *game.Screen
}

func NewMovementSystem(world *ecs.World) *MovementSystem {
	return &MovementSystem{
		world:  world,
		screen: game.NewScreen(),
	}
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
			s.wrapPosition(&pos)
			positionUpdated = true
		} else {
			// For non-player entities, destroy them if they go too far off screen
			if s.isOffScreen(pos) {
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

func (s *MovementSystem) wrapPosition(pos *components.Position) {
	width := float64(s.screen.Width())
	height := float64(s.screen.Height())

	if pos.X < 0 {
		pos.X += width
	} else if pos.X >= width {
		pos.X -= width
	}

	if pos.Y < 0 {
		pos.Y += height
	} else if pos.Y >= height {
		pos.Y -= height
	}
}

func (s *MovementSystem) isOffScreen(pos components.Position) bool {
	return pos.X < 0 || pos.X > float64(s.screen.Width()) || pos.Y < 0 || pos.Y > float64(s.screen.Height())
}
