package mongodb

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetInteractionsForQRCode(c *gin.Context) {
	qrCodeID := c.Query("id")

	// Get a handle to the "qr_code_details" collection
	collection := client.Database("qr").Collection("qr_code_details")

	// Create a filter to match documents with the given qrCodeID
	filter := bson.M{"qr_code_id": qrCodeID}

	// Initialize an empty slice to hold the retrieved interaction documents
	var interactions []QRCodeInteraction

	// Use the Find method to get a cursor over the matching documents
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve QR code details",
		})
		return
	}

	// Ensure the cursor is closed at the end of the function
	defer cursor.Close(context.Background())

	// Iterate over the cursor and decode the documents into the interactions slice
	if err = cursor.All(context.Background(), &interactions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to decode QR code details",
		})
		return
	}

	// Return the slice of interaction documents as JSON with a status
	c.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"interactions": interactions,
	})
}

func GetMostRecentInteractionForQRCode(qrCodeID string) (*QRCodeInteraction, error) {
	// Get a handle to the "qr_code_details" collection
	collection := client.Database("qr").Collection("qr_code_details")

	// Create a filter to match documents
	filter := bson.M{"qr_code_id": qrCodeID}

	// Setgh the FindOne options to sort the results by timestamp in descending order
	opts := options.FindOne().SetSort(bson.D{
		{Key: "timestamp", Value: -1}, // Specify the key-value pair for sorting
	})

	// Initialize a variable to hold the retrieved interaction document
	var interaction QRCodeInteraction

	// Use the FindOne method with the filter and sort options to get the most recent interaction
	err := collection.FindOne(context.Background(), filter, opts).Decode(&interaction)
	if err != nil {
		// If no interaction is found, return nil
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	// Return a pointer to the interaction document
	return &interaction, nil
}

func GetInteractionsForQRCodeInTimeRange(qrCodeID string, startTime, endTime time.Time) ([]QRCodeInteraction, error) {
	// Get a handle to the "qr_code_details" collection
	collection := client.Database("qr").Collection("qr_code_details")

	// Create a filter to match documents with the given qrCodeID and timestamp within the specified range
	filter := bson.M{
		"qr_code_id": qrCodeID,
		"timestamp": bson.M{
			"$gte": startTime.Format(time.RFC3339), // Greater than or equal to the start time
			"$lt":  endTime.Format(time.RFC3339),   // Less than the end time
		},
	}

	// Initialize an empty slice to hold the retrieved interaction documents
	var interactions []QRCodeInteraction

	// Use the Find method with the filter to get a cursor over the matching documents
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Iterate over the cursor and decode the documents into the interactions slice
	if err = cursor.All(context.Background(), &interactions); err != nil {
		return nil, err
	}

	// Return the slice of interaction documents
	return interactions, nil
}

func GetTotalInteractionsForQRCode(qrCodeID string) (int64, error) {
	// Get a handle to the "qr_code_details" collection
	collection := client.Database("qr").Collection("qr_code_details")

	// Create a filter to match documents with the given qrCodeID
	filter := bson.M{"qr_code_id": qrCodeID}

	// Use the Count Documents method with the filter to get the total count of matching documents
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return 0, err
	}

	// Return the total count of interactions
	return count, nil
}

func GetInteractionsByUserAgentForQRCode(qrCodeID string) ([]bson.M, error) {
	// Get a handle to the "qr_code_details" collection
	collection := client.Database("qr").Collection("qr_code_details")

	// Create a filter to match documents with the given qrCodeID
	filter := bson.M{"qr_code_id": qrCodeID}

	// Define the aggregation pipeline stages
	pipeline := []bson.M{
		// Match documents with the given qrCodeID
		{"$match": filter},
		// Group documents by user_agent and calculate the count, unique IP addresses, and unique referrers
		{"$group": bson.M{
			"_id":       "$user_agent",                      // Group by user_agent
			"count":     bson.M{"$sum": 1},                  // Count of interactions for each user_agent
			"ipAddrs":   bson.M{"$addToSet": "$ip_address"}, // Unique IP addresses for each user_agent
			"referrers": bson.M{"$addToSet": "$referer"},    // Unique referrers for each user_agent
		}},
	}

	// Use the Aggregate method with the pipeline to execute the aggregation
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Initialize a slice to hold the aggregation results
	var results []bson.M

	// Iterate over the cursor and decode the results into the results slice
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	// Return the slice of aggregation results
	return results, nil
}
