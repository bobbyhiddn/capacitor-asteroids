package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/bobbyhiddn/ecs-asteroids/components"
	"github.com/bobbyhiddn/ecs-asteroids/ecs"
	"github.com/bobbyhiddn/ecs-asteroids/game"
	"github.com/bobbyhiddn/ecs-asteroids/systems"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 800
	screenHeight = 600

	// Game settings
	initialAsteroids      = 4
	asteroidSpawnInterval = 5 * time.Second
)

type Game struct {
	world              *ecs.World
	inputSystem        *systems.InputSystem
	playerSystem       *systems.PlayerSystem
	movementSystem     *systems.MovementSystem
	collisionSystem    *systems.CollisionSystem
	renderSystem       *systems.RenderSystem
	asteroidSpawner    *systems.AsteroidSpawnerSystem
	explosionSystem    *systems.ExplosionSystem
	invulnerableSystem *systems.InvulnerableSystem
}

func NewGame() *Game {
	g := &Game{
		world: ecs.NewWorld(),
	}

	// Create systems
	g.inputSystem = systems.NewInputSystem(g.world)
	g.playerSystem = systems.NewPlayerSystem(g.world)
	g.movementSystem = systems.NewMovementSystem(g.world)
	g.collisionSystem = systems.NewCollisionSystem(g.world)
	g.renderSystem = systems.NewRenderSystem(g.world, ebiten.NewImage(screenWidth, screenHeight))
	g.asteroidSpawner = systems.NewAsteroidSpawnerSystem(g.world)
	g.explosionSystem = systems.NewExplosionSystem(g.world)
	g.invulnerableSystem = systems.NewInvulnerableSystem(g.world)

	g.world.AddSystem(g.inputSystem)
	g.world.AddSystem(g.playerSystem)
	g.world.AddSystem(g.movementSystem)
	g.world.AddSystem(g.invulnerableSystem)
	g.world.AddSystem(g.collisionSystem)
	g.world.AddSystem(g.renderSystem)
	g.world.AddSystem(g.asteroidSpawner)
	g.world.AddSystem(g.explosionSystem)

	// Create player ship
	game.CreatePlayerShip(g.world, float64(screenWidth/2), float64(screenHeight/2))

	// Create initial asteroids
	for i := 0; i < initialAsteroids; i++ {
		game.CreateAsteroid(g.world, rand.Intn(3))
	}

	return g
}

func (g *Game) Update() error {
	dt := 1.0 / 60.0

	// Always update input system to check for restart
	g.inputSystem.Update(dt)

	// Get player state
	players := g.world.Components["components.Player"]
	for _, playerInterface := range players {
		player := playerInterface.(components.Player)
		if player.IsGameOver {
			fmt.Printf("Game is over, waiting for restart input...\n")
			// Only update input system during game over
			return nil
		}
	}

	// Update all systems when game is active
	g.playerSystem.Update(dt)
	g.movementSystem.Update(dt)
	g.invulnerableSystem.Update(dt)
	g.collisionSystem.Update(dt)
	g.asteroidSpawner.Update(dt)
	g.explosionSystem.Update(dt)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen
	screen.Fill(color.Black)

	// Draw the game onto the screen
	g.renderSystem.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("ECS Asteroids")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
