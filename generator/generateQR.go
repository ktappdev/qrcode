package generator

import (
	"bytes"
	"image"
	"image/png"
	"log"

	"github.com/ktappdev/qrcode-server/helpers"
	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(
	data string,
	size int,
	qrCodeColour string,
	backgroundColour string,
	logo *image.Image,
	opacity float64,
	useDots string,
	overlayOurLogo string,
	overlayShrink int,
) ([]byte, error) {
	var qr *qrcode.QRCode
	var qrBytes []byte
	var qrImg image.Image
	if logo != nil {
		qr, _ = qrcode.New(data, qrcode.High) // NOTE: can also set Highest
		log.Println("Logo is present in the generate function")
	} else {
		qr, _ = qrcode.New(data, qrcode.Low)
	}

	bgcHEX, qrcHEX := helpers.SetColours(backgroundColour, qrCodeColour)
	foregroundColor, err := helpers.HexToColor(qrcHEX)
	if err != nil {
		// Handle the error
		return nil, err
	}
	backgroundColor, err := helpers.HexToColor(bgcHEX)
	if err != nil {
		// Handle the error
		return nil, err
	}

	qr.ForegroundColor = foregroundColor
	qr.BackgroundColor = backgroundColor

	if useDots == "true" {
		qrBytes, err := drawQRCodeWithDots(qr, size, foregroundColor, backgroundColor)
		if err != nil {
			return nil, err
		}
		qrImg, _, err = image.Decode(bytes.NewReader(qrBytes))
		if err != nil {
			return nil, err
		}
	} else {
		// Generate the QR code image
		qrImg = qr.Image(size)
	}

	// Overlay the logo image if provided
	if logo != nil {
		helpers.OverlayLogo(&qrImg, *logo, opacity, 3)
	}

	// Overlay your logo if overLayOurLogo is true
	if overlayOurLogo == "true" {
		qrImg, err = overlayOurLogoFunc(qrImg, 1)
		if err != nil {
			return nil, err
		}
	}

	buf := &bytes.Buffer{}
	err = png.Encode(buf, qrImg)
	if err != nil {
		return nil, err
	}
	qrBytes = buf.Bytes()

	return qrBytes, nil
}
