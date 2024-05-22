package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ktappdev/qrcode-server/helpers"
	"github.com/ktappdev/qrcode-server/qrcode"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
)

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

	// Convert the opacity to a float64
	opacityFloat64, err := helpers.ParseOpacity(opacity)
	if err != nil {
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
		logoFile = nil
		defer file.Close()

		// Decode the logo image
		decodedLogo, _, err := image.Decode(file)
		if err != nil {
			c.String(http.StatusBadRequest, "Failed to decode logo image")
			return
		}
		logo = &decodedLogo
		// cachedLogo = &decodedLogo
		// logo = cachedLogo
	}

	size := 256 // -10 will make each qr pixel 10x10, i can do 256 which would give 256x256px image but there is usually white space around it
	qrCodeURL := exchanger.GenerateQRCodeURL(originalLink)
	qrCodeBytes, err := qrcode.GenerateQRCode(qrCodeURL, size, qrCodeColour, backgroundColour, logo, opacityFloat64)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error generating QR code")
		return
	}

	c.Data(http.StatusOK, "image/png", qrCodeBytes)
}
