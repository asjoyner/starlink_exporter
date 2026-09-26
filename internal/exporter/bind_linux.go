//go:build linux

package exporter

import (
	"syscall"
)

// bindToDevice returns a net.Dialer Control func that pins the socket to
// iface with SO_BINDTODEVICE. Routing then only considers routes on that
// interface, so the dish (192.168.100.1) is reachable over a backup WAN
// without a static route, even while another WAN holds the default route.
// Unprivileged since Linux 5.7 for sockets that are not yet bound.
func bindToDevice(iface string) func(network, address string, c syscall.RawConn) error {
	return func(network, address string, c syscall.RawConn) error {
		var serr error
		if err := c.Control(func(fd uintptr) {
			serr = syscall.SetsockoptString(int(fd), syscall.SOL_SOCKET, syscall.SO_BINDTODEVICE, iface)
		}); err != nil {
			return err
		}
		return serr
	}
}
