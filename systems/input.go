package systems

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/bobbyhiddn/ecs-asteroids/components"
	"github.com/bobbyhiddn/ecs-asteroids/ecs"
	"github.com/bobbyhiddn/ecs-asteroids/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InputSystem struct {
	world  *ecs.World
	screen *game.Screen
}

func NewInputSystem(world *ecs.World) *InputSystem {
	return &InputSystem{
		world:  world,
		screen: game.NewScreen(),
	}
}

func (s *InputSystem) Update(dt float64) {
	players := s.world.Components["components.Player"]
	inputs := s.world.Components["components.Input"]

	for id, playerInterface := range players {
		player := playerInterface.(components.Player)
		input := inputs[id].(components.Input)

		// Check for game over restart
		if player.IsGameOver {
			// Check for any key press or new touch/click
			if len(inpututil.AppendPressedKeys(nil)) > 0 ||
				len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 ||
				inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				fmt.Printf("Input detected during game over, restarting...\n")
				s.handleGameRestart(id, &player)
				continue
			}
			continue // Skip other input processing when game over
		}

		// Reset input state
		input.Rotate = 0
		input.Forward = false
		input.Shoot = false
		input.MousePressed = false

		// Process multitouch inputs
		touchIDs := ebiten.TouchIDs()
		justPressedTouchIDs := inpututil.AppendJustPressedTouchIDs(nil)

		for _, touchID := range touchIDs {
			x, y := ebiten.TouchPosition(touchID)
			input.MouseX = x
			input.MouseY = y
			input.MousePressed = true

			// Check if touch is within the fire button area
			fireButtonY := s.screen.Height() - 100
			if x >= 20 && x <= 180 && y >= fireButtonY-80 && y <= fireButtonY+80 {
				// Only shoot if this is a new touch
				for _, justPressedID := range justPressedTouchIDs {
					if touchID == justPressedID {
						input.Shoot = true
						break
					}
				}
			} else {
				// Process directional input
				s.processDirectionalInput(id, float64(x), float64(y), &input)
			}
		}

		// Handle mouse input for desktop testing
		if len(touchIDs) == 0 {
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				x, y := ebiten.CursorPosition()
				input.MouseX = x
				input.MouseY = y
				input.MousePressed = true

				fireButtonY := s.screen.Height() - 100
				if x >= 20 && x <= 180 && y >= fireButtonY-80 && y <= fireButtonY+80 {
					// Only shoot if mouse button was just pressed
					input.Shoot = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
				} else {
					s.processDirectionalInput(id, float64(x), float64(y), &input)
				}
			}
		}

		// Handle keyboard input
		if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			input.Rotate = -1
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			input.Rotate = 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
			input.Forward = true
		}
		input.Shoot = input.Shoot || inpututil.IsKeyJustPressed(ebiten.KeySpace)

		// Update input component
		s.world.AddComponent(id, input)
	}
}

func (s *InputSystem) processDirectionalInput(id ecs.EntityID, x, y float64, input *components.Input) {
	if pos, ok := s.world.Components["components.Position"][id].(components.Position); ok {
		dx := x - pos.X
		dy := y - pos.Y
		targetAngle := math.Atan2(dy, dx)

		if rot, ok := s.world.Components["components.Rotation"][id].(components.Rotation); ok {
			diff := math.Mod(targetAngle-rot.Angle+math.Pi, 2*math.Pi) - math.Pi
			if math.Abs(diff) > 0.1 {
				if diff > 0 {
					input.Rotate = 1
				} else {
					input.Rotate = -1
				}
			}
		}
		input.Forward = true
	}
}

func (s *InputSystem) handleGameRestart(id ecs.EntityID, player *components.Player) {
	// Reset player
	player.Lives = 3
	player.Score = 0
	player.IsGameOver = false
	s.world.AddComponent(id, *player)

	// Add position component
	s.world.AddComponent(id, components.Position{
		X: float64(s.screen.Width() / 2),
		Y: float64(s.screen.Height() / 2),
	})

	// Add velocity component
	s.world.AddComponent(id, components.Velocity{
		DX:       0,
		DY:       0,
		MaxSpeed: 400,
	})

	// Add rotation component
	s.world.AddComponent(id, components.Rotation{
		Angle: 0,
	})

	// Add renderable component
	s.world.AddComponent(id, components.Renderable{
		Visible: true,
	})

	// Add collider component
	s.world.AddComponent(id, components.Collider{
		Type:   components.ColliderTypeShip,
		Radius: 20,
	})

	// Add input component
	s.world.AddComponent(id, components.Input{})

	// Clear all asteroids
	for asteroidID := range s.world.Components["components.Asteroid"] {
		s.world.DestroyEntity(asteroidID)
	}

	// Create initial asteroids
	for i := 0; i < 4; i++ {
		game.CreateAsteroid(s.world, rand.Intn(3))
	}
}
