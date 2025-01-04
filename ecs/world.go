package ecs

import (
	"fmt"
	"image/color"
)

type EntityID int

type System interface {
	Update(dt float64)
}

type World struct {
	nextEntityID    EntityID
	Components      map[string]map[EntityID]interface{}
	systems         []System
	entities        map[EntityID]bool
	BackgroundColor color.Color
}

func NewWorld() *World {
	return &World{
		nextEntityID: 1,
		Components: map[string]map[EntityID]interface{}{
			"components.Position":     make(map[EntityID]interface{}),
			"components.Velocity":     make(map[EntityID]interface{}),
			"components.Rotation":     make(map[EntityID]interface{}),
			"components.Renderable":   make(map[EntityID]interface{}),
			"components.Player":       make(map[EntityID]interface{}),
			"components.Input":        make(map[EntityID]interface{}),
			"components.Lifetime":     make(map[EntityID]interface{}),
			"components.Collider":     make(map[EntityID]interface{}),
			"components.Asteroid":     make(map[EntityID]interface{}),
			"components.Explosion":    make(map[EntityID]interface{}),
			"components.Invulnerable": make(map[EntityID]interface{}),
			"components.Bullet":       make(map[EntityID]interface{}),
		},
		systems:         make([]System, 0),
		entities:        make(map[EntityID]bool),
		BackgroundColor: color.Black,
	}
}

func (w *World) AddSystem(system System) {
	w.systems = append(w.systems, system)
}

func (w *World) CreateEntity() EntityID {
	id := w.nextEntityID
	w.nextEntityID++
	w.entities[id] = true
	return id
}

func (w *World) DestroyEntity(id EntityID) {
	if !w.entities[id] {
		return
	}

	// Remove all components for this entity
	for _, components := range w.Components {
		delete(components, id)
	}

	delete(w.entities, id)
}

func (w *World) AddComponent(id EntityID, component interface{}) {
	componentType := fmt.Sprintf("%T", component)
	if _, ok := w.Components[componentType]; !ok {
		w.Components[componentType] = make(map[EntityID]interface{})
	}
	w.Components[componentType][id] = component
}

func (w *World) RemoveComponent(id EntityID, componentType string) {
	if components, ok := w.Components[componentType]; ok {
		delete(components, id)
	}
}

func (w *World) Update(dt float64) {
	for _, system := range w.systems {
		system.Update(dt)
	}
}
