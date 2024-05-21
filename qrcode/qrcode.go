package qrcode

import (
	"bytes"
	"image"
	"image/color"
	"image/png"

	"github.com/ktappdev/qrcode-server/helpers"
	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(data string, size int, foregroundColor, backgroundColor color.Color, logo *image.Image, opacity float64) ([]byte, error) {
	qr, err := qrcode.New(data, qrcode.High) // NOTE: can also set Highest
	if err != nil {
		return nil, err
	}

	// Set the foreground and background colors
	qr.ForegroundColor = foregroundColor
	qr.BackgroundColor = backgroundColor

	// Generate the QR code image
	qrImg := qr.Image(size)

	// Overlay the logo image if provided
	if logo != nil {
		helpers.OverlayLogo(&qrImg, *logo, opacity)
	}

	buf := &bytes.Buffer{}
	err = png.Encode(buf, qrImg)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
