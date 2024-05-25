package routehandlers

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"strconv"

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
	logo, err := helpers.LoadLogo(c)
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

func GetWifiQR(c *gin.Context) {

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
