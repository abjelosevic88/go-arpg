// Package ecs provides the Entity-Component-System core: entities, component storage, and system orchestration.
package ecs

import (
	"github.com/abjeloevic88/go-game/component"
	"github.com/hajimehoshi/ebiten/v2"
)

// System is implemented by any system that runs during Update.
type System interface {
	Update(world *World)
}

// DrawSystem is implemented by any system that runs during Draw.
type DrawSystem interface {
	Draw(world *World, screen *ebiten.Image)
}

// World is the central ECS container that owns entities, components, and systems.
type World struct {
	nextEntity  Entity
	positions   map[Entity]*component.Position
	velocities  map[Entity]*component.Velocity
	sprites     map[Entity]*component.Sprite
	systems     []System
	drawSystems []DrawSystem
}

// NewWorld creates and returns an initialized World.
func NewWorld() *World {
	return &World{
		nextEntity: 0,
		positions:  make(map[Entity]*component.Position),
		velocities: make(map[Entity]*component.Velocity),
		sprites:    make(map[Entity]*component.Sprite),
	}
}

// CreateEntity allocates a new unique Entity ID and returns it.
func (w *World) CreateEntity() Entity {
	e := w.nextEntity
	w.nextEntity++
	return e
}

// AddPosition attaches a Position component to the given entity.
func (w *World) AddPosition(e Entity, c *component.Position) {
	w.positions[e] = c
}

// GetPosition retrieves the Position component for the given entity.
func (w *World) GetPosition(e Entity) (*component.Position, bool) {
	c, ok := w.positions[e]
	return c, ok
}

// AddVelocity attaches a Velocity component to the given entity.
func (w *World) AddVelocity(e Entity, c *component.Velocity) {
	w.velocities[e] = c
}

// GetVelocity retrieves the Velocity component for the given entity.
func (w *World) GetVelocity(e Entity) (*component.Velocity, bool) {
	c, ok := w.velocities[e]
	return c, ok
}

// AddSprite attaches a Sprite component to the given entity.
func (w *World) AddSprite(e Entity, c *component.Sprite) {
	w.sprites[e] = c
}

// GetSprite retrieves the Sprite component for the given entity.
func (w *World) GetSprite(e Entity) (*component.Sprite, bool) {
	c, ok := w.sprites[e]
	return c, ok
}

// AddSystem appends a System to the ordered update list.
func (w *World) AddSystem(s System) {
	w.systems = append(w.systems, s)
}

// AddDrawSystem appends a DrawSystem to the ordered draw list.
func (w *World) AddDrawSystem(s DrawSystem) {
	w.drawSystems = append(w.drawSystems, s)
}

// Update invokes each registered System's Update method in registration order.
func (w *World) Update() {
	for _, s := range w.systems {
		s.Update(w)
	}
}

// Draw invokes each registered DrawSystem's Draw method in registration order.
func (w *World) Draw(screen *ebiten.Image) {
	for _, ds := range w.drawSystems {
		ds.Draw(w, screen)
	}
}

// Entities returns a slice of all created entity IDs.
func (w *World) Entities() []Entity {
	entities := make([]Entity, 0, int(w.nextEntity))
	for i := Entity(0); i < w.nextEntity; i++ {
		entities = append(entities, i)
	}
	return entities
}
