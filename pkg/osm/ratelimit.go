// Package osm provides utilities for working with OpenStreetMap data.
package osm

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	// Service names for rate limiting
	ServiceNominatim = "nominatim"
	ServiceOverpass  = "overpass"
	ServiceOSRM      = "osrm"
)

// RateLimiter manages rate limiting for different OpenStreetMap API services
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
}

var (
	// globalRateLimiter is the singleton rate limiter instance
	globalRateLimiter *RateLimiter

	// rateLimiterOnce ensures we only create the rate limiter once
	rateLimiterOnce sync.Once
)

// GetRateLimiter returns the global rate limiter instance
func GetRateLimiter() *RateLimiter {
	rateLimiterOnce.Do(func() {
		// Initialize the global rate limiter with service-specific limits
		// according to OSM API usage policies
		limiters := make(map[string]*rate.Limiter)

		// Nominatim: 1 request per second
		// https://operations.osmfoundation.org/policies/nominatim/
		limiters[ServiceNominatim] = rate.NewLimiter(rate.Every(1*time.Second), 1)

		// Overpass: 2 requests per minute with bursts of up to 2 requests
		// https://wiki.openstreetmap.org/wiki/Overpass_API#Public_Overpass_API_instances
		limiters[ServiceOverpass] = rate.NewLimiter(rate.Every(30*time.Second), 2)

		// OSRM: 100 requests per minute (to be safe and avoid abuse)
		limiters[ServiceOSRM] = rate.NewLimiter(rate.Every(600*time.Millisecond), 5)

		globalRateLimiter = &RateLimiter{
			limiters: limiters,
		}
	})

	return globalRateLimiter
}

// Wait blocks until the rate limit for the specified service allows an event
// or the context is canceled.
func (rl *RateLimiter) Wait(ctx context.Context, service string) error {
	rl.mu.RLock()
	limiter, exists := rl.limiters[service]
	rl.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no rate limiter defined for service: %s", service)
	}

	// Wait for rate limiter or context cancellation
	err := limiter.Wait(ctx)
	if err != nil {
		slog.Debug("rate limiter wait error", "service", service, "error", err)
		return err
	}

	return nil
}

// WaitForService is a convenience function to wait for a service's rate limit
// using the global rate limiter
func WaitForService(ctx context.Context, service string) error {
	return GetRateLimiter().Wait(ctx, service)
}
