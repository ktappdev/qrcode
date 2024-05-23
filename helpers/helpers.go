package helpers

import (
	"encoding/hex"
	"fmt"
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
	logo = nil
	newLogo = nil
}

var ColorMap = map[string]string{
	"red":       "#FF0000",
	"green":     "#00FF00",
	"blue":      "#0000FF",
	"yellow":    "#FFFF00",
	"purple":    "#800080",
	"orange":    "#FFA500",
	"pink":      "#FFC0CB",
	"brown":     "#A52A2A",
	"gray":      "#808080",
	"black":     "#000000",
	"white":     "#FFFFFF",
	"turquoise": "#40E0D0",
	"indigo":    "#4B0082",
	"maroon":    "#800000",
	"lime":      "#00FF00",
	"teal":      "#008080",
}

func SetColours(backgroundColour, qrCodeColour string) (bgHex, qrHex string) {
	bgHex, ok := ColorMap[backgroundColour]
	if !ok {
		bgHex = ColorMap["white"]
	}

	qrHex, ok = ColorMap[qrCodeColour]
	if !ok {
		qrHex = ColorMap["black"]
	}

	return bgHex, qrHex
}

func HexToColor(hexString string) (color.Color, error) {
	hexBytes, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}

	switch len(hexBytes) {
	case 3:
		return color.RGBA{
			R: hexBytes[0],
			G: hexBytes[1],
			B: hexBytes[2],
			A: 0xff,
		}, nil
	case 4:
		return color.RGBA{
			R: hexBytes[0],
			G: hexBytes[1],
			B: hexBytes[2],
			A: hexBytes[3],
		}, nil
	default:
		return nil, fmt.Errorf("invalid hex color string: %s", hexString)
	}
}
