# Protection for Dragonfly

Is the utility library utility for Dragonfly that improves stability, protects against malicious packets, and
adds several useful utilities.

- [x] Recovery
- [ ] PPS Limit
- [ ] Overload protection.
- [ ] Protection against malicious packets

## Example

```go
for p := range srv.Accept() {
    monitoring.Enable(p)
}
```