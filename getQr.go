package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ktappdev/qrcode-server/helpers"
	"github.com/ktappdev/qrcode-server/qrcode"
	_ "image/jpeg"
	_ "image/png"
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
	if originalLink == "https://" {
		originalLink = ""
	}
	// Get the opacity from the form data
	opacity := c.PostForm("opacity")
	backgroundColour := c.PostForm("backgroundColour")
	qrCodeColour := c.PostForm("qrCodeColour")
	name := c.PostForm("name")
	// Get the logo image from the form data

	// Convert the opacity to a float64
	opacityFloat64, err := helpers.ParseOpacity(opacity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid opacity value"})
		return
	}
	logo, err := LoadLogo(c)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to load logo image")
		return
	}

	size := 256 // -10 will make each qr pixel 10x10, i can do 256 which would give 256x256px image but there is usually white space around it
	qrCodeURL := exchanger.GenerateQRCodeURL(originalLink, backgroundColour, qrCodeColour, name)
	qrCodeBytes, err := qrcode.GenerateQRCode(qrCodeURL, size, qrCodeColour, backgroundColour, logo, opacityFloat64)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error generating QR code")
		return
	}

	c.Data(http.StatusOK, "image/png", qrCodeBytes)

}
