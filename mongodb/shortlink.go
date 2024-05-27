// mongodb/mongodb.go
package mongodb

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oschwald/geoip2-golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertShortLink inserts a new short link mapping into the database
func InsertShortLink(uniqueId, originalURL, name, owner string) error {
	shortLink := ShortLink{
		Type:        "shortlink",
		CreatedAt:   time.Now().Format(time.RFC3339),
		ID:          uniqueId,
		OriginalURL: originalLinkEmpty(originalURL, "https://592code.vercel.app/empty"),
		Name:        name,
		Owner:       owner,
	}

	// Get a handle to the "short_links" collection in the database
	collection := client.Database("links").Collection("short_links")

	filter := bson.M{"_id": shortLink.ID}
	update := bson.M{"$set": shortLink}
	upsert := true // Create a boolean variable

	opts := options.UpdateOptions{
		Upsert: &upsert, // Pass the pointer to the boolean variable
	}

	_, err := collection.UpdateOne(context.Background(), filter, update, &opts)
	if err != nil {
		return err
	}

	log.Printf("Inserted or updated Short Link: %+v", shortLink)
	return nil
}

// GetShortLink retrieves the original URL for a given short link ID
func GetShortLink(id string) (string, error) {
	collection := client.Database("links").Collection("short_links")

	// Create a variable to hold the retrieved document
	var shortLink ShortLink

	// Find the document with the given ID and decode it into the shortLink variable
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&shortLink)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil // Return an empty string if the document is not found
		}
		return "", err
	}

	return shortLink.OriginalURL, nil
}

func LogShortLinkInteraction(shortLinkID string, c *gin.Context, locationData *geoip2.City) error {
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

	interaction := ShortLinkInteraction{
		ID:          uuid.New().String(),
		ShortLinkID: shortLinkID,
		Timestamp:   timestamp,
		UserAgent:   userAgent,
		IPAddress:   ipAddress,
		Referer:     referer,
		Location:    location,
	}

	collection := client.Database("links").Collection("short_link_details")

	_, err := collection.InsertOne(context.Background(), interaction)
	if err != nil {
		return err
	}

	log.Printf("Logged Short Link Interaction: %+v", interaction)
	return nil
}
