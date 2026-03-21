// Package game provides the top-level Ebiten Game implementation that owns the ECS World.
package game

import (
	"fmt"

	"github.com/abjeloevic88/go-game/asset"
	"github.com/abjeloevic88/go-game/component"
	"github.com/abjeloevic88/go-game/config"
	"github.com/abjeloevic88/go-game/ecs"
	"github.com/abjeloevic88/go-game/system"
	"github.com/abjeloevic88/go-game/tmx"
	"github.com/hajimehoshi/ebiten/v2"
)

// Game implements ebiten.Game and owns the ECS World.
type Game struct {
	World *ecs.World
}

// NewGame creates a Game, initializes the player entity and registers all systems.
// Returns an error if any asset fails to load.
func NewGame() (*Game, error) {
	w := ecs.NewWorld()

	// Load the player sprite sheet and extract the first 32×32 frame.
	playerSheet, err := asset.LoadImage("assets/player.png")
	if err != nil {
		return nil, err
	}
	playerFrame, err := asset.SubImage(playerSheet, 0, 0, config.SpriteSize, config.SpriteSize)
	if err != nil {
		return nil, err
	}

	// Load tilemap from TMX (automatically resolves TSX tileset and spritesheet).
	tm, err := tmx.LoadMapFromTMX("maps/map.tmx")
	if err != nil {
		return nil, fmt.Errorf("load map: %w", err)
	}

	// Create the player entity.
	player := w.CreateEntity()
	w.AddPosition(player, &component.Position{X: config.ScreenWidth / 2, Y: config.ScreenHeight / 2})
	w.AddVelocity(player, &component.Velocity{X: 0, Y: 0})
	w.AddSprite(player, &component.Sprite{Image: playerFrame})

	// Register update systems in order: Input → Movement.
	w.AddSystem(&system.InputSystem{})
	w.AddSystem(&system.MovementSystem{})

	// Register draw systems: TilemapRenderer first (background), then RenderSystem (entities).
	w.AddDrawSystem(&system.TilemapRenderer{Map: tm})
	w.AddDrawSystem(&system.RenderSystem{})

	return &Game{World: w}, nil
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
