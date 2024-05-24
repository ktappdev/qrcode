package urlexchanger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ktappdev/qrcode-server/mongodb"
)

func getEnvItems() (port string, server string) {
	port = os.Getenv("PORT")
	server = os.Getenv("SERVER")
	return port, server
}

type URLExchanger struct {
	mu         sync.RWMutex
	qrCodeURLs map[string]string
}

func NewURLExchanger() *URLExchanger {
	return &URLExchanger{
		qrCodeURLs: make(map[string]string),
	}
}

func (e *URLExchanger) GenerateQRCodeURL(originalURL string, backgroundColour, qrCodeColour string, name string) string {
	port, server := getEnvItems()
	uniqueID := uuid.New().String()

	//NOTE: Store the mapping in the Map (Keeping this for speed)
	e.mu.Lock()
	e.qrCodeURLs[uniqueID] = originalURL
	e.mu.Unlock()

	//NOTE: Store the mapping in the database
	err := mongodb.InsertQRCodeURL(uniqueID, originalURL, backgroundColour, qrCodeColour, name)
	if err != nil {
		log.Println("Error inserting URL into database")
		log.Fatal(err)
	}

	fmt.Println("list of QR Codes:", e.qrCodeURLs)
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
	uniqueID := c.Query("id")
	start := time.Now() // Record the start time
	e.mu.RLock()
	originalURL, exists := e.qrCodeURLs[uniqueID]
	e.mu.RUnlock()
	elapsed := time.Since(start) // Calculate the elapsed time
	if exists {
		//NOTE: Mapping found in the in-memory map
		mongodb.LogQRCodeInteraction(uniqueID, c.Request)
		log.Printf("Link returned from in-memory: %s (Took %v)", originalURL, elapsed)
		log.Println("this is the link returned from in-memory", originalURL)
		c.Redirect(http.StatusFound, originalURL)
		return
	}

	//NOTE: Mapping not found in the in-memory map, check the database
	start = time.Now() // Record the start time for MongoDB operation
	originalURL, err := mongodb.GetQRCodeURL(uniqueID)
	elapsed = time.Since(start) // Calculate the elapsed time
	if err == nil {
		// Mapping found in the database
		mongodb.LogQRCodeInteraction(uniqueID, c.Request)
		log.Printf("Link returned from MongoDB: %s (Took %v)", originalURL, elapsed)
		log.Println("this is the link returned from mongodb", originalURL)
		c.Redirect(http.StatusFound, originalURL)
		return
	}

	//NOTE: Mapping not found in the in-memory map or the database
	fmt.Printf("%s NOT FOUND IN MEMORY OR DATABASE", uniqueID)
	c.Status(http.StatusNotFound)
}

// func (e *URLExchanger) logInteraction(uniqueID string, r *http.Request) {
// 	timestamp := time.Now().Format(time.RFC3339)
// 	userAgent := r.UserAgent()
// 	ipAddress := r.RemoteAddr
// 	referer := r.Referer()
// 	fmt.Printf("QR Code Interaction - ID: %s, Timestamp: %s, User Agent: %s, IP Address: %s, referrer: %s", uniqueID, timestamp, userAgent, ipAddress, referer)
//
// }
