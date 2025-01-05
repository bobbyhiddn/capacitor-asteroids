package systems

import (
	"math"
	"math/rand"

	"github.com/bobbyhiddn/ecs-asteroids/components"
	"github.com/bobbyhiddn/ecs-asteroids/ecs"
	"github.com/bobbyhiddn/ecs-asteroids/game"
)

const (
	minSpawnInterval = 2.5   // Minimum time between spawns in seconds
	maxAsteroids     = 8     // Maximum number of asteroids allowed at once
	spawnPadding     = 50    // Distance beyond screen edges where asteroids spawn
	minSpeed         = 50.0  // Minimum asteroid speed
	maxSpeed         = 100.0 // Maximum asteroid speed
)

type AsteroidSpawnerSystem struct {
	world          *ecs.World
	timeSinceSpawn float64
}

func NewAsteroidSpawnerSystem(world *ecs.World) *AsteroidSpawnerSystem {
	return &AsteroidSpawnerSystem{
		world:          world,
		timeSinceSpawn: 0,
	}
}

func (s *AsteroidSpawnerSystem) Update(dt float64) {
	s.timeSinceSpawn += dt

	// Count current asteroids
	asteroidCount := 0
	for _, comp := range s.world.Components["components.Asteroid"] {
		if asteroid, ok := comp.(components.Asteroid); ok {
			// Only count large asteroids for spawn limit
			if asteroid.Size == 2 {
				asteroidCount++
			}
		}
	}

	// Spawn new asteroid if conditions are met
	if s.timeSinceSpawn >= minSpawnInterval && asteroidCount < maxAsteroids {
		s.spawnAsteroid()
		s.timeSinceSpawn = 0
	}
}

func (s *AsteroidSpawnerSystem) spawnAsteroid() {
	// Randomly choose which edge to spawn from (0=top, 1=right, 2=bottom, 3=left)
	edge := rand.Intn(4)
	var x, y float64

	switch edge {
	case 0: // Top
		x = rand.Float64() * float64(game.ScreenWidth)
		y = -spawnPadding
	case 1: // Right
		x = float64(game.ScreenWidth) + spawnPadding
		y = rand.Float64() * float64(game.ScreenHeight)
	case 2: // Bottom
		x = rand.Float64() * float64(game.ScreenWidth)
		y = float64(game.ScreenHeight) + spawnPadding
	case 3: // Left
		x = -spawnPadding
		y = rand.Float64() * float64(game.ScreenHeight)
	}

	// Create a large asteroid
	asteroid := game.CreateAsteroid(s.world, 2)

	// Set its position
	s.world.AddComponent(asteroid, components.Position{X: x, Y: y})

	// Calculate velocity towards screen center with some randomness
	centerX := float64(game.ScreenWidth) / 2
	centerY := float64(game.ScreenHeight) / 2

	// Base angle towards center
	angle := math.Atan2(centerY-y, centerX-x)

	// Add randomness to angle (±45 degrees)
	angle += (rand.Float64() - 0.5) * math.Pi / 2

	// Random speed between min and max
	speed := minSpeed + rand.Float64()*(maxSpeed-minSpeed)

	// Set velocity components
	s.world.AddComponent(asteroid, components.Velocity{
		DX:       math.Cos(angle) * speed,
		DY:       math.Sin(angle) * speed,
		MaxSpeed: maxSpeed,
	})
}
