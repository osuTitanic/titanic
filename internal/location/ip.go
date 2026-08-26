package location

import "net/netip"

// IsLocalIP reports whether the given IP address is a local / non-routable address.
func IsLocalIP(ip string) bool {
	address, err := netip.ParseAddr(ip)
	if err != nil {
		return false
	}
	address = address.Unmap()

	return address.IsPrivate() ||
		address.IsLoopback() ||
		address.IsLinkLocalUnicast() ||
		address.IsLinkLocalMulticast() ||
		address.IsMulticast() ||
		address.IsUnspecified()
}
