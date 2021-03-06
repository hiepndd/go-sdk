package r2

import (
	"net"
	"time"
)

// OptDialKeepAlive sets the dial keep alive duration.
// Only use this if you know what you're doing, the defaults are typically sufficient.
func OptDialKeepAlive(d time.Duration) DialOption {
	return func(dialer *net.Dialer) {
		dialer.KeepAlive = d
	}
}
