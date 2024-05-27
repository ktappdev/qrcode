package urlhandler

import (
	"fmt"
	"log"
	"math/rand"
	"sync"

	"github.com/ktappdev/qrcode-server/mongodb"
)

type LinkExchanger struct {
	mu       sync.RWMutex
	linksMap map[string]string
}

func NewLinkExchanger() *LinkExchanger {
	return &LinkExchanger{
		linksMap: make(map[string]string),
	}
}

func (e *LinkExchanger) GenerateShortLink(originalURL string, name, owner string) string {
	port, server := GetEnvItems()
	uniqueID := generateUniqueString(4)

	// Store the mapping in the Map (Keeping this for speed)
	e.mu.Lock()
	e.linksMap[uniqueID] = originalURL
	e.mu.Unlock()

	// Store the mapping in the database
	err := mongodb.InsertShortLink(uniqueID, originalURL, name, owner)
	if err != nil {
		log.Println("Error inserting URL into database")
		log.Fatal(err)
	}

	log.Println("list of Short Links:", e.linksMap)

	var link string
	if server != "https://short.lugetech.com" {
		fmt.Println("Using local server with port, if this is running on the remote server it will not work")
		link = fmt.Sprintf("%s:%s/%s", server, port, uniqueID)
	} else {
		link = fmt.Sprintf("%s/link?id=%s", server, uniqueID)
	}

	fmt.Println("link", link)
	return link
}

func generateUniqueString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}