package routehandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ktappdev/qrcode-server/urlhandler"
)

var linkExchanger = urlhandler.NewLinkExchanger()

func GetShortLink(c *gin.Context) {
	// Get the original link from the form data
	originalURL := c.PostForm("originalURL")
	if originalURL == "https://" {
		originalURL = ""
	}

	backhalf := c.PostForm("backhalf")
	name := c.PostForm("name")
	owner := c.PostForm("owner")

	shortLink, err := linkExchanger.GenerateShortLink(originalURL, backhalf, name, owner)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_link": shortLink,
	})
}
