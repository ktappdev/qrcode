package main

import (
	"image"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ktappdev/qrcode-server/ratelimiter"
	"github.com/ktappdev/qrcode-server/urlexchanger"
	_ "image/jpeg"
	_ "image/png"

	"golang.org/x/image/bmp"
)

var limiter = ratelimiter.NewIPRateLimiter(1)
var exchanger = urlexchanger.NewURLExchanger()

func init() {
	image.RegisterFormat("bmp", "bmp", bmp.Decode, bmp.DecodeConfig)
}

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
	// Set the maximum memory limit for parsing multipart forms
	router.MaxMultipartMemory = 10 << 20 // 10MB

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*", "https://592code.vercel.app"}
	config.AllowMethods = []string{"GET", "POST"}
	router.Use(cors.New(config))
	router.POST("/qrcode", GetQr)
	router.GET("/qr", exchanger.HandleQRCodeInteraction)

	router.Run(":" + port)
}
