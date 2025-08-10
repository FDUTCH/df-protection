package monitoring

import (
	"reflect"

	u "github.com/bedrock-gophers/unsafe/unsafe"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Enable enables monitoring for the player.
func Enable(pl *player.Player) {
	sess := u.Session(pl)
	hookConn(sess, pl.UUID())
	handlers := fetchPrivateField(sess, "handlers")
	r := handlers.MapRange()

	monitoring := newMonitoring()
	mu.Lock()
	monitorings[pl.UUID()] = monitoring
	mu.Unlock()

	for r.Next() {
		value := r.Value().Interface()
		if value == nil {
			continue
		}

		switch value.(type) {
		// we cannot hook these handlers because of their specific use within Dragonfly.
		case
			*session.CommandRequestHandler,
			*session.ItemStackRequestHandler,
			*session.NPCRequestHandler,
			*session.ModalFormResponseHandler,
			*session.ServerBoundLoadingScreenHandler:
			continue
		}

		if handler, ok := value.(packetHandler); ok {
			handlers.SetMapIndex(r.Key(), reflect.ValueOf(monitoring.hook(handler)))
		}
	}
}

// packetHandler ...
type packetHandler interface {
	Handle(p packet.Packet, s *session.Session, tx *world.Tx, c session.Controllable) error
}
