package monitoring

import (
	"fmt"
	"time"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// hook implements packetHandler to hook original handler.
type hook struct {
	original packetHandler
	m        *Monitoring
}

// Handle ...
func (h *hook) Handle(p packet.Packet, s *session.Session, tx *world.Tx, c session.Controllable) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = makeError(e)
			if Config.Recovery != nil {
				Config.Recovery(s, p, c.(*player.Player), err.(error))
			} else {
				s.Disconnect(fmt.Sprintf("error proccessing packet: %T handler: %T err: %v", p, h.original, err.(error)))
			}
		}
	}()

	start := time.Now()
	// calling original handler.
	err = h.original.Handle(p, s, tx, c)
	if err != nil {
		return
	}
	end := time.Now()

	h.m.WriteExecutionTime(start, end)
	count := h.m.updateWorld(tx.World(), c.UUID())
	if count < 1 || !Config.PreventLags {
		return
	}

	execution := h.m.ExecutionTimeForLastSecond()
	// it shouldn't take more than one tick anyway.
	if execution >= time.Millisecond*50 || time.Second/time.Duration(count) < execution {
		if Config.PerformanceReporter != nil {
			Config.PerformanceReporter(s, c.(*player.Player), p)
		} else {
			s.Disconnect(fmt.Sprintf("disconnected due to server performance issues"))
			err = fmt.Errorf("causing too much lag")
		}
	}
	return
}

func makeError(e any) error {
	switch t := e.(type) {
	case error:
		return t
	default:
		return fmt.Errorf("error: %v", e)
	}
}
