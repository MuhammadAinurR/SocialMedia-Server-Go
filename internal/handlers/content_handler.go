package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cms-server/internal/database"
	"cms-server/internal/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Helper to get collection
func getContentCollection() *mongo.Collection {
	return database.GetCollection("contents")
}

func getStackCollection() *mongo.Collection {
	return database.GetCollection("stacks")
}

// handleError sends an error response
func handleError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// getUserIDFromContext retrieves the user ID from the request context
func getUserIDFromContext(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value("userID").(string)
	return userID, ok
}

// createContent inserts the content into the database
func createContent(content models.Content) error {
	collection := database.GetCollection("contents")
	_, err := collection.InsertOne(context.TODO(), content)
	return err
}

// fetchContents retrieves content based on a filter
func fetchContents(filter interface{}) ([]models.Content, error) {
	collection := getContentCollection()

	// Set up a context with a timeout for querying MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var contents []models.Content

	// Retrieve documents based on the filter
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode documents into the contents slice
	for cursor.Next(ctx) {
		var content models.Content
		if err := cursor.Decode(&content); err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}

	// Check if there was an error during the cursor iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return contents, nil
}

// fetchStacks checks if all stack names exist and returns their details
func fetchStacks(stackNames []string) ([]models.Stack, error) {
	var stackDetails []models.Stack

	// Get the collection
	collection := getStackCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find the stack documents based on the provided names
	cursor, err := collection.Find(ctx, bson.M{"name": bson.M{"$in": stackNames}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var stack models.Stack
		if err := cursor.Decode(&stack); err != nil {
			return nil, err
		}
		stackDetails = append(stackDetails, stack)
	}

	// Check if we found all the stacks
	if len(stackDetails) != len(stackNames) {
		return nil, fmt.Errorf("one or more stacks not found")
	}

	return stackDetails, nil
}

// CreateContentHandler handles the creation of new content
func CreateContentHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		handleError(w, "Unable to retrieve user ID", http.StatusInternalServerError)
		return
	}

	var requestBody struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Url         string   `json:"url"`
		ImgUrl      string   `json:"imgUrl"`
		Stack       []string `json:"stack"` // Array of stack IDs
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		handleError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate and fetch the stack data
	stackDetails, err := fetchStacks(requestBody.Stack)
	if err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	content := models.Content{
		UserID:      userID,
		Name:        requestBody.Name,
		Description: requestBody.Description,
		Url:         requestBody.Url,
		ImgUrl:      requestBody.ImgUrl,
		Stack:       stackDetails, // Use the fetched stack details
	}

	if err := createContent(content); err != nil {
		handleError(w, "Error creating content", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Content created successfully"})
}

// GetContentHandler retrieves content for a specific user
func GetContentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := getUserIDFromContext(r)
	if !ok {
		handleError(w, "Unable to retrieve user ID", http.StatusInternalServerError)
		return
	}

	filter := bson.D{{Key: "user_id", Value: userID}}
	contents, err := fetchContents(filter)
	if err != nil {
		handleError(w, "Error fetching content", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(contents)
}

// GetContentsHandler retrieves all the content from the database
func GetContentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	contents, err := fetchContents(bson.M{})
	if err != nil {
		handleError(w, "Error fetching content", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(contents)
}

func EditContentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the user ID from the context
	userID, ok := getUserIDFromContext(r)
	if !ok {
		handleError(w, "Unable to retrieve user ID", http.StatusInternalServerError)
		return
	}

	// Extract the content ID from the request URL (assuming it's passed as a path parameter)
	contentID := mux.Vars(r)["id"]

	// Check if content exists and belongs to the user
	collection := database.GetCollection("contents")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert string contentID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(contentID)
	if err != nil {
		handleError(w, "Invalid content ID", http.StatusBadRequest)
		return
	}

	// Find the content by ID and user ID to ensure ownership
	filter := bson.M{"_id": objectID, "user_id": userID}
	var content models.Content
	err = collection.FindOne(ctx, filter).Decode(&content)
	if err != nil {
		handleError(w, "Content not found or unauthorized", http.StatusNotFound)
		return
	}

	// Decode the new content data from the request body
	var updatedContent models.Content
	err = json.NewDecoder(r.Body).Decode(&updatedContent)
	if err != nil {
		handleError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the content in the database
	update := bson.M{
		"$set": bson.M{
			"name":   updatedContent.Name,
			"url":    updatedContent.Url,
			"imgUrl": updatedContent.ImgUrl,
			"stack":  updatedContent.Stack,
		},
	}
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		handleError(w, "Error updating content", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Content updated successfully"))
}

func DeleteContentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the user ID from the context
	userID, ok := getUserIDFromContext(r)
	if !ok {
		handleError(w, "Unable to retrieve user ID", http.StatusInternalServerError)
		return
	}

	// Extract the content ID from the request URL
	contentID := mux.Vars(r)["id"]

	// Convert the contentID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(contentID)
	if err != nil {
		handleError(w, "Invalid content ID", http.StatusBadRequest)
		return
	}

	collection := database.GetCollection("contents")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find the content by ID and ensure it belongs to the current user
	filter := bson.M{"_id": objectID, "user_id": userID}
	var content models.Content
	err = collection.FindOne(ctx, filter).Decode(&content)
	if err != nil {
		handleError(w, "Content not found or unauthorized", http.StatusNotFound)
		return
	}

	// Delete the content
	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		handleError(w, "Error deleting content", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Content deleted successfully"))
}
