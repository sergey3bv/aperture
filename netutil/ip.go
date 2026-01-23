package netutil

import "net"

var (
	// ipv4Mask24 masks IPv4 addresses to /24 (last octet zeroed).
	// This groups clients on the same subnet together.
	ipv4Mask24 = net.CIDRMask(24, 32)

	// ipv6Mask48 masks IPv6 addresses to /48.
	// Residential connections typically receive /48 to /64 allocations,
	// so /48 provides reasonable grouping for rate limiting purposes.
	ipv6Mask48 = net.CIDRMask(48, 128)
)

// MaskIP returns a masked version of the IP address for grouping purposes.
// IPv4 addresses are masked to /24 (zeroing the last octet).
// IPv6 addresses are masked to /48.
//
// This is useful for rate limiting and freebie tracking where we want to
// group requests from the same network segment rather than individual IPs,
// reducing abuse potential from users with multiple addresses.
func MaskIP(ip net.IP) net.IP {
	if ip4 := ip.To4(); ip4 != nil {
		return ip4.Mask(ipv4Mask24)
	}

	return ip.Mask(ipv6Mask48)
}
