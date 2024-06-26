package helpers

import (
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	// qrcodeeffect "github.com/ktappdev/qrcode-server/helpers/qr_code_effect"
	"golang.org/x/image/draw"
)

// ParseOpacity parses a string to a float64 value.
// It returns the parsed float64 value and any error that occurred during parsing.
func ParseOpacity(opacityStr string) (float64, error) {
	wholeNumber, err := strconv.ParseFloat(opacityStr, 64)
	return wholeNumber / 100, err
}

func OverlayLogo(qrImg *image.Image, logo image.Image, opacity float64, overlayShrink int) {
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
	desiredLogoWidth := qrWidth / overlayShrink
	desiredLogoHeight := qrHeight / overlayShrink

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

func SetColours(backgroundColour, qrCodeColour string) (bgHex, qrHex string) {
	bgHex, err := IsValidHexColor(backgroundColour)
	if err != nil {
		log.Println("Error validating background colour")
		log.Fatal(err)
	}
	qrHex, err = IsValidHexColor(qrCodeColour)
	if err != nil {
		log.Println("Error validating qr code colour")
		log.Fatal(err)
	}
	return bgHex, qrHex
}

func IsValidHexColor(s string) (string, error) {
	log.Print(s)
	if !strings.HasPrefix(s, "#") {
		return "", fmt.Errorf("invalid hex color string: %s", s)
	}
	s = s[1:]

	if len(s) != 3 && len(s) != 6 {
		return "", fmt.Errorf("invalid hex color string length: %s", s)
	}

	for _, char := range s {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
			return "", fmt.Errorf("invalid hex color string: %s", s)
		}
	}

	return s, nil
}

func HexToColor(hexString string) (color.Color, error) {
	if !strings.HasPrefix(hexString, "#") {
		return nil, fmt.Errorf("invalid hex color string: %s", hexString)
	}
	hexString = hexString[1:] // Remove leading "#"

	var r, g, b, a uint8
	var err error

	switch len(hexString) {
	case 3:
		r, err = hexToByte(hexString[0:1] + hexString[0:1])
		if err != nil {
			return nil, err
		}
		g, err = hexToByte(hexString[1:2] + hexString[1:2])
		if err != nil {
			return nil, err
		}
		b, err = hexToByte(hexString[2:3] + hexString[2:3])
		if err != nil {
			return nil, err
		}
		a = 0xff
	case 4:
		r, err = hexToByte(hexString[0:1] + hexString[0:1])
		if err != nil {
			return nil, err
		}
		g, err = hexToByte(hexString[1:2] + hexString[1:2])
		if err != nil {
			return nil, err
		}
		b, err = hexToByte(hexString[2:3] + hexString[2:3])
		if err != nil {
			return nil, err
		}
		a, err = hexToByte(hexString[3:4] + hexString[3:4])
		if err != nil {
			return nil, err
		}
	case 6:
		r, err = hexToByte(hexString[0:2])
		if err != nil {
			return nil, err
		}
		g, err = hexToByte(hexString[2:4])
		if err != nil {
			return nil, err
		}
		b, err = hexToByte(hexString[4:6])
		if err != nil {
			return nil, err
		}
		a = 0xff
	case 8:
		r, err = hexToByte(hexString[0:2])
		if err != nil {
			return nil, err
		}
		g, err = hexToByte(hexString[2:4])
		if err != nil {
			return nil, err
		}
		b, err = hexToByte(hexString[4:6])
		if err != nil {
			return nil, err
		}
		a, err = hexToByte(hexString[6:8])
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid hex color string length: %s", hexString)
	}

	return color.RGBA{R: r, G: g, B: b, A: a}, nil
}

// hexToByte converts a hex string to a byte.
func hexToByte(hexStr string) (byte, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil || len(bytes) != 1 {
		return 0, fmt.Errorf("invalid hex byte: %s", hexStr)
	}
	return bytes[0], nil
}

func LoadLogo(c *gin.Context, effect bool) (*image.Image, error) {
	logoFile, err := c.FormFile("logo")
	if err != nil {
		if err != http.ErrMissingFile {
			return nil, err
		}
		// No logo file provided, return nil
		return nil, nil
	}

	file, err := logoFile.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decodedLogo, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	if effect {
		// Apply QR code effect to the decoded logo
		modifiedLogo := ApplyQRCodeEffect(decodedLogo, 10, 0.3)

		return &modifiedLogo, nil
	} else {
		return &decodedLogo, nil
	}
}
