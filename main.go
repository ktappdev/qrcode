package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ktappdev/qrcode/qrcode"
	"github.com/ktappdev/qrcode/urlexchanger"
	"golang.org/x/time/rate"
)

var limiter = NewIPRateLimiter(1)
var exchanger = urlexchanger.NewURLExchanger()

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8081"
	}

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*", "http://localhost:3000", "http://localhost:8081", "http://192.168.1.14"}
	config.AllowMethods = []string{"GET"}
	router.Use(cors.New(config))
	router.GET("/qrcode", getQr)
	router.GET("/qr", exchanger.HandleQRCodeInteraction)
	router.Run(port)
}

func getQr(c *gin.Context) {
	clientIP := c.ClientIP()
	if !limiter.Allow(clientIP) {
		c.String(http.StatusOK, "Slowdown cowboy!, %v request per second", limiter.rateLimit)
		return
	}
	// data := "https://www.lugetech.com"
	size := -10 // -10 will make each qr pixel 10x10, i can do 256 which would give 256x256px image but there is usually white space around it
	originalURL := "https://www.lugetech.com"
	qrCodeURL := exchanger.GenerateQRCodeURL(originalURL)
	qrCodeBytes, err := qrcode.GenerateQRCode(qrCodeURL, size)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error generating QR code")
		return
	}

	c.Data(http.StatusOK, "image/png", qrCodeBytes)
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
