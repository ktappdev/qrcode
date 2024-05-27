package routehandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ktappdev/qrcode-server/mongodb"
)

func HandleBackhalfCheck(c *gin.Context) {
	mongodb.CheckBackhalfAvailability(c)
}
