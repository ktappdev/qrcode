package urlhandler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ktappdev/qrcode-server/geoip"
	"github.com/ktappdev/qrcode-server/mongodb"
)

func (e *LinkExchanger) HandleShortLinkInteraction(c *gin.Context, path string) {
	log.Println("this is the path param for the link", path)
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

	// uniqueID := c.Query("id")

	e.mu.RLock()
	originalURL, exists := e.linksMap[path]
	e.mu.RUnlock()

	if exists {
		// Mapping found in the in-memory map
		go mongodb.LogShortLinkInteraction(path, c, locationData)
		log.Println("this is the link returned from in-memory", originalURL)
		c.Redirect(http.StatusFound, originalURL)
		return
	}

	// Mapping not found in the in-memory map, check the database
	originalURL, err = mongodb.GetShortLink(path)
	if err == nil {
		// Mapping found in the database
		go mongodb.LogShortLinkInteraction(path, c, locationData)
		log.Println("this is the link returned from mongodb", originalURL)
		c.Redirect(http.StatusFound, originalURL)
		return
	}

	// Mapping not found in the in-memory map or the database
	fmt.Printf("%s NOT FOUND IN MEMORY OR DATABASE", path)
	c.Status(http.StatusNotFound)
}
