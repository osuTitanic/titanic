package location

import "net"

// IsLocalIP reports whether the given IP address is a local / non-routable address.
func IsLocalIP(ip string) bool {
	address := net.ParseIP(ip)
	if address == nil {
		return false
	}

	return address.IsPrivate() ||
		address.IsLoopback() ||
		address.IsLinkLocalUnicast() ||
		address.IsLinkLocalMulticast() ||
		address.IsMulticast() ||
		address.IsUnspecified()
}
