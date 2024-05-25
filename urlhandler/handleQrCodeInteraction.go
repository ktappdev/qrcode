package urlhandler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ktappdev/qrcode-server/geoip"
	"github.com/ktappdev/qrcode-server/mongodb"
)

func (e *URLExchanger) HandleQRCodeInteraction(c *gin.Context) {
	clientIP := c.ClientIP()
	geoIP, err := geoip.New("./geo-lite/GeoLite2-City.mmdb")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer geoIP.Close()
	// The `geoIP` instance now contains the initialized `geoip2.Reader`
	// You can use it to perform lookups
	locationData, err := geoIP.LookupCity(clientIP)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("THIS IS THE CITY", locationData)
	uniqueID := c.Query("id")
	e.mu.RLock()
	originalURL, exists := e.qrCodeURLsMap[uniqueID]
	e.mu.RUnlock()
	if exists {
		//NOTE: Mapping found in the in-memory map
		go mongodb.LogQRCodeInteraction(uniqueID, c, locationData)
		log.Println("this is the link returned from in-memory", originalURL)
		c.Redirect(http.StatusFound, originalURL)
		return
	}

	//NOTE: Mapping not found in the in-memory map, check the database
	originalURL, err = mongodb.GetQRCodeURL(uniqueID)
	if err == nil {
		// Mapping found in the database
		go mongodb.LogQRCodeInteraction(uniqueID, c, locationData)
		log.Println("this is the link returned from mongodb", originalURL)
		c.Redirect(http.StatusFound, originalURL)
		return
	}

	//NOTE: Mapping not found in the in-memory map or the database
	fmt.Printf("%s NOT FOUND IN MEMORY OR DATABASE", uniqueID)
	c.Status(http.StatusNotFound)
}
