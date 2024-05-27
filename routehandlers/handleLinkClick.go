package routehandlers

import (
	"github.com/gin-gonic/gin"
)

func HandleLinkClick(c *gin.Context, path string) {
	linkExchanger.HandleShortLinkInteraction(c, path)
}
