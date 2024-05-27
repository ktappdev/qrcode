package routehandlers

import (
	"github.com/gin-gonic/gin"
)

func HandleLinkClick(c *gin.Context) {
	linkExchanger.HandleShortLinkInteraction(c)
}
