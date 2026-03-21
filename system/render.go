package system

import (
	"github.com/abjeloevic88/go-game/ecs"
	"github.com/hajimehoshi/ebiten/v2"
)

// RenderSystem draws entities that have both a Position and a Sprite component.
type RenderSystem struct{}

// Draw iterates all entities and draws each sprite centered on its position.
func (s *RenderSystem) Draw(world *ecs.World, screen *ebiten.Image) {
	for _, e := range world.Entities() {
		pos, hasPos := world.GetPosition(e)
		sprite, hasSprite := world.GetSprite(e)
		if !hasPos || !hasSprite || sprite.Image == nil {
			continue
		}

		w := float64(sprite.Image.Bounds().Dx())
		h := float64(sprite.Image.Bounds().Dy())

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(pos.X-w/2, pos.Y-h/2)
		screen.DrawImage(sprite.Image, opts)
	}
}
