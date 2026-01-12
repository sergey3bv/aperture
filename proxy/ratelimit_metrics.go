package proxy

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// rateLimitAllowed counts requests that passed rate limiting.
	rateLimitAllowed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "aperture",
			Subsystem: "ratelimit",
			Name:      "allowed_total",
			Help:      "Total number of requests allowed by rate limiter",
		},
		[]string{"service", "path_pattern"},
	)

	// rateLimitDenied counts requests denied by rate limiting.
	rateLimitDenied = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "aperture",
			Subsystem: "ratelimit",
			Name:      "denied_total",
			Help:      "Total number of requests denied by rate limiter",
		},
		[]string{"service", "path_pattern"},
	)

	// rateLimitCacheSize tracks the current size of the rate limiter cache.
	rateLimitCacheSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "aperture",
			Subsystem: "ratelimit",
			Name:      "cache_size",
			Help:      "Current number of entries in the rate limiter cache",
		},
		[]string{"service"},
	)

	// rateLimitEvictions counts LRU cache evictions.
	rateLimitEvictions = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "aperture",
			Subsystem: "ratelimit",
			Name:      "evictions_total",
			Help:      "Total number of rate limiter cache evictions",
		},
		[]string{"service"},
	)
)
