# ARPG

A Diablo 1-style action RPG built with Go and [Ebiten](https://ebitengine.org/) v2.9.9.

Currently features a player character (green square) that moves via WASD input within a 1024×768 window. The architecture follows an Entity-Component-System (ECS) inspired pattern designed for easy extension with combat, inventory, enemies, maps, and more.

## Getting Started

```bash
go run main.go
```

### Controls

| Key | Action    |
|-----|-----------|
| W   | Move up   |
| A   | Move left |
| S   | Move down |
| D   | Move right|

Diagonal movement is normalized so speed stays consistent in all directions. The character is clamped to screen bounds.

## Project Structure

```
├── main.go              # Entry point — window config and game launch
├── config/
│   └── config.go        # Shared constants (screen size, speed, sprite size)
├── ecs/
│   ├── entity.go        # Entity type (int ID)
│   └── world.go         # World: entity creation, component storage, system orchestration
├── component/
│   ├── position.go      # Position component (X, Y float64)
│   ├── velocity.go      # Velocity component (X, Y float64)
│   └── sprite.go        # Sprite component (*ebiten.Image)
├── system/
│   ├── input.go         # InputSystem: WASD → velocity
│   ├── movement.go      # MovementSystem: velocity → position + boundary clamping
│   └── render.go        # RenderSystem: draw sprites centered on position
└── game/
    └── game.go          # Game struct implementing ebiten.Game, owns the World
```

## Architecture

The game loop delegates to a `World` that owns entities, components, and systems. Each tick, systems run in a fixed order:

1. **InputSystem** — reads keyboard state, sets velocity
2. **MovementSystem** — applies velocity to position, clamps to screen bounds
3. **RenderSystem** — draws sprites centered on entity positions

Components are plain data structs. Systems contain all logic. New features can be added by defining new components and systems without modifying existing ones.

## Requirements

- Go 1.26+
- Ebiten v2.9.9
