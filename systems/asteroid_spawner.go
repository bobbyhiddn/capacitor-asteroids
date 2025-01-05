package systems

import (
	"math"
	"math/rand"
	"time"

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
	world        *ecs.World
	screen       *game.Screen
	lastSpawnTime time.Time
}

func NewAsteroidSpawnerSystem(world *ecs.World) *AsteroidSpawnerSystem {
	return &AsteroidSpawnerSystem{
		world:        world,
		screen:       game.NewScreen(),
		lastSpawnTime: time.Now(),
	}
}

func (s *AsteroidSpawnerSystem) Update(dt float64) {
	currentTime := time.Now()
	elapsedTime := currentTime.Sub(s.lastSpawnTime).Seconds()

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
	if elapsedTime >= minSpawnInterval && asteroidCount < maxAsteroids {
		s.spawnAsteroid()
		s.lastSpawnTime = currentTime
	}
}

func (s *AsteroidSpawnerSystem) spawnAsteroid() {
	// Randomly choose a side of the screen to spawn from
	side := rand.Intn(4)
	var x, y float64

	switch side {
	case 0: // Top
		x = float64(rand.Intn(s.screen.Width()))
		y = 0
	case 1: // Right
		x = float64(s.screen.Width())
		y = float64(rand.Intn(s.screen.Height()))
	case 2: // Bottom
		x = float64(rand.Intn(s.screen.Width()))
		y = float64(s.screen.Height())
	case 3: // Left
		x = 0
		y = float64(rand.Intn(s.screen.Height()))
	}

	// Create asteroid at random size (0-2)
	size := rand.Intn(3)
	asteroid := game.CreateAsteroid(s.world, size)

	// Set its position
	s.world.AddComponent(asteroid, components.Position{X: x, Y: y})

	// Calculate velocity towards screen center with some randomness
	angle := math.Atan2(s.screen.CenterY()-y, s.screen.CenterX()-x)

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
