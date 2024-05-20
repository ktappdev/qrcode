package helpers

import (
	"image"
	"image/draw"
)

func OverlayLogo(qrImg *image.Image, logo image.Image) {
	// Calculate the center position of the QR code
	qrBounds := (*qrImg).Bounds()
	qrWidth := qrBounds.Max.X - qrBounds.Min.X
	qrHeight := qrBounds.Max.Y - qrBounds.Min.Y
	centerX := qrWidth / 2
	centerY := qrHeight / 2

	// Calculate the size of the logo
	logoBounds := logo.Bounds()
	logoWidth := logoBounds.Max.X - logoBounds.Min.X
	logoHeight := logoBounds.Max.Y - logoBounds.Min.Y

	// Calculate the position to place the logo
	logoX := centerX - (logoWidth / 2)
	logoY := centerY - (logoHeight / 2)

	// Create a new image with the same dimensions as the QR code
	merged := image.NewRGBA(qrBounds)

	// Draw the QR code onto the new image
	draw.Draw(merged, qrBounds, *qrImg, image.Point{0, 0}, draw.Src)

	// Draw the logo onto the new image at the calculated position
	draw.Draw(merged, image.Rect(logoX, logoY, logoX+logoWidth, logoY+logoHeight), logo, image.Point{0, 0}, draw.Over)

	// Update the QR code image with the merged image
	*qrImg = merged
}
