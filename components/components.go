package components

import (
	"time"
	"image/color"
)

type Position struct {
	X, Y float64
}

type Velocity struct {
	DX, DY    float64
	MaxSpeed  float64
}

type Rotation struct {
	Angle         float64
	RotationSpeed float64
}

type RenderableType int

const (
	RenderableTypeShip RenderableType = iota
	RenderableTypeBullet
	RenderableTypeAsteroid
	RenderableTypeExplosion
)

type Renderable struct {
	Type    RenderableType
	Scale   float64
	Visible bool
}

type Lifetime struct {
	Created  time.Time
	Duration time.Duration
}

type Player struct {
	IsThrusting bool
	Score       int
	Lives       int
	IsGameOver  bool
}

type ColliderType int

const (
	ColliderTypeShip ColliderType = iota
	ColliderTypeBullet
	ColliderTypeAsteroid
)

type Collider struct {
	Radius float64
	Type   ColliderType
}

type Input struct {
	Rotate  float64 // -1 for left, 1 for right
	Forward bool
	Shoot   bool
	MouseX  int
	MouseY  int
	MousePressed bool
}

type Asteroid struct {
	Size int // 0 = small, 1 = medium, 2 = large
}

type Explosion struct {
	Age     float64 // Time since explosion started
	MaxAge  float64 // When to remove the explosion
	Radius  float64 // Current radius of explosion
	Pieces  int     // Number of particles
}

type Invulnerable struct {
	Duration float64 // How long the invulnerability lasts
	Timer    float64 // Current time left
}

type Bullet struct {
	ShooterID int
}

// UI button component
type UIButton struct {
	X, Y, Width, Height int
	Color              color.Color
}
