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
func (h *hook) Handle(p packet.Packet, s *session.Session, tx *world.Tx, c session.Controllable) error {
	defer func() {
		err := recover()
		if err != nil {
			if Config.Recovery != nil {
				Config.Recovery(s, p, c.(*player.Player), err.(error))
			} else {
				s.Disconnect(fmt.Sprintf("error proccessing packet: %T handler: %T err: %v", p, h.original, err.(error)))
			}
		}
	}()

	start := time.Now()
	// calling original handler.
	err := h.original.Handle(p, s, tx, c)
	end := time.Now()

	h.m.writeExecutionTime(start, end)

	if Config.PreventLags && time.Second/time.Duration(h.m.updateWorld(tx.World())) < h.m.ExecutionTimeForLastSecond() {
		if Config.PerformanceReporter != nil {
			Config.PerformanceReporter(s, c.(*player.Player))
		} else {
			s.Disconnect(fmt.Sprintf("disconnected due to server performance issues"))
		}
	}
	return err
}
