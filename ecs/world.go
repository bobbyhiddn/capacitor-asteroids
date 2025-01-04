package ecs

import (
	"github.com/samuel-pratt/ebiten-asteroids/components"
	"fmt"
	"image/color"
)

type EntityID int

type System interface {
	Update(dt float64)
}

type World struct {
	nextID     EntityID
	Components map[string]map[EntityID]interface{}
	systems    []System
	entities   map[EntityID]bool
	BackgroundColor color.Color
}

func NewWorld() *World {
	return &World{
		nextID: 1,
		Components: map[string]map[EntityID]interface{}{
			"components.Position":     make(map[EntityID]interface{}),
			"components.Velocity":     make(map[EntityID]interface{}),
			"components.Rotation":     make(map[EntityID]interface{}),
			"components.Input":        make(map[EntityID]interface{}),
			"components.Player":       make(map[EntityID]interface{}),
			"components.Renderable":   make(map[EntityID]interface{}),
			"components.Lifetime":     make(map[EntityID]interface{}),
			"components.Collider":     make(map[EntityID]interface{}),
			"components.Asteroid":     make(map[EntityID]interface{}),
			"components.Explosion":    make(map[EntityID]interface{}),
			"components.Invulnerable": make(map[EntityID]interface{}),
		},
		systems:    make([]System, 0),
		entities:   make(map[EntityID]bool),
		BackgroundColor: color.Black,
	}
}

func (w *World) AddSystem(system System) {
	w.systems = append(w.systems, system)
}

func (w *World) CreateEntity() EntityID {
	id := w.nextID
	w.nextID++
	w.entities[id] = true
	return id
}

func (w *World) DestroyEntity(id EntityID) {
	// Remove all components for this entity
	for _, components := range w.Components {
		delete(components, id)
	}
	// Remove the entity itself
	delete(w.entities, id)
}

func (w *World) AddComponent(id EntityID, component interface{}) {
	componentType := getComponentType(component)
	if w.Components[componentType] == nil {
		w.Components[componentType] = make(map[EntityID]interface{})
	}
	w.Components[componentType][id] = component
}

func (w *World) RemoveComponent(entityID EntityID, component interface{}) {
	componentType := getComponentType(component)
	if _, ok := w.Components[componentType]; ok {
		delete(w.Components[componentType], entityID)
	}
}

func (w *World) GetComponent(id EntityID, componentType string) (interface{}, bool) {
	if components, ok := w.Components[componentType]; ok {
		if component, ok := components[id]; ok {
			return component, true
		}
	}
	return nil, false
}

func (w *World) Update(dt float64) {
	for _, system := range w.systems {
		system.Update(dt)
	}
}

func getComponentType(v interface{}) string {
	switch v.(type) {
	case components.Position:
		return "components.Position"
	case components.Velocity:
		return "components.Velocity"
	case components.Rotation:
		return "components.Rotation"
	case components.Input:
		return "components.Input"
	case components.Player:
		return "components.Player"
	case components.Renderable:
		return "components.Renderable"
	case components.Lifetime:
		return "components.Lifetime"
	case components.Collider:
		return "components.Collider"
	case components.Asteroid:
		return "components.Asteroid"
	case components.Explosion:
		return "components.Explosion"
	case components.Invulnerable:
		return "components.Invulnerable"
	default:
		return fmt.Sprintf("%T", v)
	}
}

func getType(v interface{}) string {
	switch v.(type) {
	case components.Position:
		return "Position"
	case components.Velocity:
		return "Velocity"
	case components.Rotation:
		return "Rotation"
	case components.Renderable:
		return "Renderable"
	case components.Player:
		return "Player"
	case components.Input:
		return "Input"
	case components.Collider:
		return "Collider"
	case components.Lifetime:
		return "Lifetime"
	case components.Asteroid:
		return "Asteroid"
	case components.Explosion:
		return "Explosion"
	default:
		return fmt.Sprintf("%T", v)
	}
}
