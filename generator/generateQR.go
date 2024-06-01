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
	qrCodeUrl string,
	size int,
	formData *helpers.FormDataStruct,
	logo *image.Image,
	overlayShrink int,
) ([]byte, error) {
	var qr *qrcode.QRCode
	var qrBytes []byte
	var qrImg image.Image
	if logo != nil {
		qr, _ = qrcode.New(qrCodeUrl, qrcode.High) // NOTE: can also set Highest
		log.Println("Logo is present in the generate function")
	} else {
		qr, _ = qrcode.New(qrCodeUrl, qrcode.Low)
	}

	foregroundColor, err := helpers.HexToColor(formData.QRCodeColour)
	if err != nil {
		// Handle the error
		return nil, err
	}
	backgroundColor, err := helpers.HexToColor(formData.BackgroundColour)
	if err != nil {
		// Handle the error
		return nil, err
	}

	qr.ForegroundColor = foregroundColor
	qr.BackgroundColor = backgroundColor

	if formData.UseDots == "true" {
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
		helpers.OverlayLogo(&qrImg, *logo, formData.Opacityf64, 3)
	}

	// Overlay your logo if overLayOurLogo is true
	if formData.OverlayOurLogo == "true" {
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
