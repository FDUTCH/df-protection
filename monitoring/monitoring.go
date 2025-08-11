package monitoring

import (
	"container/list"
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
)

var (
	monitorings = make(map[uuid.UUID]*Monitoring)
	mu          sync.RWMutex
)

// Monitoring stores execution time.
type Monitoring struct {
	timeSpent time.Duration
	calls     *list.List

	previousWorld *world.World
}

// callRecord is the record of the execution time of the handler call.
type callRecord struct {
	duration time.Duration
	endTime  time.Time
}

// newMonitoring ...
func newMonitoring() *Monitoring {
	return &Monitoring{
		calls: list.New(),
	}
}

// ExecutionTimeForLastSecond returns the processor time spent processing player packets over the last second.
func (m *Monitoring) ExecutionTimeForLastSecond() time.Duration {
	return m.executionTimeSince(time.Second)
}

// executionTimeSince returns the processor time spent processing player packets over a specified period of time.
func (m *Monitoring) executionTimeSince(period time.Duration) time.Duration {
	cutoff := time.Now().Add(-period)
	var total time.Duration
	for e := m.calls.Front(); e != nil; e = e.Next() {
		record := e.Value.(callRecord)
		if record.endTime.Before(cutoff) {
			break
		}
		total += record.duration
	}
	return total
}

// ExecutionTime returns total processor time spent processing player packets.
func (m *Monitoring) ExecutionTime() time.Duration {
	return m.timeSpent
}

// writeExecutionTime writes callRecord into Monitoring.
func (m *Monitoring) writeExecutionTime(start, end time.Time) {
	duration := end.Sub(start)
	m.timeSpent += duration
	m.calls.PushFront(callRecord{
		duration: duration,
		endTime:  end,
	})
	m.gc()
}

// gc clears all outdated records.
func (m *Monitoring) gc() {
	cutoff := time.Now().Add(-time.Second)
	for {
		last := m.calls.Back()
		if last == nil {
			return
		}
		if last.Value.(callRecord).endTime.Before(cutoff) {
			m.calls.Remove(last)
		} else {
			return
		}
	}
}

// hook hooks original handler.
func (m *Monitoring) hook(h packetHandler) packetHandler {
	return &hook{original: h, m: m}
}

// deletePlayer deletes player from the internal map.
func deletePlayer(uuid uuid.UUID) {
	mu.Lock()
	previousWorld := monitorings[uuid].previousWorld
	delete(monitorings, uuid)
	mu.Unlock()
	addToWorld(previousWorld, -1)
}

// GetMonitoring returns player's Monitoring.
func GetMonitoring(pl *player.Player) *Monitoring {
	mu.RLock()
	defer mu.RUnlock()
	return monitorings[pl.UUID()]
}
