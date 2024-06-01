package routehandlers

import (
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ktappdev/qrcode-server/generator"
	"github.com/ktappdev/qrcode-server/helpers"
	"github.com/ktappdev/qrcode-server/urlhandler"
)

var exchanger = urlhandler.NewURLExchanger()

func GetQr(c *gin.Context) {
	// Parse form data into FormData struct
	formData := helpers.FormDataStruct{
		OriginalLink:     c.PostForm("originalLink"),
		Opacity:          c.PostForm("opacity"),
		BackgroundColour: c.PostForm("backgroundColour"),
		QRCodeColour:     c.PostForm("qrCodeColour"),
		Name:             c.PostForm("name"),
		UseDots:          c.PostForm("useDots"),
		OverlayOurLogo:   c.PostForm("overlayOurLogo"),
	}

	if formData.OriginalLink == "https://" {
		formData.OriginalLink = ""
	}

	// Convert the opacity to a float64
	opacityFloat64, err := helpers.ParseOpacity(formData.Opacity)
	formData.Opacityf64 = opacityFloat64
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid opacity value"})
		return
	}

	// Load logo image
	logo, err := helpers.LoadLogo(c, false) // NOTE: Plugin system for effects later - false no effect
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to load logo image")
		return
	}

	size := 512 // -10 will make each qr pixel 10x10, i can do 256 which would give 256x256px image but there is usually white space around it
	qrCodeURL := exchanger.GenerateQRCodeURL(&formData)
	log.Println("qrCodeURL in getQr: ", qrCodeURL)
	qrCodeBytes, err := generator.GenerateQRCode(
		qrCodeURL,
		size,
		&formData,
		logo,
		0,
	)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error generating QR code")
		return
	}

	c.Data(http.StatusOK, "image/png", qrCodeBytes)
}
