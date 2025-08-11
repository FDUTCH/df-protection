package monitoring

import (
	"sync"

	"github.com/df-mc/dragonfly/server/world"
)

var worldPlayerCount sync.Map

// playerCount returns player count.
func playerCount(w *world.World) int {
	val, ok := worldPlayerCount.Load(w)
	if ok {
		return val.(int)
	}
	return 0
}

// addToWorld adds value to world player count.
func addToWorld(w *world.World, val int) (count int) {
	if w == nil {
		return 0
	}
	previous, ok := worldPlayerCount.Load(w)
	if ok {
		count = previous.(int)
	}
	count += val
	if count <= 0 {
		// avoiding memory leak.
		worldPlayerCount.Delete(w)
	}
	worldPlayerCount.Store(w, count)
	return count
}

// updateWorld updates player count in worlds and returns current count.
func (m *Monitoring) updateWorld(current *world.World) int {
	if m.previousWorld != current {
		addToWorld(m.previousWorld, -1)
		m.previousWorld = current
		return addToWorld(current, 1)
	}
	return playerCount(current)
}
