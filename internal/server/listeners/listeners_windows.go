//go:build windows

package listeners

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/Microsoft/go-winio"

	"github.com/xmapst/osreapi/internal/sockets"
)

var PipeName = "npipe:////./pipe/osreapi"

// Init creates new listeners for the server.
func Init(proto, addr string, tlsConfig *tls.Config) (net.Listener, error) {
	switch proto {
	case "tcp":
		return sockets.NewTCPSocket(addr, tlsConfig)
	case "npipe":
		// allow Administrators and SYSTEM, plus whatever additional users or groups were specified
		sddl := "D:P(A;;GA;;;BA)(A;;GA;;;SY)"
		c := winio.PipeConfig{
			SecurityDescriptor: sddl,
			MessageMode:        true,  // Use message mode so that CloseWrite() is supported
			InputBufferSize:    65536, // Use 64KB buffers to improve performance
			OutputBufferSize:   65536,
		}
		return winio.ListenPipe(addr, &c)
	default:
		return nil, fmt.Errorf("invalid protocol format: windows only supports tcp and npipe")
	}
}
