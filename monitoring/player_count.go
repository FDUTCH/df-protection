package monitoring

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
	"github.com/puzpuzpuz/xsync/v4"
)

var worldPlayerCount = xsync.NewMap[*world.World, *xsync.Map[uuid.UUID, struct{}]]()

// playerCount returns player count.
func playerCount(w *world.World) int {
	val, ok := worldPlayerCount.Load(w)
	if ok {
		return val.Size()
	}
	return 0
}

// addToWorld adds value to world player count.
func addToWorld(w *world.World, val uuid.UUID) (count int) {
	if w == nil {
		return 0
	}
	players, ok := worldPlayerCount.Load(w)
	if !ok {
		players = xsync.NewMap[uuid.UUID, struct{}]()
		worldPlayerCount.Store(w, players)
	}
	players.Store(val, struct{}{})
	return players.Size()
}

// deleteFromWorld removes value from world player count.
func deleteFromWorld(w *world.World, val uuid.UUID) (count int) {
	if w == nil {
		return 0
	}
	players, ok := worldPlayerCount.Load(w)
	if !ok {
		return 0
	}
	players.Delete(val)
	count = players.Size()
	if count == 0 {
		worldPlayerCount.Delete(w)
	}
	return count
}

// updateWorld updates player count in worlds and returns current count.
func (m *Monitoring) updateWorld(current *world.World, uuid uuid.UUID) int {
	if m.previousWorld != current {
		deleteFromWorld(m.previousWorld, uuid)
		m.previousWorld = current
		return addToWorld(current, uuid)
	}
	return playerCount(current)
}
