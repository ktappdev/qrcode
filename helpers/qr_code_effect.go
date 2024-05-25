package helpers

import (
	"image"
	"image/color"
	"image/draw"
	"math/rand"
)

func ApplyQRCodeEffect(img image.Image, blockSize int, probability float64) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	draw.Draw(result, bounds, img, bounds.Min, draw.Src)

	for y := bounds.Min.Y; y < bounds.Max.Y; y += blockSize {
		for x := bounds.Min.X; x < bounds.Max.X; x += blockSize {
			if rand.Float64() < probability {
				fillBlock(result, x, y, blockSize, color.Black)
			}
		}
	}

	return result
}

func fillBlock(img *image.RGBA, x, y, size int, c color.Color) {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			img.Set(x+i, y+j, c)
		}
	}
}
