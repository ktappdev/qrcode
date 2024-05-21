package main

import (
	"image"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ktappdev/qrcode-server/helpers"
	"github.com/ktappdev/qrcode-server/qrcode"
	"github.com/ktappdev/qrcode-server/ratelimiter"
	"github.com/ktappdev/qrcode-server/urlexchanger"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
)

var limiter = ratelimiter.NewIPRateLimiter(1)
var exchanger = urlexchanger.NewURLExchanger()

// var cachedLogo *image.Image
func init() {
	image.RegisterFormat("bmp", "bmp", bmp.Decode, bmp.DecodeConfig)
	image.RegisterFormat("tiff", "tiff", tiff.Decode, tiff.DecodeConfig)
	image.RegisterFormat("webp", "webp", webp.Decode, webp.DecodeConfig)
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
	backgroundColour := c.PostForm("backgroundColour")
	qrCodeColour := c.PostForm("qrCodeColour")

	// Get the logo image from the form data
	logoFile, err := c.FormFile("logo")
	if err != nil && err != http.ErrMissingFile {
		log.Println("no logo", err)
	}

	// Get the opacity from the form data
	opacity := c.PostForm("opacity")

	// Convert the opacity to a float64 using the utility function
	opacityFloat64, err := helpers.ParseOpacity(opacity)
	if err != nil {
		log.Println("invalid opacity value:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid opacity value"})
		return
	}
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
			log.Println("Failed to decode logo image:", err)
			c.String(http.StatusBadRequest, "Failed to decode logo image")
			return
		}
		logo = &decodedLogo
		// cachedLogo = &decodedLogo
		// logo = cachedLogo
	}

	size := -10 // -10 will make each qr pixel 10x10, i can do 256 which would give 256x256px image but there is usually white space around it
	qrCodeURL := exchanger.GenerateQRCodeURL(originalLink)
	qrCodeBytes, err := qrcode.GenerateQRCode(qrCodeURL, size, qrCodeColour, backgroundColour, logo, opacityFloat64)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error generating QR code")
		return
	}

	c.Data(http.StatusOK, "image/png", qrCodeBytes)
}
