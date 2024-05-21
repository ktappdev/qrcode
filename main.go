package main

import (
	"image"
	"image/color"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ktappdev/qrcode-server/qrcode"
	"github.com/ktappdev/qrcode-server/ratelimiter"
	"github.com/ktappdev/qrcode-server/urlexchanger"
)

var limiter = ratelimiter.NewIPRateLimiter(1)
var exchanger = urlexchanger.NewURLExchanger()
var cachedLogo *image.Image

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
	config.AllowMethods = []string{"GET", "POST"}
	router.Use(cors.New(config))
	router.GET("/qrcode", GetQr)
	router.POST("/qrcode", GetQr)
	router.GET("/qr", exchanger.HandleQRCodeInteraction)

	router.Run(":" + port)
}

func GetQr(c *gin.Context) {
	clientIP := c.ClientIP()
	if !limiter.Allow(clientIP) {
		c.String(http.StatusOK, "Slowdown cowboy!, %v request per second", limiter.RateLimit)
		return
	}

	// Get the original link from the form data
	originalLink := c.PostForm("originalLink")
	log.Println("Original link:", originalLink)

	// Get the logo image from the form data
	logoFile, err := c.FormFile("logo")
	if err != nil && err != http.ErrMissingFile {
		log.Println("no logo", err)
	}
	log.Println("YES logo")

	var logo *image.Image
	if logoFile != nil {
		file, err := logoFile.Open()
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to open logo file")
			return
		}
		defer file.Close()

		// Decode the logo image
		decodedLogo, _, err := image.Decode(file)
		if err != nil {
			c.String(http.StatusBadRequest, "Failed to decode logo image")
			return
		}
		cachedLogo = &decodedLogo
		logo = cachedLogo
	} else if cachedLogo != nil {
		logo = cachedLogo
	}

	size := -10        // -10 will make each qr pixel 10x10, i can do 256 which would give 256x256px image but there is usually white space around it
	fgc := color.White //RGBA{R: 255, G: 0, B: 0, A: 255} // Red color
	bgc := color.Black

	qrCodeURL := exchanger.GenerateQRCodeURL(originalLink)
	qrCodeBytes, err := qrcode.GenerateQRCode(qrCodeURL, size, fgc, bgc, logo)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error generating QR code")
		return
	}

	c.Data(http.StatusOK, "image/png", qrCodeBytes)
}
