package qrreader

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/tuotoo/qrcode"
)

func ReadQRCodeFromFile(filename string) (string, error) {
	// Open the QR code image file
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("error opening image file: %v", err)
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("error decoding image: %v", err)
	}

	// Convert image to byte slice
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return "", fmt.Errorf("error encoding image: %v", err)
	}

	// Create io.Reader from byte slice
	reader := bytes.NewReader(buf.Bytes())

	// Create a new QRCode instance
	qrMatrix, err := qrcode.Decode(reader)
	if err != nil {
		return "", fmt.Errorf("error decoding QR code: %v", err)
	}

	return qrMatrix.Content, nil
}
