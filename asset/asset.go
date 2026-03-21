package asset

import (
	"fmt"
	"image"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// LoadImage reads a PNG file from the given path and returns an *ebiten.Image.
// Returns a descriptive error if the file cannot be opened or decoded.
func LoadImage(path string) (*ebiten.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("asset: open %s: %w", path, err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("asset: decode %s: %w", path, err)
	}

	return ebiten.NewImageFromImage(img), nil
}

// SubImage extracts a rectangular region from a source image.
// x, y is the top-left corner; w, h is the width and height in pixels.
// Returns an error if the rectangle extends beyond the source bounds.
func SubImage(source *ebiten.Image, x, y, w, h int) (*ebiten.Image, error) {
	bounds := source.Bounds()
	if x < 0 || y < 0 || x+w > bounds.Dx() || y+h > bounds.Dy() {
		return nil, fmt.Errorf("asset: sub-image rect (%d,%d,%d,%d) exceeds source bounds (%d,%d)", x, y, w, h, bounds.Dx(), bounds.Dy())
	}

	sub := source.SubImage(image.Rect(x, y, x+w, y+h)).(*ebiten.Image)
	return sub, nil
}
