package types

// GameDimensions provides access to the actual game window dimensions
type GameDimensions interface {
	GetActualWidth() int
	GetActualHeight() int
}
