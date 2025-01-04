package systems

import (
	"github.com/samuel-pratt/ebiten-asteroids/components"
	"github.com/samuel-pratt/ebiten-asteroids/ecs"
	"github.com/samuel-pratt/ebiten-asteroids/game"
	"math"
	"math/rand"
	"fmt"
)

type CollisionSystem struct {
	world *ecs.World
}

func NewCollisionSystem(world *ecs.World) *CollisionSystem {
	return &CollisionSystem{world: world}
}

func (s *CollisionSystem) Update(dt float64) {
	positions := s.world.Components["components.Position"]
	colliders := s.world.Components["components.Collider"]
	asteroids := s.world.Components["components.Asteroid"]

	// Check all pairs of entities with colliders
	for id1, collider1Interface := range colliders {
		collider1 := collider1Interface.(components.Collider)
		pos1 := positions[id1].(components.Position)

		for id2, collider2Interface := range colliders {
			if id1 >= id2 {
				continue // Skip self and already checked pairs
			}

			collider2 := collider2Interface.(components.Collider)
			pos2 := positions[id2].(components.Position)

			// Check collision
			dx := pos2.X - pos1.X
			dy := pos2.Y - pos1.Y
			distance := math.Sqrt(dx*dx + dy*dy)
			minDistance := collider1.Radius + collider2.Radius

			if distance < minDistance {
				// Handle asteroid-asteroid collisions
				_, isAsteroid1 := asteroids[id1]
				_, isAsteroid2 := asteroids[id2]
				if isAsteroid1 && isAsteroid2 {
					s.handleAsteroidCollision(id1, id2, pos1, pos2, collider1, collider2)
					continue
				}

				// Handle ship-asteroid collisions
				if _, isShip := s.world.Components["components.Player"][id1]; isShip {
					s.handleShipHit(id1)
				} else if _, isShip := s.world.Components["components.Player"][id2]; isShip {
					s.handleShipHit(id2)
				}

				// Handle bullet-asteroid collisions
				if _, isBullet := s.world.Components["components.Bullet"][id1]; isBullet {
					if shooter := s.findShooter(id1); shooter != 0 {
						s.awardPoints(shooter, id2)
					}
					s.handleAsteroidHit(id2)
					s.world.DestroyEntity(id1)
				} else if _, isBullet := s.world.Components["components.Bullet"][id2]; isBullet {
					if shooter := s.findShooter(id2); shooter != 0 {
						s.awardPoints(shooter, id1)
					}
					s.handleAsteroidHit(id1)
					s.world.DestroyEntity(id2)
				}
			}
		}
	}
}

func (s *CollisionSystem) handleAsteroidCollision(id1, id2 ecs.EntityID, pos1, pos2 components.Position, col1, col2 components.Collider) {
	// Get velocities
	vel1 := s.world.Components["components.Velocity"][id1].(components.Velocity)
	vel2 := s.world.Components["components.Velocity"][id2].(components.Velocity)

	fmt.Printf("Before collision - Asteroid 1: vel=(%f, %f), pos=(%f, %f)\n", vel1.DX, vel1.DY, pos1.X, pos1.Y)
	fmt.Printf("Before collision - Asteroid 2: vel=(%f, %f), pos=(%f, %f)\n", vel2.DX, vel2.DY, pos2.X, pos2.Y)

	// Calculate normal vector
	dx := pos2.X - pos1.X
	dy := pos2.Y - pos1.Y
	dist := math.Sqrt(dx*dx + dy*dy)
	nx := dx / dist
	ny := dy / dist

	// Calculate relative velocity
	dvx := vel2.DX - vel1.DX
	dvy := vel2.DY - vel1.DY

	// Calculate relative velocity along normal
	velAlongNormal := dvx*nx + dvy*ny

	// If asteroids are moving apart, don't bounce
	if velAlongNormal > 0 {
		return
	}

	// Elastic collision with some energy loss (90% elastic)
	restitution := 0.9
	impulse := 2.0 * velAlongNormal * restitution

	// Apply impulse
	vel1New := components.Velocity{
		DX:       vel1.DX + impulse*nx,
		DY:       vel1.DY + impulse*ny,
		MaxSpeed: vel1.MaxSpeed,
	}
	
	vel2New := components.Velocity{
		DX:       vel2.DX - impulse*nx,
		DY:       vel2.DY - impulse*ny,
		MaxSpeed: vel2.MaxSpeed,
	}

	// Add some random variation to prevent asteroids from getting stuck
	randomAngle1 := rand.Float64() * math.Pi * 2
	randomAngle2 := rand.Float64() * math.Pi * 2
	randomForce := 20.0

	vel1New.DX += math.Cos(randomAngle1) * randomForce
	vel1New.DY += math.Sin(randomAngle1) * randomForce
	vel2New.DX += math.Cos(randomAngle2) * randomForce
	vel2New.DY += math.Sin(randomAngle2) * randomForce

	// Enforce minimum speed
	minSpeed := 100.0
	for _, vel := range []*components.Velocity{&vel1New, &vel2New} {
		speed := math.Sqrt(vel.DX*vel.DX + vel.DY*vel.DY)
		if speed < minSpeed {
			scale := minSpeed / speed
			vel.DX *= scale
			vel.DY *= scale
		} else if speed > vel.MaxSpeed {
			scale := vel.MaxSpeed / speed
			vel.DX *= scale
			vel.DY *= scale
		}
	}

	// Update velocities
	s.world.AddComponent(id1, vel1New)
	s.world.AddComponent(id2, vel2New)

	// Separate asteroids to prevent sticking
	overlap := (col1.Radius + col2.Radius) - dist
	if overlap > 0 {
		separation := overlap / 2
		pos1.X -= nx * separation
		pos1.Y -= ny * separation
		pos2.X += nx * separation
		pos2.Y += ny * separation

		s.world.AddComponent(id1, pos1)
		s.world.AddComponent(id2, pos2)
	}

	fmt.Printf("After collision - Asteroid 1: vel=(%f, %f), pos=(%f, %f)\n", vel1New.DX, vel1New.DY, pos1.X, pos1.Y)
	fmt.Printf("After collision - Asteroid 2: vel=(%f, %f), pos=(%f, %f)\n", vel2New.DX, vel2New.DY, pos2.X, pos2.Y)
}

func (s *CollisionSystem) handleShipHit(shipID ecs.EntityID) {
	// Check if ship is invulnerable
	if _, ok := s.world.Components["components.Invulnerable"][shipID]; ok {
		return // Skip collision if ship is invulnerable
	}

	player, ok := s.world.Components["components.Player"][shipID].(components.Player)
	if !ok {
		return
	}

	// Get current position for explosion and respawn
	pos := s.world.Components["components.Position"][shipID].(components.Position)

	// Create explosion at ship's position
	game.CreateExplosion(s.world, pos.X, pos.Y, 30.0) // Size matches ship roughly

	// Reduce lives
	player.Lives--

	// Check for game over
	if player.Lives <= 0 {
		player.IsGameOver = true
		player.Lives = 0 // Ensure it doesn't go negative
	}

	// Update player component
	s.world.AddComponent(shipID, player)

	// Reset velocity
	s.world.AddComponent(shipID, components.Velocity{MaxSpeed: 400})

	// Reset rotation
	s.world.AddComponent(shipID, components.Rotation{})

	// If not game over, move ship back to center and make invulnerable
	if !player.IsGameOver {
		pos.X = float64(game.ScreenWidth / 2)
		pos.Y = float64(game.ScreenHeight / 2)
		s.world.AddComponent(shipID, pos)

		// Add invulnerability
		s.world.AddComponent(shipID, components.Invulnerable{
			Duration: 3.0, // 3 seconds of invulnerability
			Timer: 3.0,
		})
	} else {
		// Make ship invisible on game over
		if renderable, ok := s.world.Components["components.Renderable"][shipID].(components.Renderable); ok {
			renderable.Visible = false
			s.world.AddComponent(shipID, renderable)
		}
	}
}

func (s *CollisionSystem) handleAsteroidHit(asteroidID ecs.EntityID) {
	asteroid, ok := s.world.Components["components.Asteroid"][asteroidID].(components.Asteroid)
	if !ok {
		return
	}

	pos := s.world.Components["components.Position"][asteroidID].(components.Position)

	// Destroy the hit asteroid
	s.world.DestroyEntity(asteroidID)

	// If it wasn't the smallest size, spawn two smaller asteroids
	if asteroid.Size > 0 {
		for i := 0; i < 2; i++ {
			newAsteroid := game.CreateAsteroid(s.world, asteroid.Size-1)
			if newPos, ok := s.world.Components["components.Position"][newAsteroid].(components.Position); ok {
				newPos.X = pos.X
				newPos.Y = pos.Y
				s.world.AddComponent(newAsteroid, newPos)
			}
		}
	}
}

func (s *CollisionSystem) findShooter(bulletID ecs.EntityID) ecs.EntityID {
	// In this simple implementation, we assume there's only one player
	// In a multiplayer game, you'd want to store the shooter's ID with the bullet
	for id := range s.world.Components["components.Player"] {
		return id
	}
	return 0
}

func (s *CollisionSystem) awardPoints(playerID, asteroidID ecs.EntityID) {
	player, ok := s.world.Components["components.Player"][playerID].(components.Player)
	if !ok {
		return
	}

	asteroid, ok := s.world.Components["components.Asteroid"][asteroidID].(components.Asteroid)
	if !ok {
		return
	}

	// Award points based on asteroid size
	// Small asteroids are worth more points
	points := 0
	switch asteroid.Size {
	case 0: // Small
		points = 100
	case 1: // Medium
		points = 50
	case 2: // Large
		points = 20
	}

	player.Score += points
	s.world.AddComponent(playerID, player)
}
