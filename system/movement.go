package system

import (
	"github.com/abjeloevic88/go-game/config"
	"github.com/abjeloevic88/go-game/ecs"
)

// MovementSystem applies velocity to position and clamps to screen bounds.
type MovementSystem struct{}

// Update adds velocity to position for all entities with both components,
// then clamps position to screen bounds using sprite dimensions.
func (s *MovementSystem) Update(world *ecs.World) {
	for _, e := range world.Entities() {
		pos, hasPos := world.GetPosition(e)
		vel, hasVel := world.GetVelocity(e)
		if !hasPos || !hasVel {
			continue
		}

		// Apply velocity to position.
		pos.X += vel.X
		pos.Y += vel.Y

		// Clamp to screen bounds using sprite dimensions if available.
		sprite, hasSprite := world.GetSprite(e)
		if hasSprite && sprite.Image != nil {
			sw := float64(sprite.Image.Bounds().Dx())
			sh := float64(sprite.Image.Bounds().Dy())

			maxX := float64(config.ScreenWidth) - sw
			maxY := float64(config.ScreenHeight) - sh

			if pos.X < 0 {
				pos.X = 0
			} else if pos.X > maxX {
				pos.X = maxX
			}

			if pos.Y < 0 {
				pos.Y = 0
			} else if pos.Y > maxY {
				pos.Y = maxY
			}
		}
	}
}
