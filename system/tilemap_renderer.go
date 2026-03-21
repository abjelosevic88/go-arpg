package system

import (
	"github.com/abjeloevic88/go-game/ecs"
	"github.com/abjeloevic88/go-game/tilemap"
	"github.com/hajimehoshi/ebiten/v2"
)

// TilemapRenderer is a DrawSystem that draws a Tilemap as the background.
type TilemapRenderer struct {
	Map *tilemap.Tilemap
}

// Draw iterates every cell in the tilemap and draws the corresponding
// tile image from the TileSet at the correct pixel position.
func (r *TilemapRenderer) Draw(world *ecs.World, screen *ebiten.Image) {
	tm := r.Map
	for row := 0; row < tm.Rows; row++ {
		for col := 0; col < tm.Cols; col++ {
			idx := tm.TileAt(col, row)
			img, ok := tm.TileSet.Tiles[idx]
			if !ok || img == nil {
				continue
			}
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(col*tm.TileW), float64(row*tm.TileH))
			screen.DrawImage(img, opts)
		}
	}
}
