package main

import (
	"github.com/gin-gonic/gin"
	"image"
	"net/http"
)

func LoadLogo(c *gin.Context) (*image.Image, error) {
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

	return &decodedLogo, nil
}
