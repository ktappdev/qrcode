package routehandlers

import (
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ktappdev/qrcode-server/helpers"
	"github.com/ktappdev/qrcode-server/qrcode"
	"github.com/ktappdev/qrcode-server/urlexchanger"
)

var exchanger = urlexchanger.NewURLExchanger()

func GetQr(c *gin.Context) {

	// Get the original link from the form data
	originalLink := c.PostForm("originalLink")
	if originalLink == "https://" {
		originalLink = ""
	}
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
	logo, err := helpers.LoadLogo(c, false) // NOTE: Plugin system for effects later
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to load logo image")
		return
	}

	size := 512 // -10 will make each qr pixel 10x10, i can do 256 which would give 256x256px image but there is usually white space around it
	qrCodeURL := exchanger.GenerateQRCodeURL(originalLink, backgroundColour, qrCodeColour, name)
	qrCodeBytes, err := qrcode.GenerateQRCode(qrCodeURL, size, qrCodeColour, backgroundColour, logo, opacityFloat64)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error generating QR code")
		return
	}

	c.Data(http.StatusOK, "image/png", qrCodeBytes)

}
