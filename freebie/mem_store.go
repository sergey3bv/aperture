package freebie

import (
	"net"
	"net/http"

	"github.com/lightninglabs/aperture/netutil"
)

type Count uint16

type memStore struct {
	numFreebies    Count
	freebieCounter map[string]Count
}

func (m *memStore) getKey(ip net.IP) string {
	return netutil.MaskIP(ip).String()
}

func (m *memStore) currentCount(ip net.IP) Count {
	counter, ok := m.freebieCounter[m.getKey(ip)]
	if !ok {
		return 0
	}
	return counter
}

func (m *memStore) CanPass(r *http.Request, ip net.IP) (bool, error) {
	return m.currentCount(ip) < m.numFreebies, nil
}

func (m *memStore) TallyFreebie(r *http.Request, ip net.IP) (bool, error) {
	counter := m.currentCount(ip) + 1
	m.freebieCounter[m.getKey(ip)] = counter
	return true, nil
}

// NewMemIPMaskStore creates a new in-memory freebie store that masks IP
// addresses to keep track of free requests. IPv4 addresses are masked to /24
// and IPv6 addresses to /48. This reduces risk of abuse by users that have a
// whole range of IPs at their disposal.
func NewMemIPMaskStore(numFreebies Count) DB {
	return &memStore{
		numFreebies:    numFreebies,
		freebieCounter: make(map[string]Count),
	}
}
