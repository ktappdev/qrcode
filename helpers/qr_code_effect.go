package helpers

import (
	"image"
	"image/color"
	"math/rand"
)

func ApplyQRCodeEffect(img image.Image, blockSize int, probability float64) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			gray := uint8(0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8))
			sepia := color.RGBA{
				uint8(float64(gray)*0.9 + 20),
				uint8(float64(gray)*0.7 + 20),
				uint8(float64(gray)*0.4 + 20),
				255,
			}
			result.Set(x, y, sepia)
		}
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y += blockSize {
		for x := bounds.Min.X; x < bounds.Max.X; x += blockSize {
			if rand.Float64() < probability {
				fillBlock(result, x, y, blockSize, color.Black)
			} else if rand.Float64() < 0.1 {
				fillBlock(result, x, y, blockSize, color.White)
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
