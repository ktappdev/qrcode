package qrcode

import (
	"bytes"
	"image"
	"image/png"

	"github.com/ktappdev/qrcode-server/helpers"
	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(data string, size int, qrCodeColour string, backgroundColour string, logo *image.Image, opacity float64) ([]byte, error) {
	var qr *qrcode.QRCode
	if logo != nil {
		qr, _ = qrcode.New(data, qrcode.High) // NOTE: can also set Highest
	} else {
		qr, _ = qrcode.New(data, qrcode.Low)
	}
	bgc, qrc := helpers.SetColours(backgroundColour, qrCodeColour)

	// Set the foreground and background colors
	qr.ForegroundColor = qrc
	qr.BackgroundColor = bgc

	// Generate the QR code image
	qrImg := qr.Image(size)

	// Overlay the logo image if provided
	if logo != nil {
		helpers.OverlayLogo(&qrImg, *logo, opacity)
	}

	buf := &bytes.Buffer{}
	err := png.Encode(buf, qrImg)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
