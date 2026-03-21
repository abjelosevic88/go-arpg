// Package component defines ECS data components: Position, Velocity, and Sprite.
package component

// Position stores an Entity's X and Y coordinates in world space.
type Position struct {
	X, Y float64
}
