package netutil

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestMaskIP verifies that MaskIP correctly applies /24 masks to IPv4 and /48
// masks to IPv6 addresses.
func TestMaskIP(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "IPv4 masks last octet",
			input:    "192.168.1.123",
			expected: "192.168.1.0",
		},
		{
			name:     "IPv4 already masked",
			input:    "10.0.0.0",
			expected: "10.0.0.0",
		},
		{
			name:     "IPv4 different last octet same result",
			input:    "192.168.1.255",
			expected: "192.168.1.0",
		},
		{
			name:     "IPv6 masks to /48",
			input:    "2001:db8:1234:5678:9abc:def0:1234:5678",
			expected: "2001:db8:1234::",
		},
		{
			name:     "IPv6 already masked",
			input:    "2001:db8:abcd::",
			expected: "2001:db8:abcd::",
		},
		{
			name:     "IPv6 loopback",
			input:    "::1",
			expected: "::",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ip := net.ParseIP(tc.input)
			require.NotNil(t, ip, "failed to parse input IP")

			result := MaskIP(ip)
			require.Equal(t, tc.expected, result.String())
		})
	}
}

// TestMaskIP_SameSubnetGroupsTogether verifies that IPv4 addresses in the same
// /24 subnet produce identical masked results.
func TestMaskIP_SameSubnetGroupsTogether(t *testing.T) {
	// Verify that IPs in the same /24 subnet produce the same masked result.
	ips := []string{
		"192.168.1.1",
		"192.168.1.100",
		"192.168.1.255",
	}

	results := make([]string, 0, len(ips))
	for _, ipStr := range ips {
		ip := net.ParseIP(ipStr)
		results = append(results, MaskIP(ip).String())
	}

	// All should be the same.
	for i := 1; i < len(results); i++ {
		require.Equal(t, results[0], results[i],
			"IPs in same /24 should have same masked result")
	}
}

// TestMaskIP_DifferentSubnetsDiffer verifies that IPv4 addresses in different
// /24 subnets produce distinct masked results.
func TestMaskIP_DifferentSubnetsDiffer(t *testing.T) {
	ip1 := net.ParseIP("192.168.1.100")
	ip2 := net.ParseIP("192.168.2.100")

	result1 := MaskIP(ip1).String()
	result2 := MaskIP(ip2).String()

	require.NotEqual(t, result1, result2,
		"IPs in different /24 subnets should have different masked results")
}

// TestMaskIP_IPv6SamePrefix48GroupsTogether verifies that IPv6 addresses
// sharing the same /48 prefix produce identical masked results.
func TestMaskIP_IPv6SamePrefix48GroupsTogether(t *testing.T) {
	// IPs in the same /48 should produce the same masked result.
	ips := []string{
		"2001:db8:1234:0001::",
		"2001:db8:1234:ffff::",
		"2001:db8:1234:abcd:1234:5678:9abc:def0",
	}

	results := make([]string, 0, len(ips))
	for _, ipStr := range ips {
		ip := net.ParseIP(ipStr)
		results = append(results, MaskIP(ip).String())
	}

	for i := 1; i < len(results); i++ {
		require.Equal(t, results[0], results[i],
			"IPs in same /48 should have same masked result")
	}
}
