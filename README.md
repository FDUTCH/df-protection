# Protection for Dragonfly

Is the utility library utility for Dragonfly that improves stability, protects against malicious packets, and
adds several useful utilities.

- [x] Recovery
- [x] Overload protection.
- [ ] PPS Limit
- [ ] Protection against malicious packets

## Example

```go
// enable antilag protection
monitoring.Config.PreventLags = true

monitoring.Config.PerformanceReporter = func(s *session.Session, c session.Controllable) {
	// your actions against server lag
}

// customize your Recovery 
monitoring.Config.Recovery = func(s *session.Session, c session.Controllable, err error) {
	// your recovery
}

for p := range srv.Accept() {
    monitoring.Enable(p)
}
```