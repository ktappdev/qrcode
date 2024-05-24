package main

import (
	"image"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ktappdev/qrcode-server/helpers"
	"github.com/ktappdev/qrcode-server/qrcode"
)

func GetWifiQR(c *gin.Context) {
	clientIP := c.ClientIP()

	if !limiter.Allow(clientIP) {
		c.String(http.StatusOK, "Slowdown cowboy!, %v request per second", limiter.RateLimit)
		return
	}

	// Get the Wi-Fi network information from the form data
	ssid := c.PostForm("ssid")
	password := c.PostForm("password")
	encryption := c.PostForm("encryption")
	isHidden := c.PostForm("isHidden") == "true"

	// Create the Wi-Fi QR code data string
	wifiQRData := "WIFI:T:" + encryption + ";S:" + ssid + ";P:" + password + ";H:" + strconv.FormatBool(isHidden) + ";;"

	// Get the QR code color and background color from the form data
	qrCodeColour := c.PostForm("qrCodeColour")
	backgroundColour := c.PostForm("backgroundColour")
	name := c.PostForm("name")

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
		defer file.Close()

		// Decode the logo image
		decodedLogo, _, err := image.Decode(file)
		if err != nil {
			c.String(http.StatusBadRequest, "Failed to decode logo image")
			return
		}

		logo = &decodedLogo
	}

	size := 256
	qrCodeURL := exchanger.GenerateQRCodeURL(wifiQRData, backgroundColour, qrCodeColour, name)
	qrCodeBytes, err := qrcode.GenerateQRCode(qrCodeURL, size, qrCodeColour, backgroundColour, logo, opacityFloat64)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error generating QR code")
		return
	}

	c.Data(http.StatusOK, "image/png", qrCodeBytes)
}
