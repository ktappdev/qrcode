package routehandlers

import (
	"github.com/gin-gonic/gin"
	"log"
)

func HandleScan(c *gin.Context) {
	log.Println("HandleScan")
	exchanger.HandleQRCodeInteraction(c)
}
