package urlexchanger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ktappdev/qrcode-server/geoip"
	"github.com/ktappdev/qrcode-server/mongodb"
)

func getEnvItems() (port string, server string) {
	port = os.Getenv("PORT")
	server = os.Getenv("SERVER")
	return port, server
}

type URLExchanger struct {
	mu            sync.RWMutex
	qrCodeURLsMap map[string]string
}

func NewURLExchanger() *URLExchanger {
	return &URLExchanger{
		qrCodeURLsMap: make(map[string]string),
	}
}

func (e *URLExchanger) GenerateQRCodeURL(originalLink string, backgroundColour, qrCodeColour string, name string) string {
	port, server := getEnvItems()
	uniqueID := uuid.New().String()

	//NOTE: Store the mapping in the Map (Keeping this for speed)
	e.mu.Lock()
	e.qrCodeURLsMap[uniqueID] = originalLink
	e.mu.Unlock()

	//NOTE: Store the mapping in the database
	err := mongodb.InsertQRCodeURL(uniqueID, originalLink, backgroundColour, qrCodeColour, name)
	if err != nil {
		log.Println("Error inserting URL into database")
		log.Fatal(err)
	}

	log.Println("list of QR Codes:", e.qrCodeURLsMap)
	var link string
	if server != "https://qr.lugetech.com" {
		fmt.Println("Using local server with port, if this is running on the remote server it will not work")
		link = fmt.Sprintf("%s:%s/qr?id=%s", server, port, uniqueID)
	} else {
		link = fmt.Sprintf("%s/qr?id=%s", server, uniqueID)
	}
	fmt.Println("link", link)

	return link
}

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
	city, err := geoIP.LookupCity(clientIP)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("THIS IS THE CITY", city)
	uniqueID := c.Query("id")
	e.mu.RLock()
	originalURL, exists := e.qrCodeURLsMap[uniqueID]
	e.mu.RUnlock()
	if exists {
		//NOTE: Mapping found in the in-memory map
		mongodb.LogQRCodeInteraction(uniqueID, c.Request)
		log.Println("this is the link returned from in-memory", originalURL)
		c.Redirect(http.StatusFound, originalURL)
		return
	}

	//NOTE: Mapping not found in the in-memory map, check the database
	originalURL, err = mongodb.GetQRCodeURL(uniqueID)
	if err == nil {
		// Mapping found in the database
		mongodb.LogQRCodeInteraction(uniqueID, c)
		log.Println("this is the link returned from mongodb", originalURL)
		c.Redirect(http.StatusFound, originalURL)
		return
	}

	//NOTE: Mapping not found in the in-memory map or the database
	fmt.Printf("%s NOT FOUND IN MEMORY OR DATABASE", uniqueID)
	c.Status(http.StatusNotFound)
}
