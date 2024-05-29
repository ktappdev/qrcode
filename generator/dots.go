package generator

import (
	"bytes"
	"image/color"
	"image/png"

	"github.com/fogleman/gg"
	"github.com/skip2/go-qrcode"
)

// drawQRCodeWithDots draws the QR code with dots
func drawQRCodeWithDots(qr *qrcode.QRCode, size int, foregroundColor color.Color, backgroundColor color.Color) ([]byte, error) {
	dotSize := size / len(qr.Bitmap())
	dc := gg.NewContext(size, size)

	// Set background color
	dc.SetColor(backgroundColor)
	dc.Clear()

	// Set foreground color
	dc.SetColor(foregroundColor)

	// Draw the QR code with dots
	matrix := qr.Bitmap()
	for y := 0; y < len(matrix); y++ {
		for x := 0; x < len(matrix); x++ {
			if matrix[y][x] {
				dc.DrawCircle(float64(x*dotSize+dotSize/2), float64(y*dotSize+dotSize/2), float64(dotSize/2))
				dc.Fill()
			}
		}
	}

	// Encode the image to PNG
	buf := &bytes.Buffer{}
	err := png.Encode(buf, dc.Image())
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
