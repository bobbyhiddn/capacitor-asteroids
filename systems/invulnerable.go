package systems

import (
	"github.com/bobbyhiddn/ecs-asteroids/components"
	"github.com/bobbyhiddn/ecs-asteroids/ecs"
)

type InvulnerableSystem struct {
	world *ecs.World
}

func NewInvulnerableSystem(world *ecs.World) *InvulnerableSystem {
	return &InvulnerableSystem{world: world}
}

func (s *InvulnerableSystem) Update(dt float64) {
	invulnerables := s.world.Components["components.Invulnerable"]
	renderables := s.world.Components["components.Renderable"]

	for id, invulnerableInterface := range invulnerables {
		invulnerable := invulnerableInterface.(components.Invulnerable)

		// Update timer
		invulnerable.Timer -= dt

		// Handle blinking effect
		if renderable, ok := renderables[id].(components.Renderable); ok {
			// Blink 4 times per second
			renderable.Visible = int(invulnerable.Timer*4)%2 == 0
			s.world.AddComponent(id, renderable)
		}

		// Remove invulnerability when timer expires
		if invulnerable.Timer <= 0 {
			s.world.RemoveComponent(id, "components.Invulnerable")
			// Make sure ship is visible when invulnerability ends
			if renderable, ok := renderables[id].(components.Renderable); ok {
				renderable.Visible = true
				s.world.AddComponent(id, renderable)
			}
			continue
		}

		s.world.AddComponent(id, invulnerable)
	}
}
