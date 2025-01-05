package systems

import (
	"fmt"
	"math"
	"math/rand"
	"github.com/samuel-pratt/ebiten-asteroids/components"
	"github.com/samuel-pratt/ebiten-asteroids/ecs"
	"github.com/samuel-pratt/ebiten-asteroids/game"
)

type CollisionSystem struct {
	world *ecs.World
}

func NewCollisionSystem(world *ecs.World) *CollisionSystem {
	return &CollisionSystem{world: world}
}

func (s *CollisionSystem) Update(dt float64) {
	colliders := s.world.Components["components.Collider"]
	positions := s.world.Components["components.Position"]

	// Create a list of all entities with colliders
	var entities []ecs.EntityID
	for id := range colliders {
		entities = append(entities, id)
	}

	// Check each pair of entities for collisions
	for i := 0; i < len(entities); i++ {
		id1 := entities[i]
		pos1, ok1 := positions[id1].(components.Position)
		col1, ok2 := colliders[id1].(components.Collider)
		if !ok1 || !ok2 {
			continue
		}

		for j := i + 1; j < len(entities); j++ {
			id2 := entities[j]
			pos2, ok1 := positions[id2].(components.Position)
			col2, ok2 := colliders[id2].(components.Collider)
			if !ok1 || !ok2 {
				continue
			}

			// Calculate distance between entities
			dx := pos1.X - pos2.X
			dy := pos1.Y - pos2.Y
			distance := math.Sqrt(dx*dx + dy*dy)

			// Check if collision occurred
			if distance < col1.Radius+col2.Radius {
				// Determine collision type
				isAsteroid1 := col1.Type == components.ColliderTypeAsteroid
				isAsteroid2 := col2.Type == components.ColliderTypeAsteroid
				isShip1 := col1.Type == components.ColliderTypeShip
				isShip2 := col2.Type == components.ColliderTypeShip
				isBullet1 := col1.Type == components.ColliderTypeBullet
				isBullet2 := col2.Type == components.ColliderTypeBullet

				// Handle different collision types
				switch {
				// Asteroid-Asteroid collisions
				case isAsteroid1 && isAsteroid2:
					s.handleAsteroidCollision(id1, id2, pos1, pos2, col1, col2)

				// Ship-Asteroid collisions
				case isShip1 && isAsteroid2:
					if !s.isInvulnerable(id1) {
						s.handleShipHit(id1)
					}
				case isShip2 && isAsteroid1:
					if !s.isInvulnerable(id2) {
						s.handleShipHit(id2)
					}

				// Bullet-Asteroid collisions
				case isBullet1 && isAsteroid2:
					if shooter := s.findShooter(id1); shooter != 0 {
						s.awardPoints(shooter, id2)
					}
					s.handleAsteroidHit(id2)
					s.world.DestroyEntity(id1)
				case isBullet2 && isAsteroid1:
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
	vel1Interface := s.world.Components["components.Velocity"][id1]
	vel2Interface := s.world.Components["components.Velocity"][id2]
	
	// Check if velocities exist
	if vel1Interface == nil || vel2Interface == nil {
		// One of the asteroids was probably already destroyed
		return
	}
	
	vel1 := vel1Interface.(components.Velocity)
	vel2 := vel2Interface.(components.Velocity)

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

	// Calculate collision response
	restitution := 0.8 // Coefficient of restitution (0.8 = 80% energy conservation)
	j := -(1.0 + restitution) * velAlongNormal
	j *= 0.5 // Assuming equal masses

	// Calculate new velocities
	vel1New := components.Velocity{
		DX:       vel1.DX - j*nx,
		DY:       vel1.DY - j*ny,
		MaxSpeed: vel1.MaxSpeed,
	}
	
	vel2New := components.Velocity{
		DX:       vel2.DX + j*nx,
		DY:       vel2.DY + j*ny,
		MaxSpeed: vel2.MaxSpeed,
	}

	// Add perpendicular impulse to create more interesting collisions
	// This adds some angular momentum to the collision
	tx := -ny // Tangent vector is perpendicular to normal
	ty := nx
	tangentImpulse := (rand.Float64() - 0.5) * 50.0

	vel1New.DX += tx * tangentImpulse
	vel1New.DY += ty * tangentImpulse
	vel2New.DX -= tx * tangentImpulse
	vel2New.DY -= ty * tangentImpulse

	// Enforce speed constraints for both asteroids
	minSpeed := 100.0
	for _, vel := range []*components.Velocity{&vel1New, &vel2New} {
		speed := math.Sqrt(vel.DX*vel.DX + vel.DY*vel.DY)
		
		// Apply minimum speed
		if speed < minSpeed {
			scale := minSpeed / speed
			vel.DX *= scale
			vel.DY *= scale
			continue // Skip max speed check if we just scaled up
		}
		
		// Apply maximum speed
		if speed > vel.MaxSpeed {
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

		// Update player component
		s.world.AddComponent(shipID, player)

		// Remove all components except Player and Input to effectively disable the ship
		// but keep the player state for restart
		for componentName, components := range s.world.Components {
			if componentName != "components.Player" && componentName != "components.Input" {
				if _, exists := components[shipID]; exists {
					delete(s.world.Components[componentName], shipID)
				}
			}
		}
		return
	}

	// Update player component
	s.world.AddComponent(shipID, player)

	// Reset velocity
	s.world.AddComponent(shipID, components.Velocity{MaxSpeed: 400})

	// Reset rotation
	s.world.AddComponent(shipID, components.Rotation{})

	// Move ship back to center and make invulnerable
	pos.X = float64(game.ScreenWidth / 2)
	pos.Y = float64(game.ScreenHeight / 2)
	s.world.AddComponent(shipID, pos)

	// Add invulnerability
	s.world.AddComponent(shipID, components.Invulnerable{
		Duration: 3.0, // 3 seconds of invulnerability
		Timer: 3.0,
	})
}

func (s *CollisionSystem) handleAsteroidHit(asteroidID ecs.EntityID) {
	asteroid, ok := s.world.Components["components.Asteroid"][asteroidID].(components.Asteroid)
	if !ok {
		return
	}

	pos := s.world.Components["components.Position"][asteroidID].(components.Position)
	vel := s.world.Components["components.Velocity"][asteroidID].(components.Velocity)

	// Create explosion effect
	game.CreateExplosion(s.world, pos.X, pos.Y, float64(20+asteroid.Size*10))

	// Destroy the hit asteroid
	s.world.DestroyEntity(asteroidID)

	// If it wasn't the smallest size, spawn two smaller asteroids
	if asteroid.Size > 0 {
		// Calculate base angle from current velocity
		angle := math.Atan2(vel.DY, vel.DX)

		for i := 0; i < 2; i++ {
			newAsteroid := game.CreateAsteroid(s.world, asteroid.Size-1)
			
			// Position at split point
			s.world.AddComponent(newAsteroid, components.Position{
				X: pos.X,
				Y: pos.Y,
			})

			// Calculate new velocity
			splitAngle := angle
			if i == 1 {
				splitAngle = angle - math.Pi/3
			} else {
				splitAngle = angle + math.Pi/3
			}

			speed := math.Sqrt(vel.DX*vel.DX + vel.DY*vel.DY) * 1.5
			s.world.AddComponent(newAsteroid, components.Velocity{
				DX: math.Cos(splitAngle) * speed,
				DY: math.Sin(splitAngle) * speed,
				MaxSpeed: speed * 1.2,
			})
		}
	}
}

func (s *CollisionSystem) findShooter(bulletID ecs.EntityID) ecs.EntityID {
	if bullet, ok := s.world.Components["components.Bullet"][bulletID].(components.Bullet); ok {
		// Find player with matching ID
		for id := range s.world.Components["components.Player"] {
			if int(id) == bullet.ShooterID {
				return id
			}
		}
	}
	return 0
}

func (s *CollisionSystem) awardPoints(playerID ecs.EntityID, asteroidID ecs.EntityID) {
	player, ok := s.world.Components["components.Player"][playerID].(components.Player)
	if !ok {
		return
	}

	asteroid, ok := s.world.Components["components.Asteroid"][asteroidID].(components.Asteroid)
	if !ok {
		return
	}

	// Award points based on asteroid size
	switch asteroid.Size {
	case 0: // Small
		player.Score += 100
	case 1: // Medium
		player.Score += 50
	case 2: // Large
		player.Score += 20
	}

	// Update player score
	s.world.AddComponent(playerID, player)
}

func (s *CollisionSystem) isInvulnerable(entityID ecs.EntityID) bool {
	_, hasInvulnerable := s.world.Components["components.Invulnerable"][entityID]
	return hasInvulnerable
}
