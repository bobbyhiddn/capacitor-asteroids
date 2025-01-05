//go:build !(js && wasm)
// +build !js !wasm

package highscore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type Score struct {
	Value     int       `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

type HighScores struct {
	Scores []Score `json:"scores"`
	mu     sync.Mutex
}

var instance *HighScores

func GetInstance() *HighScores {
	if instance == nil {
		instance = &HighScores{}
		instance.load()
	}
	return instance
}

func (hs *HighScores) AddScore(value int) {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	hs.Scores = append(hs.Scores, Score{
		Value:     value,
		Timestamp: time.Now(),
	})

	sort.Slice(hs.Scores, func(i, j int) bool {
		return hs.Scores[i].Value > hs.Scores[j].Value
	})

	if len(hs.Scores) > 10 {
		hs.Scores = hs.Scores[:10]
	}

	hs.save()
}

func (hs *HighScores) GetTopScores() []Score {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	return hs.Scores
}

func (hs *HighScores) IsHighScore(value int) bool {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	if len(hs.Scores) < 10 {
		return true
	}
	return value > hs.Scores[len(hs.Scores)-1].Value
}

func (hs *HighScores) load() {
	hs.loadFile()
}

func (hs *HighScores) save() {
	hs.saveFile()
}

func (hs *HighScores) loadFile() {
	dataDir := getDataDir()
	filePath := filepath.Join(dataDir, "highscores.json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &hs.Scores)
	if err != nil {
		return
	}
}

func (hs *HighScores) saveFile() {
	dataDir := getDataDir()
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return
	}

	filePath := filepath.Join(dataDir, "highscores.json")
	data, err := json.Marshal(hs.Scores)
	if err != nil {
		return
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return
	}
}

func getDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(homeDir, ".ecs-asteroids")
}
