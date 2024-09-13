package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"cms-server/internal/database"
	"cms-server/internal/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateStackHandler handles the creation of a new stack
func CreateStackHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Create a new instance of Stack
	var stack models.Stack

	// Decode the request body into the stack struct
	err := json.NewDecoder(r.Body).Decode(&stack)
	if err != nil || stack.Name == "" || stack.Color == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the stacks collection
	collection := database.GetCollection("stacks")

	// Set a timeout context for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if the stack name already exists
	var existingStack models.Stack
	err = collection.FindOne(ctx, bson.M{"name": stack.Name}).Decode(&existingStack)
	if err == nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"message": "Stack with the same name already exists"})
		return
	}

	// Assign a new ObjectID to the stack
	stack.ID = primitive.NewObjectID()

	// Insert the stack into the database
	_, err = collection.InsertOne(ctx, stack)
	if err != nil {
		http.Error(w, "Error inserting stack", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(stack)
}

// GetStacksHandler retrieves all stacks from the database
func GetStacksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the stacks collection
	collection := database.GetCollection("stacks")

	// Set up a context with a timeout for querying MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a slice to hold the stack documents
	var stacks []models.Stack

	// Retrieve all documents
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Error fetching stacks", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	// Decode documents into the stacks slice
	for cursor.Next(ctx) {
		var stack models.Stack
		if err := cursor.Decode(&stack); err != nil {
			http.Error(w, "Error decoding stack", http.StatusInternalServerError)
			return
		}
		stacks = append(stacks, stack)
	}

	// Check for cursor iteration errors
	if err := cursor.Err(); err != nil {
		http.Error(w, "Error fetching stacks", http.StatusInternalServerError)
		return
	}

	// Return the stacks as JSON
	json.NewEncoder(w).Encode(stacks)
}

// EditStackHandler updates an existing stack by ID
func EditStackHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract stack ID from the request URL
	params := mux.Vars(r)
	stackID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid stack ID", http.StatusBadRequest)
		return
	}

	// Decode the request body into a stack struct
	var updatedStack models.Stack
	err = json.NewDecoder(r.Body).Decode(&updatedStack)
	if err != nil || updatedStack.Name == "" || updatedStack.Color == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set the update fields
	update := bson.M{
		"$set": bson.M{
			"name":  updatedStack.Name,
			"color": updatedStack.Color,
		},
	}

	// Get the stacks collection
	collection := database.GetCollection("stacks")

	// Set a timeout context for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the stack in the database
	_, err = collection.UpdateOne(ctx, bson.M{"_id": stackID}, update)
	if err != nil {
		http.Error(w, "Error updating stack", http.StatusInternalServerError)
		return
	}

	updatedStack.ID = stackID

	// Return the updated stack as a response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedStack)
}

// DeleteStackHandler deletes a stack by ID
func DeleteStackHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract stack ID from the request URL
	params := mux.Vars(r)
	stackID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid stack ID", http.StatusBadRequest)
		return
	}

	// Get the stacks collection
	collection := database.GetCollection("stacks")

	// Set a timeout context for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete the stack from the database
	_, err = collection.DeleteOne(ctx, bson.M{"_id": stackID})
	if err != nil {
		http.Error(w, "Error deleting stack", http.StatusInternalServerError)
		return
	}

	// Return success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Stack deleted successfully"})
}
