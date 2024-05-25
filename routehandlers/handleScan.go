package routehandlers

import (
	"github.com/gin-gonic/gin"
)

func HandleScan(c *gin.Context) {
	exchanger.HandleQRCodeInteraction(c)
}
