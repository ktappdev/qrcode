package qrcode

import (
	"bytes"
	"image"
	"image/color"
	"image/png"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(data string, size int, foregroundColor, backgroundColor color.Color, logo *image.Image) ([]byte, error) {
	qr, err := qrcode.New(data, qrcode.Medium)
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
		overlayLogo(&qrImg, *logo)
	}

	buf := &bytes.Buffer{}
	err = png.Encode(buf, qrImg)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func overlayLogo(qrImg *image.Image, logo image.Image) {
	// ... (the overlayLogo function remains the same)
}
