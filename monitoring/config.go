package monitoring

import "github.com/df-mc/dragonfly/server/session"

// Config provides a way to customize monitoring settings.
var Config struct {
	// Recovery recovers from panic, if not specified default will be used.
	Recovery func(s *session.Session, c session.Controllable, err error)
}
