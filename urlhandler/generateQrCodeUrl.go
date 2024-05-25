package urlhandler

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
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
