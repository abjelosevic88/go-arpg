// Package system implements ECS systems for input, movement, and rendering.
package system

import (
	"math"

	"github.com/abjeloevic88/go-game/config"
	"github.com/abjeloevic88/go-game/ecs"
	"github.com/hajimehoshi/ebiten/v2"
)

// InputSystem reads WASD key state and sets velocity on all entities
// that have a Velocity component.
type InputSystem struct{}

// Update reads keyboard input and sets velocity for all entities with a Velocity component.
func (s *InputSystem) Update(world *ecs.World) {
	for _, e := range world.Entities() {
		vel, ok := world.GetVelocity(e)
		if !ok {
			continue
		}

		// Determine X axis: A (left) and D (right). Opposing keys cancel.
		var dx float64
		left := ebiten.IsKeyPressed(ebiten.KeyA)
		right := ebiten.IsKeyPressed(ebiten.KeyD)
		if left && !right {
			dx = -config.MoveSpeed
		} else if right && !left {
			dx = config.MoveSpeed
		}

		// Determine Y axis: W (up) and S (down). Opposing keys cancel.
		var dy float64
		up := ebiten.IsKeyPressed(ebiten.KeyW)
		down := ebiten.IsKeyPressed(ebiten.KeyS)
		if up && !down {
			dy = -config.MoveSpeed
		} else if down && !up {
			dy = config.MoveSpeed
		}

		// Normalize diagonal movement so magnitude equals MoveSpeed.
		if dx != 0 && dy != 0 {
			mag := math.Sqrt(dx*dx + dy*dy)
			dx = (dx / mag) * config.MoveSpeed
			dy = (dy / mag) * config.MoveSpeed
		}

		vel.X = dx
		vel.Y = dy
	}
}
