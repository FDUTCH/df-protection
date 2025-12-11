package monitoring

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Config provides a way to customize monitoring settings.
var Config struct {
	// Recovery recovers from panic, if not specified default will be used.
	Recovery func(s *session.Session, pk packet.Packet, p *player.Player, err error)
	// PreventLags will disconnect or report players that are causing more lags than the server could handle into PerformanceReporter.
	// Enabling this feature may limit the number of players in your worlds.
	PreventLags bool
	// PerformanceReporter reports players that are causing lags,
	// if not specified player will be disconnected.
	PerformanceReporter func(s *session.Session, p *player.Player, pk packet.Packet)
}
