package systems

import (
	"math"

	"github.com/bobbyhiddn/ecs-asteroids/components"
	"github.com/bobbyhiddn/ecs-asteroids/ecs"
	"github.com/bobbyhiddn/ecs-asteroids/game"
)

const (
	thrustForce   = 150.0
	rotationSpeed = 0.1
)

type PlayerSystem struct {
	world *ecs.World
}

func NewPlayerSystem(world *ecs.World) *PlayerSystem {
	return &PlayerSystem{world: world}
}

func (s *PlayerSystem) Update(dt float64) {
	players := s.world.Components["components.Player"]
	inputs := s.world.Components["components.Input"]
	velocities := s.world.Components["components.Velocity"]
	rotations := s.world.Components["components.Rotation"]
	positions := s.world.Components["components.Position"]

	for id := range players {
		if input, ok := inputs[id].(components.Input); ok {
			// Handle rotation
			if rot, ok := rotations[id].(components.Rotation); ok {
				rot.Angle += float64(input.Rotate) * rotationSpeed
				s.world.AddComponent(id, rot)
			}

			// Handle thrust
			if vel, ok := velocities[id].(components.Velocity); ok {
				if input.Forward {
					// Update player velocity based on current rotation
					if rot, ok := rotations[id].(components.Rotation); ok {
						// Apply thrust force scaled by dt
						vel.DX += math.Cos(rot.Angle) * thrustForce * dt
						vel.DY += math.Sin(rot.Angle) * thrustForce * dt
					}
				}

				// Apply velocity limits
				speed := math.Sqrt(vel.DX*vel.DX + vel.DY*vel.DY)
				if speed > vel.MaxSpeed {
					scale := vel.MaxSpeed / speed
					vel.DX *= scale
					vel.DY *= scale
				}

				s.world.AddComponent(id, vel)
			}

			// Update thruster visibility
			player := players[id].(components.Player)
			player.IsThrusting = input.Forward
			s.world.AddComponent(id, player)

			// Handle shooting
			if input.Shoot {
				if pos, ok := positions[id].(components.Position); ok {
					if rot, ok := rotations[id].(components.Rotation); ok {
						// Create bullet at ship's position with ship's rotation
						game.CreateBullet(s.world, pos.X, pos.Y, rot.Angle, id)
					}
				}
				// Reset shoot flag
				input.Shoot = false
				s.world.AddComponent(id, input)
			}
		}
	}
}
