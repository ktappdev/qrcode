package helpers

import (
	"image"
	"image/color"
	"strconv"

	"golang.org/x/image/draw"
)

// ParseOpacity parses a string to a float64 value.
// It returns the parsed float64 value and any error that occurred during parsing.
func ParseOpacity(opacityStr string) (float64, error) {
	wholeNumber, err := strconv.ParseFloat(opacityStr, 64)
	return wholeNumber / 100, err
}

func OverlayLogo(qrImg *image.Image, logo image.Image, opacity float64) {
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

	// Calculate the desired logo size (1/10th of QR code size)
	desiredLogoWidth := qrWidth / 3
	desiredLogoHeight := qrHeight / 3

	// Create a new image with the adjusted logo size
	newLogoWidth, newLogoHeight := desiredLogoWidth, desiredLogoHeight
	if logoWidth < desiredLogoWidth && logoHeight < desiredLogoHeight {
		// If the original logo is smaller than the desired size, use the original size
		newLogoWidth, newLogoHeight = logoWidth, logoHeight
	}
	newLogo := image.NewRGBA(image.Rect(0, 0, newLogoWidth, newLogoHeight))
	draw.BiLinear.Scale(newLogo, newLogo.Rect, logo, logo.Bounds(), draw.Over, nil)

	// Calculate the position to place the resized logo
	logoX := centerX - (newLogoWidth / 2)
	logoY := centerY - (newLogoHeight / 2)

	// Create a new image with the same dimensions as the QR code
	merged := image.NewRGBA(qrBounds)

	// Draw the QR code onto the new image
	draw.Draw(merged, qrBounds, *qrImg, image.ZP, draw.Src)

	// Create a new image with the same dimensions as the resized logo
	logoMask := image.NewRGBA(image.Rect(0, 0, newLogoWidth, newLogoHeight))

	// Fill the logo mask with the specified opacity
	opacity255 := uint8(opacity * 255)
	for y := 0; y < newLogoHeight; y++ {
		for x := 0; x < newLogoWidth; x++ {
			logoMask.Set(x, y, color.RGBA{0, 0, 0, opacity255})
		}
	}

	// Draw the resized logo onto the merged image using the logo mask
	draw.DrawMask(merged, image.Rect(logoX, logoY, logoX+newLogoWidth, logoY+newLogoHeight),
		newLogo, image.ZP, logoMask, image.ZP, draw.Over)

	// Update the QR code image with the merged image
	*qrImg = merged
}
