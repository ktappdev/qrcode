package qrcode

import (
	"bytes"
	"image/png"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(data string, size int) ([]byte, error) {
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return nil, err
	}

	img := qr.Image(size)

	buf := &bytes.Buffer{}
	err = png.Encode(buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
