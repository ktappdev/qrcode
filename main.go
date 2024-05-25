package main

import (
	"image"
	"log"
	"net/http"
	"os"

	_ "image/jpeg"
	_ "image/png"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ktappdev/qrcode-server/mongodb"
	"github.com/ktappdev/qrcode-server/ratelimiter"
	"github.com/ktappdev/qrcode-server/routehandlers"

	"golang.org/x/image/bmp"
)

var limiter = ratelimiter.NewIPRateLimiter(1)

func init() {
	image.RegisterFormat("bmp", "bmp", bmp.Decode, bmp.DecodeConfig)
}

func initMongoDB(MONGO_URL string) {
	err := mongodb.Connect(MONGO_URL)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB Atlas: %v", err)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	MONGO_URL := os.Getenv("MONGO_URL")
	if MONGO_URL == "" {
		log.Fatal("MONGO_URL environment variable is not set")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8081"
	}

	initMongoDB(MONGO_URL)
	router := gin.Default()
	// Set the maximum memory limit for parsing multipart forms
	router.MaxMultipartMemory = 10 << 20 // 10MB

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*", "http://192.168.1.13", "https://592code.vercel.app"}
	config.AllowMethods = []string{"GET", "POST"}
	router.Use(cors.New(config))

	// Apply rate limiter middleware to all routes
	router.Use(func(c *gin.Context) {
		clientIP := c.ClientIP()
		if !limiter.Allow(clientIP) {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	})

	router.POST("/qrcode", routehandlers.GetQr)
	router.GET("/qr", routehandlers.HandleScan)
	router.GET("/qrcode-details", mongodb.GetInteractionsForQRCode)

	router.Run(":" + port)
}
