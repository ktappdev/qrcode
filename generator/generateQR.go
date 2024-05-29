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
	useDots bool,
) ([]byte, error) {
	log.Println("this is generate and the use dots is", useDots)
	var qr *qrcode.QRCode
	var qrBytes []byte
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

	if useDots {
		log.Println("Dots are being used")
		qrBytes, err = drawQRCodeWithDots(qr, size, foregroundColor, backgroundColor)
		if err != nil {
			return nil, err
		}
	} else {
		// Generate the QR code image
		qrImg := qr.Image(size)

		// Overlay the logo image if provided
		if logo != nil {
			log.Println("Logo is present and about to overlay")
			helpers.OverlayLogo(&qrImg, *logo, opacity)
		}

		log.Println("this is after overlay and about to encode png")
		buf := &bytes.Buffer{}
		err = png.Encode(buf, qrImg)
		if err != nil {
			return nil, err
		}

		qrBytes = buf.Bytes()
	}

	return qrBytes, nil
}
