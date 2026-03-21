// Package tilemap provides a tile-based map system for the game world.
package tilemap

import "github.com/hajimehoshi/ebiten/v2"

// TileSet maps tile indices to their Ebiten images.
type TileSet struct {
	Tiles map[int]*ebiten.Image
}

// Tilemap represents a 2D grid of tile indices.
type Tilemap struct {
	Cols    int       // number of columns
	Rows    int       // number of rows
	TileW   int       // tile width in pixels
	TileH   int       // tile height in pixels
	Grid    [][]int   // Grid[row][col] = tile index
	TileSet *TileSet
}

// TileAt returns the tile index at (col, row).
// Returns 0 (grass) for out-of-bounds coordinates.
func (tm *Tilemap) TileAt(col, row int) int {
	if row < 0 || row >= tm.Rows || col < 0 || col >= tm.Cols {
		return 0
	}
	return tm.Grid[row][col]
}

// NewDefaultMap creates a Tilemap sized to cover the screen with a crossroads
// pattern: path tiles (index 1) at the center row and center column, grass
// tiles (index 0) everywhere else.
func NewDefaultMap(screenWidth, screenHeight, tileSize int, ts *TileSet) *Tilemap {
	cols := screenWidth / tileSize
	rows := screenHeight / tileSize

	grid := make([][]int, rows)
	centerRow := rows / 2
	centerCol := cols / 2

	for r := 0; r < rows; r++ {
		grid[r] = make([]int, cols)
		for c := 0; c < cols; c++ {
			if r == centerRow || c == centerCol {
				grid[r][c] = 1 // path
			}
			// default zero value is 0 (grass)
		}
	}

	return &Tilemap{
		Cols:    cols,
		Rows:    rows,
		TileW:   tileSize,
		TileH:   tileSize,
		Grid:    grid,
		TileSet: ts,
	}
}
