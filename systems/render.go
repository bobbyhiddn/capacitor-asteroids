package systems

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/samuel-pratt/ebiten-asteroids/components"
	"github.com/samuel-pratt/ebiten-asteroids/ecs"
	"github.com/samuel-pratt/ebiten-asteroids/game"
	"github.com/samuel-pratt/ebiten-asteroids/render"
	"image/color"
	"math"
)

type RenderSystem struct {
	world  *ecs.World
	screen *ebiten.Image
}

func NewRenderSystem(world *ecs.World, screen *ebiten.Image) *RenderSystem {
	return &RenderSystem{
		world:  world,
		screen: screen,
	}
}

func (s *RenderSystem) Update(dt float64) {
	// No update logic needed for rendering
}

func drawDottedCircle(screen *ebiten.Image, x, y, radius float64, c color.Color) {
	numSegments := 16 // Reduced for larger dots
	for i := 0; i < numSegments; i++ {
		angle := float64(i) * 2 * math.Pi / float64(numSegments)
		nextAngle := float64(i+1) * 2 * math.Pi / float64(numSegments)
		
		// Only draw every other segment for dotted effect
		if i%2 == 0 {
			x1 := x + radius*math.Cos(angle)
			y1 := y + radius*math.Sin(angle)
			x2 := x + radius*math.Cos(nextAngle)
			y2 := y + radius*math.Sin(nextAngle)
			ebitenutil.DrawLine(screen, x1, y1, x2, y2, c)
		}
	}
}

func (s *RenderSystem) Draw(screen *ebiten.Image) {
	positions := s.world.Components["components.Position"]
	renderables := s.world.Components["components.Renderable"]
	rotations := s.world.Components["components.Rotation"]
	players := s.world.Components["components.Player"]
	explosions := s.world.Components["components.Explosion"]

	// Draw all renderable entities
	for id, renderableInterface := range renderables {
		renderable := renderableInterface.(components.Renderable)
		if !renderable.Visible {
			continue
		}

		position, ok := positions[id].(components.Position)
		if !ok {
			continue
		}

		rotation := float64(0)
		if r, ok := rotations[id].(components.Rotation); ok {
			rotation = r.Angle
		}

		switch renderable.Type {
		case components.RenderableTypeShip:
			isThrusting := false
			if player, ok := players[id].(components.Player); ok {
				isThrusting = player.IsThrusting
			}
			render.DrawShip(screen, position.X, position.Y, rotation, isThrusting)
		case components.RenderableTypeBullet:
			render.DrawBullet(screen, position.X, position.Y)
		case components.RenderableTypeAsteroid:
			render.DrawAsteroid(screen, position.X, position.Y, rotation, renderable.Scale)
		case components.RenderableTypeExplosion:
			if explosion, ok := explosions[id].(components.Explosion); ok {
				render.DrawExplosion(screen, position.X, position.Y, explosion)
			}
		}
	}

	// Draw UI
	for _, player := range players {
		if p, ok := player.(components.Player); ok {
			// Draw score
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", p.Score), 10, 10)
			
			// Draw lives
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Lives: %d", p.Lives), 10, 30)
			
			// Draw fire button (red dotted circle)
			drawDottedCircle(screen, 60, float64(game.ScreenHeight-60), 40, color.RGBA{255, 0, 0, 255})
			
			// Draw game over
			if p.IsGameOver {
				gameOverText := "GAME OVER - Press Any Key to Restart"
				textWidth := len(gameOverText) * 6 // Approximate width of text
				x := (game.ScreenWidth - textWidth) / 2 // Center horizontally
				ebitenutil.DebugPrintAt(screen, gameOverText, x, 300)
			}
			break // Only show UI for first player
		}
	}
}
