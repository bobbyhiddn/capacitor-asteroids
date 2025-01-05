package main

import (
	"github.com/bobbyhiddn/ecs-asteroids/game"
	"github.com/hajimehoshi/ebiten/v2/mobile"
)

func init() {
	// Initialize the game for mobile
	mobile.SetGame(game.NewGame())
}

func main() {}
