package mongodb

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RequestData struct {
	Backhalf string `json:"backhalf"`
	UserID   int    `json:"userId"`
}

// handlers.go
func CheckBackhalfAvailability(c *gin.Context) {
	var requestData RequestData

	// Bind the JSON request body to the requestData struct
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to bind request body",
		})
		return
	}

	// Access the bound data from the requestData struct
	backhalf := requestData.Backhalf
	_ = requestData.UserID

	// Get a handle to the "short_links" collection
	collection := client.Database("links").Collection("short_links")
	// Check if the backhalf exists in the collection
	filter := bson.M{"_id": backhalf}
	var result bson.M
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusOK, gin.H{
				"available": true,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check backhalf availability",
		})
		return
	}
	c.JSON(http.StatusConflict, gin.H{
		"available": false,
	})
}
