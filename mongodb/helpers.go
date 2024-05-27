package mongodb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// client represents the MongoDB client instance
var client *mongo.Client

// Connect establishes a connection to the MongoDB Atlas cluster
func Connect(uri string) error {
	var err error
	// Create a new MongoDB client using the provided connection string
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	// Ping the MongoDB server to verify the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	// collection := client.Database("links").Collection("short_links")
	// // Create a unique index on the "ID" field
	// indexModel := mongo.IndexModel{
	// 	Keys:    bson.M{"email": 1}, // Index key
	// 	Options: options.Index().SetUnique(true),
	// }
	// _, err = collection.Indexes().CreateOne(context.Background(), indexModel)
	// if err != nil {
	// 	return fmt.Errorf("failed to create mongodb index for backhalf")
	// }

	log.Println("Connected to MongoDB Atlas")
	return nil
}

func originalLinkEmpty(originalLink string, defaultLink string) string {
	if originalLink != "" {
		return originalLink
	}
	return defaultLink
}
func IsDuplicateKeyError(err error) bool {
	// Check if the error is a WriteException and the error code is 11000 (duplicate key error)
	writeException, ok := err.(mongo.WriteException)
	if ok && len(writeException.WriteErrors) > 0 && writeException.WriteErrors[0].Code == 11000 {
		return true
	}
	return false
}

// QRCodeURL represents the data model for a QR code URL mapping
type QRCodeURL struct {
	Type          string `bson:"type"`
	ID            string `bson:"_id"`          // This field maps to the "_id" field in the MongoDB document
	OriginalURL   string `bson:"original_url"` // This field maps to the "original_url" field in the MongoDB document
	ForegroundHex string `bson:"foreground_hex"`
	BackgroundHex string `bson:"background_hex"`
	Name          string `bson:"name"`
	Owner         string `bson:"owner"`
	CreatedAt     string `bson:"timestamp"`
}
type Location struct {
	Latitude       float64
	Longitude      float64
	TimeZone       string
	MetroCode      uint
	AccuracyRadius uint16
	City           string
	PostalCode     string
	Continent      string
	CountryName    string
	CountryIsoCode string
	Regions        []string
}
type QRCodeInteraction struct {
	Type      string `bson:"type"`
	ID        string `bson:"_id"`
	QRCodeID  string `bson:"qr_code_id"`
	Timestamp string `bson:"timestamp"`
	UserAgent string `bson:"user_agent"`
	IPAddress string `bson:"ip_address"`
	Referer   string `bson:"referer"`
	Location  Location
}
type User struct {
	ID        string   `bson:"_id"`
	Name      string   `bson:"name"`
	Email     string   `bson:"email"`
	QRCodeIDs []string `bson:"qr_code_ids"`
}

type ShortLink struct {
	Type        string `bson:"type"`
	ID          string `bson:"_id"`
	OriginalURL string `bson:"original_url"`
	Name        string `bson:"name"`
	Owner       string `bson:"owner"`
	CreatedAt   string `bson:"timestamp"`
}

type ShortLinkInteraction struct {
	Type        string `bson:"type"`
	ID          string `bson:"_id"`
	ShortLinkID string `bson:"short_link_id"`
	Timestamp   string `bson:"timestamp"`
	UserAgent   string `bson:"user_agent"`
	IPAddress   string `bson:"ip_address"`
	Referer     string `bson:"referer"`
	Location    Location
}
