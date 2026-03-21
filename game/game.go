// Package game provides the top-level Ebiten Game implementation that owns the ECS World.
package game

import (
	"image/color"

	"github.com/abjeloevic88/go-game/component"
	"github.com/abjeloevic88/go-game/config"
	"github.com/abjeloevic88/go-game/ecs"
	"github.com/abjeloevic88/go-game/system"
	"github.com/hajimehoshi/ebiten/v2"
)

// Game implements ebiten.Game and owns the ECS World.
type Game struct {
	World *ecs.World
}

// NewGame creates a Game, initializes the player entity and registers all systems.
func NewGame() *Game {
	w := ecs.NewWorld()

	// Create the player entity.
	player := w.CreateEntity()
	w.AddPosition(player, &component.Position{X: config.ScreenWidth / 2, Y: config.ScreenHeight / 2})
	w.AddVelocity(player, &component.Velocity{X: 0, Y: 0})

	// Create a 32×32 bright-green placeholder sprite.
	img := ebiten.NewImage(config.SpriteSize, config.SpriteSize)
	img.Fill(color.RGBA{0, 255, 0, 255})
	w.AddSprite(player, &component.Sprite{Image: img})

	// Register update systems in order: Input → Movement.
	w.AddSystem(&system.InputSystem{})
	w.AddSystem(&system.MovementSystem{})

	// Register draw systems.
	w.AddDrawSystem(&system.RenderSystem{})

	return &Game{World: w}
}

// Update delegates to the World's update systems.
func (g *Game) Update() error {
	g.World.Update()
	return nil
}

// Draw delegates to the World's draw systems.
func (g *Game) Draw(screen *ebiten.Image) {
	g.World.Draw(screen)
}

// Layout returns the logical screen resolution.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
