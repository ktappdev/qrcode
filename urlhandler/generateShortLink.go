package urlhandler

import (
	"fmt"
	"log"
	"sync"

	"github.com/ktappdev/qrcode-server/helpers"
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

func (e *LinkExchanger) GenerateShortLink(originalURL, backhalf, name, owner string) (string, error) {
	port, server := GetEnvItems()
	var uniqueID string
	if backhalf != "" {
		uniqueID = backhalf
	} else {
		uniqueID = helpers.GenerateUniqueString(4)
	}

	// Store the mapping in the Map (Keeping this for speed)
	e.mu.Lock()
	e.linksMap[uniqueID] = originalURL
	e.mu.Unlock()

	// Store the mapping in the database
	err := mongodb.InsertShortLink(uniqueID, originalURL, name, owner)
	if err != nil {
		// log.Println("Error inserting URL into database")
		return "Error inserting URL into database", err

	}

	log.Println("list of Short Links:", e.linksMap)

	var link string
	if server != "https://gr.lugetech.com" {
		fmt.Println("Using local server with port, if this is running on the remote server it will not work")
		link = fmt.Sprintf("%s:%s/%s", server, port, uniqueID)
	} else {
		link = fmt.Sprintf("%s/%s", server, uniqueID)
	}

	fmt.Println("link", link)
	return link, nil
}
