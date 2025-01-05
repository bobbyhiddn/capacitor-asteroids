package game

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/bobbyhiddn/ecs-asteroids/components"
	"github.com/bobbyhiddn/ecs-asteroids/ecs"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

func CreatePlayerShip(world *ecs.World, x, y float64) ecs.EntityID {
	id := world.CreateEntity()
	fmt.Printf("Creating player ship with ID %d at (%f, %f)\n", id, x, y)

	world.AddComponent(id, components.Position{X: x, Y: y})
	world.AddComponent(id, components.Velocity{MaxSpeed: 400})
	world.AddComponent(id, components.Rotation{})
	world.AddComponent(id, components.Input{})
	world.AddComponent(id, components.Player{Lives: 3})
	world.AddComponent(id, components.Renderable{Type: components.RenderableTypeShip, Visible: true})
	world.AddComponent(id, components.Collider{Radius: 15})
	world.AddComponent(id, components.Invulnerable{Duration: 3.0, Timer: 3.0})

	return id
}

func CreateBullet(world *ecs.World, x, y, angle float64, shooterID ecs.EntityID) ecs.EntityID {
	id := world.CreateEntity()

	speed := 500.0
	world.AddComponent(id, components.Position{X: x, Y: y})
	world.AddComponent(id, components.Velocity{
		DX:       math.Cos(angle) * speed,
		DY:       math.Sin(angle) * speed,
		MaxSpeed: speed,
	})
	world.AddComponent(id, components.Renderable{
		Type:    components.RenderableTypeBullet,
		Scale:   1.0,
		Visible: true,
	})
	world.AddComponent(id, components.Lifetime{
		Created:  time.Now(),
		Duration: 750 * time.Millisecond,
	})
	world.AddComponent(id, components.Collider{
		Radius: 2,
		Type:   components.ColliderTypeBullet,
	})
	world.AddComponent(id, components.Bullet{
		ShooterID: int(shooterID),
	})

	return id
}

func CreateAsteroid(world *ecs.World, size int) ecs.EntityID {
	id := world.CreateEntity()

	var scale float64
	var radius float64
	var maxSpeed float64
	switch size {
	case 0: // Small
		scale = 0.5
		radius = 10
		maxSpeed = 300
	case 1: // Medium
		scale = 1.0
		radius = 20
		maxSpeed = 200
	case 2: // Large
		scale = 2.0
		radius = 40
		maxSpeed = 100
	}

	world.AddComponent(id, components.Position{X: 0, Y: 0})
	world.AddComponent(id, components.Velocity{
		DX:       0,
		DY:       0,
		MaxSpeed: maxSpeed,
	})
	world.AddComponent(id, components.Rotation{
		Angle:         rand.Float64() * math.Pi * 2,
		RotationSpeed: (rand.Float64() - 0.5) * 2,
	})
	world.AddComponent(id, components.Renderable{
		Type:    components.RenderableTypeAsteroid,
		Scale:   scale,
		Visible: true,
	})
	world.AddComponent(id, components.Collider{
		Radius: radius,
		Type:   components.ColliderTypeAsteroid,
	})
	world.AddComponent(id, components.Asteroid{
		Size: size,
	})

	return id
}

func CreateExplosion(world *ecs.World, x, y float64, size float64) ecs.EntityID {
	id := world.CreateEntity()

	world.AddComponent(id, components.Position{X: x, Y: y})
	world.AddComponent(id, components.Renderable{
		Type:    components.RenderableTypeExplosion,
		Scale:   1.0,
		Visible: true,
	})
	world.AddComponent(id, components.Explosion{
		Age:    0,
		MaxAge: 0.5, // Explosion lasts 0.5 seconds
		Radius: size,
		Pieces: 12, // Number of particles in explosion
	})

	return id
}
