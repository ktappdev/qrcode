package generator

import (
	"image"
	"image/png"
	"os"

	"github.com/ktappdev/qrcode-server/helpers"
)

func overlayOurLogoFunc(qrImg image.Image, overlayShrink int) (image.Image, error) {
	// Open the logo file
	//TODO: make this dynamic when i have other qr sizes
	file, err := os.Open("assets/512_logo_overlay.png")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the logo image
	logo, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	// Overlay the logo on the QR code image
	helpers.OverlayLogo(&qrImg, logo, 1.0, overlayShrink) // Adjust the opacity as needed

	return qrImg, nil
}
