package monitoring

import (
	"github.com/df-mc/dragonfly/server/session"
	"github.com/google/uuid"
)

// hookConn hooks Conn interface in order to call deletePlayer automatically.
func hookConn(s *session.Session, uuid uuid.UUID) {
	conn := fetchPrivateField(s, "conn").Interface().(session.Conn)
	updatePrivateField(s, "conn", connection{conn, uuid})
}

// connection ...
type connection struct {
	session.Conn
	uuid uuid.UUID
}

// Close ...
func (c connection) Close() error {
	deletePlayer(c.uuid)
	return c.Conn.Close()
}
