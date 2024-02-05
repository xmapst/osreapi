//go:build !windows

package listeners

import (
	"crypto/tls"
	"net"
	"os"

	"github.com/pkg/errors"

	"github.com/xmapst/osreapi/internal/sockets"
)

var PipeName = "unix:///var/run/osreapi.sock"

// Init creates new listeners for the server.
func Init(proto, addr string, tlsConfig *tls.Config) (net.Listener, error) {
	switch proto {
	case "tcp":
		return sockets.NewTCPSocket(addr, tlsConfig)
	case "unix":
		return sockets.NewUnixSocket(addr, os.Getegid())
	default:
		return nil, errors.Errorf("invalid protocol format: %q", proto)
	}
}
