package routehandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ktappdev/qrcode-server/urlhandler"
)

var linkExchanger = urlhandler.NewLinkExchanger()

func GetShortLink(c *gin.Context) {
	// Get the original link from the form data
	originalLink := c.PostForm("originalLink")
	if originalLink == "https://" {
		originalLink = ""
	}

	name := c.PostForm("name")
	owner := c.PostForm("owner")

	// Generate the short link
	shortLink := linkExchanger.GenerateShortLink(originalLink, name, owner)

	// Return the short link as the response
	c.JSON(http.StatusOK, gin.H{
		"short_link": shortLink,
	})
}
