package systems

import (
	"math"
	"math/rand"

	"github.com/samuel-pratt/ebiten-asteroids/components"
	"github.com/samuel-pratt/ebiten-asteroids/ecs"
	"github.com/samuel-pratt/ebiten-asteroids/game"
)

const (
	minSpawnInterval = 1.0  // Minimum time between spawns in seconds
	maxAsteroids    = 12    // Maximum number of asteroids allowed at once
	spawnPadding    = 50    // Distance beyond screen edges where asteroids spawn
	minSpeed        = 50.0  // Minimum asteroid speed
	maxSpeed        = 100.0 // Maximum asteroid speed
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
	for range s.world.Components["components.Asteroid"] {
		asteroidCount++
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

	// Create asteroid with random size (prefer larger ones)
	size := rand.Intn(3)
	asteroid := game.CreateAsteroid(s.world, size)

	// Set position
	s.world.AddComponent(asteroid, components.Position{X: x, Y: y})

	// Calculate direction vector towards center of screen with some randomness
	centerX := float64(game.ScreenWidth) / 2
	centerY := float64(game.ScreenHeight) / 2
	
	// Add some randomness to target point
	targetX := centerX + (rand.Float64()*200 - 100)
	targetY := centerY + (rand.Float64()*200 - 100)
	
	// Calculate direction vector
	dx := targetX - x
	dy := targetY - y

	// Normalize the direction vector
	length := math.Sqrt(dx*dx + dy*dy)
	if length != 0 {
		dx /= length
		dy /= length
	}

	// Apply random speed
	speed := minSpeed + rand.Float64()*(maxSpeed-minSpeed)
	dx *= speed
	dy *= speed

	s.world.AddComponent(asteroid, components.Velocity{
		DX: dx,
		DY: dy,
	})
}
