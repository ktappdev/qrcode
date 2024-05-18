package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ktappdev/qrcode/qrcode"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

var limiter = NewIPRateLimiter(1)

func main() {
	router := gin.Default()

	router.GET("/qrcode", func(c *gin.Context) {
		clientIP := c.ClientIP()
		if !limiter.Allow(clientIP) {
			c.String(http.StatusOK, "Slowdown cowboy!, %v request per second", limiter.rateLimit)
			return
		}
		data := "https://www.maad97.com"
		size := -10 // -10 will make each qr pixel 10x10, i can do 256 which would give 256x256px image but there is usually white space around it

		qrCodeBytes, err := qrcode.GenerateQRCode(data, size)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error generating QR code")
			return
		}

		c.Data(http.StatusOK, "image/png", qrCodeBytes)
	})

	router.Run(":8081")
}

type IPRateLimiter struct {
	ipLimiters map[string]*rate.Limiter // A map to store rate limiters for each IP address
	mutex      sync.Mutex               // A mutex to synchronize access to the ipLimiters map
	rateLimit  rate.Limit               // The rate limit (requests per second)
}

// NewIPRateLimiter creates a new instance of IPRateLimiter with the given rate limit.
func NewIPRateLimiter(rateLimit rate.Limit) *IPRateLimiter {
	return &IPRateLimiter{
		ipLimiters: make(map[string]*rate.Limiter), // Initialize the ipLimiters map
		rateLimit:  rateLimit,                      // Set the rate limit
	}
}

// Allow checks if a request from the specified IP address is allowed based on the rate limit.
// It returns true if the request is allowed and false if it should be rejected.
func (limiter *IPRateLimiter) Allow(ipAddress string) bool {
	limiter.mutex.Lock()         // Acquire the lock to safely access the ipLimiters map
	defer limiter.mutex.Unlock() // Release the lock when the function returns

	if _, exists := limiter.ipLimiters[ipAddress]; !exists {
		// If no rate limiter exists for the IP address, create a new one
		limiter.ipLimiters[ipAddress] = rate.NewLimiter(limiter.rateLimit, 1)
	}

	// Check if the request is allowed by the rate limiter for the IP address
	return limiter.ipLimiters[ipAddress].Allow()
}
