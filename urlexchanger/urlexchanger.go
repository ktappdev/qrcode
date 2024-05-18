package urlexchanger

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (e *URLExchanger) GenerateQRCodeURL(originalURL string) string {
	port, server := getEnvItems()
	uniqueID := uuid.New().String()

	e.mu.Lock()
	e.qrCodeURLs[uniqueID] = originalURL
	e.mu.Unlock()
	fmt.Println("list of QR Codes:", e.qrCodeURLs)
	if server != "https://qr.lugetech.com" {
		fmt.Println("Using local server with port, if this is running on the remote server it will not work")
		link := fmt.Sprintf("%s:%s/qr?id=%s", server, port, uniqueID)
		fmt.Println("link", link)
		return link
	}
	link := fmt.Sprintf("%s/qr?id=%s", server, uniqueID)
	fmt.Println("link", link)

	return link
}

func (e *URLExchanger) HandleQRCodeInteraction(c *gin.Context) {
	uniqueID := c.Query("id")

	e.mu.RLock()
	originalURL, exists := e.qrCodeURLs[uniqueID]
	e.mu.RUnlock()

	if !exists {
		fmt.Printf("Not found %s", uniqueID)
		c.Status(http.StatusNotFound)
		return
	}

	e.logInteraction(uniqueID, c.Request)

	c.Redirect(http.StatusFound, originalURL)
}

func (e *URLExchanger) logInteraction(uniqueID string, r *http.Request) {
	timestamp := time.Now().Format(time.RFC3339)
	userAgent := r.UserAgent()
	ipAddress := r.RemoteAddr
	referer := r.Referer()
	// ... Log or store the interaction details ...
	fmt.Printf("QR Code Interaction - ID: %s, Timestamp: %s, User Agent: %s, IP Address: %s, referrer: %s", uniqueID, timestamp, userAgent, ipAddress, referer)

}
