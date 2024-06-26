// mongodb/mongodb.go
package mongodb

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ktappdev/qrcode-server/helpers"
	"github.com/oschwald/geoip2-golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertQRCodeURL inserts a new QR code URL mapping into the database
func InsertQRCodeURL(id, originalLink string, backgroundColour, qrCodeColour string, name string) error {
	foregroundHex, backgroundHex := helpers.SetColours(backgroundColour, qrCodeColour)
	timestamp := time.Now().Format(time.RFC3339)
	qrCodeURL := QRCodeURL{
		ID:            id,
		OriginalURL:   originalLinkEmpty(originalLink, "https://592code.vercel.app/empty"),
		ForegroundHex: foregroundHex,
		BackgroundHex: backgroundHex,
		Name:          name,
		Type:          "qrcode",
		CreatedAt:     timestamp,
	}

	// Get a handle to the "qr_code_urls" collection in the database
	collection := client.Database("qr").Collection("qr_codes")
	filter := bson.M{"_id": id}
	update := bson.M{"$set": qrCodeURL}
	upsert := true // Create a boolean variable
	opts := options.UpdateOptions{
		Upsert: &upsert, // Pass the pointer to the boolean variable
	}

	_, err := collection.UpdateOne(context.Background(), filter, update, &opts)
	if err != nil {
		return err
	}

	log.Printf("Inserted or updated QR Code URL: %+v", qrCodeURL)
	return nil
}

// GetQRCodeURL retrieves the original URL for a given QR code ID
func GetQRCodeURL(id string) (string, error) {
	collection := client.Database("qr").Collection("qr_codes")

	// Create a variable to hold the retrieved document
	var qrCodeURL QRCodeURL

	// Find the document with the given ID and decode it into the qrCodeURL variable
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&qrCodeURL)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil // Return an empty string if the document is not found
		}
		return "", err
	}

	return qrCodeURL.OriginalURL, nil
}

func LogQRCodeInteraction(qrCodeID string, c *gin.Context, locationData *geoip2.City) error {
	timestamp := time.Now().Format(time.RFC3339)
	userAgent := c.Request.UserAgent()
	ipAddress := c.ClientIP()
	referer := c.Request.Referer()

	var regions []string
	for _, subdivision := range locationData.Subdivisions {
		if name, ok := subdivision.Names["en"]; ok {
			regions = append(regions, name)
		}
	}

	location := Location{
		Latitude:       locationData.Location.Latitude,
		Longitude:      locationData.Location.Longitude,
		TimeZone:       locationData.Location.TimeZone,
		MetroCode:      locationData.Location.MetroCode,
		AccuracyRadius: locationData.Location.AccuracyRadius,
		City:           locationData.City.Names["en"],
		PostalCode:     locationData.Postal.Code,
		Continent:      locationData.Continent.Names["en"],
		CountryName:    locationData.Country.Names["en"],
		CountryIsoCode: locationData.Country.IsoCode,
		Regions:        regions,
	}

	interaction := QRCodeInteraction{
		ID:        uuid.New().String(),
		QRCodeID:  qrCodeID,
		Timestamp: timestamp,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		Referer:   referer,
		Location:  location,
	}

	collection := client.Database("qr").Collection("qr_code_details")
	_, err := collection.InsertOne(context.Background(), interaction)
	if err != nil {
		return err
	}

	log.Printf("Logged QR Code Interaction: %+v", interaction)
	return nil
}
