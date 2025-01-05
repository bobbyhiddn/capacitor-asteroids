package systems

import (
	"github.com/bobbyhiddn/ecs-asteroids/components"
	"github.com/bobbyhiddn/ecs-asteroids/ecs"
	"github.com/bobbyhiddn/ecs-asteroids/highscore"
)

type ScoreSystem struct {
	world      *ecs.World
	highScores *highscore.HighScores
}

func NewScoreSystem(world *ecs.World) *ScoreSystem {
	return &ScoreSystem{
		world:      world,
		highScores: highscore.GetInstance(),
	}
}

func (s *ScoreSystem) Update(dt float64) {
	players := s.world.Components["components.Player"]
	for _, playerInterface := range players {
		player := playerInterface.(components.Player)
		if player.IsGameOver {
			// When game is over, check and save high score
			if s.highScores.IsHighScore(player.Score) {
				s.highScores.AddScore(player.Score)
			}
		}
	}
}

func (s *ScoreSystem) GetTopScores() []highscore.Score {
	return s.highScores.GetTopScores()
}
