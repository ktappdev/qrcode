package urlhandler

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/ktappdev/qrcode-server/helpers"
	"github.com/ktappdev/qrcode-server/mongodb"
)

type URLExchanger struct {
	mu            sync.RWMutex
	qrCodeURLsMap map[string]string
}

func NewURLExchanger() *URLExchanger {
	return &URLExchanger{
		qrCodeURLsMap: make(map[string]string),
	}
}

func (e *URLExchanger) GenerateQRCodeURL(formData *helpers.FormDataStruct) string {
	log.Println("GenerateQRCodeURL ******", formData)
	port, server := GetEnvItems()
	uniqueID := uuid.New().String()

	//NOTE: Store the mapping in the Map (Keeping this for speed)
	e.mu.Lock()
	e.qrCodeURLsMap[uniqueID] = formData.OriginalLink
	e.mu.Unlock()

	//NOTE: Store the mapping in the database
	err := mongodb.InsertQRCodeURL(
		uniqueID,
		formData.OriginalLink,
		formData.BackgroundColour,
		formData.QRCodeColour,
		formData.Name,
	)
	if err != nil {
		log.Println("Error inserting URL into database")
		log.Fatal(err)
	}

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
