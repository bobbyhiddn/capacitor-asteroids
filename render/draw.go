package render

import (
	"image/color"
	"math"

	"github.com/bobbyhiddn/ecs-asteroids/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 800
	screenHeight = 600

	// Ship scale factors
	shipLength = screenHeight / 36 // Length from nose to rear
	shipWidth  = screenHeight / 40 // Width at the rear
	rearOffset = shipLength * 0.5  // Distance from center to rear
)

func DrawShip(screen *ebiten.Image, x, y, angle float64, isThrusting bool) {
	// Ship vertices (triangle)
	shipPoints := []struct{ x, y float64 }{
		{-10, 10},  // Bottom left
		{20, 0},    // Tip
		{-10, -10}, // Bottom right
	}

	// Draw ship lines
	for i := 0; i < len(shipPoints); i++ {
		p1 := shipPoints[i]
		p2 := shipPoints[(i+1)%len(shipPoints)]

		transformed1 := transformPoint(p1.x, p1.y, x, y, angle)
		transformed2 := transformPoint(p2.x, p2.y, x, y, angle)

		drawLine(screen, transformed1, transformed2, color.White)
	}

	// Draw thrusters if the ship is thrusting
	if isThrusting {
		// Calculate the base points for the flames (slightly behind the ship)
		backOffset := 12.0
		leftBase := struct{ x, y float64 }{-backOffset, 5}
		rightBase := struct{ x, y float64 }{-backOffset, -5}

		// Calculate flame tips (random length for animation effect)
		flameLength := 15.0
		leftTip := struct{ x, y float64 }{-(backOffset + flameLength), 0}
		rightTip := struct{ x, y float64 }{-(backOffset + flameLength), 0}

		// Draw flames
		drawLine(screen,
			transformPoint(leftBase.x, leftBase.y, x, y, angle),
			transformPoint(leftTip.x, leftTip.y, x, y, angle),
			color.RGBA{R: 255, G: 100, B: 0, A: 255},
		)
		drawLine(screen,
			transformPoint(rightTip.x, rightTip.y, x, y, angle),
			transformPoint(rightBase.x, rightBase.y, x, y, angle),
			color.RGBA{R: 255, G: 100, B: 0, A: 255},
		)
	}
}

func DrawBullet(screen *ebiten.Image, x, y float64) {
	size := 2.0
	// Draw a small filled square
	ebitenutil.DrawLine(screen, x-size, y-size, x+size, y-size, color.White) // Top
	ebitenutil.DrawLine(screen, x+size, y-size, x+size, y+size, color.White) // Right
	ebitenutil.DrawLine(screen, x+size, y+size, x-size, y+size, color.White) // Bottom
	ebitenutil.DrawLine(screen, x-size, y+size, x-size, y-size, color.White) // Left
	// Add diagonal lines to make it look more solid
	ebitenutil.DrawLine(screen, x-size, y-size, x+size, y+size, color.White)
	ebitenutil.DrawLine(screen, x-size, y+size, x+size, y-size, color.White)
}

func DrawAsteroid(screen *ebiten.Image, x, y, angle, scale float64) {
	// Define asteroid shape as a rough circle with some variation
	numPoints := 12
	baseRadius := 20.0 * scale
	points := make([]struct{ x, y float64 }, numPoints)

	// Generate points
	for i := 0; i < numPoints; i++ {
		pointAngle := float64(i) * 2 * math.Pi / float64(numPoints)
		// Add some randomness to the radius for each vertex
		radius := baseRadius * (0.8 + 0.4*math.Sin(float64(i)*3))
		points[i] = struct{ x, y float64 }{
			x: radius * math.Cos(pointAngle),
			y: radius * math.Sin(pointAngle),
		}
	}

	// Draw lines between points
	for i := 0; i < len(points); i++ {
		p1 := points[i]
		p2 := points[(i+1)%len(points)]

		transformed1 := transformPoint(p1.x, p1.y, x, y, angle)
		transformed2 := transformPoint(p2.x, p2.y, x, y, angle)

		drawLine(screen, transformed1, transformed2, color.White)
	}
}

func DrawExplosion(screen *ebiten.Image, x, y float64, explosion components.Explosion) {
	// Calculate current radius based on age
	progress := explosion.Age / explosion.MaxAge
	currentRadius := explosion.Radius * (1 + progress) // Expand over time

	// Calculate alpha (fade out)
	alpha := 1.0 - progress

	// Draw explosion particles
	for i := 0; i < explosion.Pieces; i++ {
		angle := float64(i) * (2 * math.Pi / float64(explosion.Pieces))

		// Particles move outward
		particleX := x + math.Cos(angle)*currentRadius
		particleY := y + math.Sin(angle)*currentRadius

		// Draw particle line from center
		ebitenutil.DrawLine(
			screen,
			x, y,
			particleX, particleY,
			color.RGBA{255, 200, 50, uint8(255 * alpha)},
		)
	}
}

type point struct {
	x, y float64
}

func transformPoint(x, y, centerX, centerY, angle float64) point {
	// Rotate
	rotatedX := x*math.Cos(angle) - y*math.Sin(angle)
	rotatedY := x*math.Sin(angle) + y*math.Cos(angle)

	// Translate
	return point{
		x: rotatedX + centerX,
		y: rotatedY + centerY,
	}
}

func drawLine(screen *ebiten.Image, p1, p2 point, clr color.Color) {
	ebitenutil.DrawLine(screen, p1.x, p1.y, p2.x, p2.y, clr)
}
