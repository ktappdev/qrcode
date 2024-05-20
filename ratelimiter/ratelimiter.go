package ratelimiter

import (
	"sync"

	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ipLimiters map[string]*rate.Limiter // A map to store rate limiters for each IP address
	mutex      sync.Mutex               // A mutex to synchronize access to the ipLimiters map
	RateLimit  rate.Limit               // The rate limit (requests per second)
}

// NewIPRateLimiter creates a new instance of IPRateLimiter with the given rate limit.
func NewIPRateLimiter(rateLimit rate.Limit) *IPRateLimiter {
	return &IPRateLimiter{
		ipLimiters: make(map[string]*rate.Limiter), // Initialize the ipLimiters map
		RateLimit:  rateLimit,                      // Set the rate limit
	}
}

// Allow checks if a request from the specified IP address is allowed based on the rate limit.
// It returns true if the request is allowed and false if it should be rejected.
func (limiter *IPRateLimiter) Allow(ipAddress string) bool {
	limiter.mutex.Lock()         // Acquire the lock to safely access the ipLimiters map
	defer limiter.mutex.Unlock() // Release the lock when the function returns

	if _, exists := limiter.ipLimiters[ipAddress]; !exists {
		// If no rate limiter exists for the IP address, create a new one
		limiter.ipLimiters[ipAddress] = rate.NewLimiter(limiter.RateLimit, 1)
	}

	// Check if the request is allowed by the rate limiter for the IP address
	return limiter.ipLimiters[ipAddress].Allow()
}
