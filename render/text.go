package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

var (
	// Default font face for regular text
	DefaultFace = basicfont.Face7x13

	// Large font face for high scores and titles (using basicfont)
	LargeFace = basicfont.Face7x13
)

// DrawCenteredText draws text centered horizontally at the specified y position
func DrawCenteredText(screen *ebiten.Image, str string, y int, clr color.Color, face font.Face) {
	bound := text.BoundString(face, str)
	x := (screen.Bounds().Dx() - bound.Dx()) / 2
	text.Draw(screen, str, face, x, y, clr)
}

// DrawCenteredScaledText draws text centered horizontally at the specified y position with scaling
func DrawCenteredScaledText(screen *ebiten.Image, str string, y int, scale float64, clr color.Color, face font.Face) {
	// Create a temporary image to draw the text
	bound := text.BoundString(face, str)
	w := bound.Dx()
	h := bound.Dy()
	tmpImg := ebiten.NewImage(w, h)

	// Draw the text onto the temporary image
	text.Draw(tmpImg, str, face, 0, -bound.Min.Y, clr)

	// Calculate the scaled dimensions
	scaledW := int(float64(w) * scale)

	// Calculate position to center the scaled text
	x := (screen.Bounds().Dx() - scaledW) / 2

	// Create options for drawing the scaled image
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(scale, scale)
	opts.GeoM.Translate(float64(x), float64(y))

	screen.DrawImage(tmpImg, opts)
}

// DrawText draws text at the specified position
func DrawText(screen *ebiten.Image, str string, x, y int, clr color.Color, face font.Face) {
	text.Draw(screen, str, face, x, y, clr)
}
