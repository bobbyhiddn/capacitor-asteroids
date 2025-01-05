package highscore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"syscall/js"
)

type Score struct {
	Value     int       `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

type HighScores struct {
	Scores []Score `json:"scores"`
	mu     sync.Mutex
}

const (
	maxScores        = 10
	highScoreKey     = "asteroids_high_scores"
	highScoreFile    = "highscores.json"
	defaultDirectory = "asteroids_data"
)

var (
	instance *HighScores
	once     sync.Once
)

// GetInstance returns the singleton instance of HighScores
func GetInstance() *HighScores {
	once.Do(func() {
		instance = &HighScores{
			Scores: make([]Score, 0),
		}
		instance.load()
	})
	return instance
}

// AddScore adds a new score and maintains only the top scores
func (hs *HighScores) AddScore(value int) {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	score := Score{
		Value:     value,
		Timestamp: time.Now(),
	}

	hs.Scores = append(hs.Scores, score)
	sort.Slice(hs.Scores, func(i, j int) bool {
		return hs.Scores[i].Value > hs.Scores[j].Value
	})

	if len(hs.Scores) > maxScores {
		hs.Scores = hs.Scores[:maxScores]
	}

	hs.save()
}

// GetTopScores returns the top scores
func (hs *HighScores) GetTopScores() []Score {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	
	scores := make([]Score, len(hs.Scores))
	copy(scores, hs.Scores)
	return scores
}

// IsHighScore checks if a score would be a new high score
func (hs *HighScores) IsHighScore(value int) bool {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	if len(hs.Scores) < maxScores {
		return true
	}
	return value > hs.Scores[len(hs.Scores)-1].Value
}

func (hs *HighScores) load() {
	if isWasm() {
		hs.loadWasm()
	} else {
		hs.loadFile()
	}
}

func (hs *HighScores) save() {
	if isWasm() {
		hs.saveWasm()
	} else {
		hs.saveFile()
	}
}

func (hs *HighScores) loadWasm() {
	window := js.Global().Get("window")
	localStorage := window.Get("localStorage")
	
	if localStorage.IsUndefined() {
		return
	}

	data := localStorage.Call("getItem", highScoreKey)
	if data.IsNull() {
		return
	}

	if err := json.Unmarshal([]byte(data.String()), &hs.Scores); err != nil {
		// If there's an error, start with empty scores
		hs.Scores = make([]Score, 0)
	}
}

func (hs *HighScores) saveWasm() {
	window := js.Global().Get("window")
	localStorage := window.Get("localStorage")
	
	if localStorage.IsUndefined() {
		return
	}

	data, err := json.Marshal(hs.Scores)
	if err != nil {
		return
	}

	localStorage.Call("setItem", highScoreKey, string(data))
}

func (hs *HighScores) loadFile() {
	dir := getDataDir()
	path := filepath.Join(dir, highScoreFile)

	data, err := os.ReadFile(path)
	if err != nil {
		// If file doesn't exist, start with empty scores
		hs.Scores = make([]Score, 0)
		return
	}

	if err := json.Unmarshal(data, &hs.Scores); err != nil {
		// If there's an error, start with empty scores
		hs.Scores = make([]Score, 0)
	}
}

func (hs *HighScores) saveFile() {
	dir := getDataDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return
	}

	path := filepath.Join(dir, highScoreFile)
	data, err := json.Marshal(hs.Scores)
	if err != nil {
		return
	}

	_ = os.WriteFile(path, data, 0644)
}

func getDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return filepath.Join(homeDir, defaultDirectory)
}

func isWasm() bool {
	return isWasmInit
}

// Initialize WASM detection
var isWasmInit = false

func init() {
	// This will be true when running in WebAssembly
	isWasmInit = isRunningInWasm()
}

func isRunningInWasm() bool {
	// Check if we're running in WebAssembly
	return js.Global().Get("window").Truthy()
}
