package server

import (
	"net"
	"strings"
)

// GetClientAddress truncates client port from address.
func GetClientAddress(addr net.Addr) string {
	strs := strings.Split(addr.String(), ":")
	if len(strs) > 0 {
		return strs[0]
	}

	return ""
}
