package systems

import (
	"github.com/bobbyhiddn/ecs-asteroids/components"
	"github.com/bobbyhiddn/ecs-asteroids/ecs"
)

type ExplosionSystem struct {
	world *ecs.World
}

func NewExplosionSystem(world *ecs.World) *ExplosionSystem {
	return &ExplosionSystem{world: world}
}

func (s *ExplosionSystem) Update(dt float64) {
	explosions := s.world.Components["components.Explosion"]

	for id, explosionInterface := range explosions {
		explosion := explosionInterface.(components.Explosion)
		explosion.Age += dt

		if explosion.Age >= explosion.MaxAge {
			s.world.DestroyEntity(id)
			continue
		}

		// Update explosion
		s.world.AddComponent(id, explosion)
	}
}
