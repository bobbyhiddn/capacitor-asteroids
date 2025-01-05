package game

// Screen dimensions
const (
	DefaultWidth  = 1200
	DefaultHeight = 600
)

// Screen provides access to screen dimensions and utilities
type Screen struct {
	width  int
	height int
}

// NewScreen creates a new Screen instance
func NewScreen() *Screen {
	return &Screen{
		width:  DefaultWidth,
		height: DefaultHeight,
	}
}

// Width returns the current screen width
func (s *Screen) Width() int {
	return s.width
}

// Height returns the current screen height
func (s *Screen) Height() int {
	return s.height
}

// SetDimensions updates the screen dimensions
func (s *Screen) SetDimensions(w, h int) {
	s.width = w
	s.height = h
}

// GetActualWidth implements types.GameDimensions
func (s *Screen) GetActualWidth() int {
	return s.width
}

// GetActualHeight implements types.GameDimensions
func (s *Screen) GetActualHeight() int {
	return s.height
}

// CenterX returns the horizontal center of the screen
func (s *Screen) CenterX() float64 {
	return float64(s.width) / 2
}

// CenterY returns the vertical center of the screen
func (s *Screen) CenterY() float64 {
	return float64(s.height) / 2
}

// InBounds checks if the given coordinates are within the screen bounds
func (s *Screen) InBounds(x, y float64) bool {
	return x >= 0 && x < float64(s.width) && y >= 0 && y < float64(s.height)
}

// WrapCoordinates wraps coordinates around screen edges
func (s *Screen) WrapCoordinates(x, y float64) (float64, float64) {
	w, h := float64(s.width), float64(s.height)
	wrappedX := x
	wrappedY := y

	if x < 0 {
		wrappedX = x + w
	} else if x >= w {
		wrappedX = x - w
	}

	if y < 0 {
		wrappedY = y + h
	} else if y >= h {
		wrappedY = y - h
	}

	return wrappedX, wrappedY
}
