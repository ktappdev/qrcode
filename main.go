package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ktappdev/qrcode-server/urlexchanger"
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
	router.GET("/qrcode", GetQr)
	router.GET("/qr", exchanger.HandleQRCodeInteraction)

	router.Run(":" + port)
}
