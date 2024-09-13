## Table of Contents

1. [Introduction](#introduction)
2. [Environment Setup](#environment-setup)
3. [Database Connection](#database-connection)
4. [API Endpoints](#api-endpoints)
    - [Public Routes](#public-routes)
    - [Private Routes](#private-routes)
5. [Models](#models)
    - [User Model](#user-model)
    - [Content Model](#content-model)
    - [Stack Model](#stack-model)
6. [Handlers](#handlers)
    - [Content Handlers](#content-handlers)
    - [Stack Handlers](#stack-handlers)
7. [Middleware](#middleware)

## Introduction

This project is a CMS server built with Go, Gorilla Mux, and MongoDB. It provides a set of RESTful API endpoints for managing users, content, and stacks.

## Environment Setup

Ensure you have a `.env` file in the root directory with the following variables:

```
PORT=8080
MONGO_URI=your_mongo_uri
MONGO_DBNAME=your_db_name
```

## Database Connection

The database connection is handled in the `database` package. It connects to MongoDB using the URI provided in the `.env` file.

## API Endpoints

### Public Routes

-   `POST /register` - Register a new user

    -   **Request Body:**
        ```json
        {
            "username": "string",
            "email": "string",
            "password": "string"
        }
        ```
    -   **Response:**
        ```json
        {
            "message": "User registered successfully"
        }
        ```

-   `POST /login` - Login a user

    -   **Request Body:**
        ```json
        {
            "email": "string",
            "password": "string"
        }
        ```
    -   **Response:**
        ```json
        {
            "token": "jwt_token"
        }
        ```
    -   **Cookies:** Not needed

-   `POST /logout` - Logout a user

    -   **Response:**
        ```json
        {
            "message": "User logged out successfully"
        }
        ```
    -   **Cookies:** Not needed

-   `GET /contents` - Get all contents

    -   **Response:**
        ```json
        [
            {
                "id": "string",
                "user_id": "string",
                "name": "string",
                "url": "string",
                "imgUrl": "string",
                "stack": [
                    {
                        "id": "string",
                        "name": "string",
                        "color": "string"
                    }
                ]
            }
        ]
        ```
    -   **Cookies:** Not needed

-   `GET /stacks` - Get all stacks
    -   **Response:**
        ```json
        [
            {
                "id": "string",
                "name": "string",
                "color": "string"
            }
        ]
        ```
    -   **Cookies:** Not needed

### Private Routes

-   `POST /content` - Create new content

    -   **Request Body:**
        ```json
        {
            "name": "string",
            "url": "string",
            "imgUrl": "string",
            "stack": ["string"]
        }
        ```
    -   **Response:**
        ```json
        {
            "message": "Content created successfully"
        }
        ```
    -   **Cookies:** JWT token required in Authorization header

-   `GET /content` - Get content for the authenticated user

    -   **Response:**
        ```json
        [
            {
                "id": "string",
                "user_id": "string",
                "name": "string",
                "url": "string",
                "imgUrl": "string",
                "stack": [
                    {
                        "id": "string",
                        "name": "string",
                        "color": "string"
                    }
                ]
            }
        ]
        ```
    -   **Cookies:** JWT token required in Authorization header

-   `PUT /content/{id}` - Edit content by ID

    -   **Request Body:**
        ```json
        {
            "name": "string",
            "url": "string",
            "imgUrl": "string",
            "stack": ["string"]
        }
        ```
    -   **Response:**
        ```json
        {
            "message": "Content updated successfully"
        }
        ```
    -   **Cookies:** JWT token required in Authorization header

-   `DELETE /content/{id}` - Delete content by ID

    -   **Response:**
        ```json
        {
            "message": "Content deleted successfully"
        }
        ```
    -   **Cookies:** JWT token required in Authorization header

-   `POST /stacks` - Create a new stack

    -   **Request Body:**
        ```json
        {
            "name": "string",
            "color": "string"
        }
        ```
    -   **Response:**
        ```json
        {
            "message": "Stack created successfully"
        }
        ```
    -   **Cookies:** JWT token required in Authorization header

-   `PUT /stacks/{id}` - Edit stack by ID

    -   **Request Body:**
        ```json
        {
            "name": "string",
            "color": "string"
        }
        ```
    -   **Response:**
        ```json
        {
            "message": "Stack updated successfully"
        }
        ```
    -   **Cookies:** JWT token required in Authorization header

-   `DELETE /stacks/{id}` - Delete stack by ID
    -   **Response:**
        ```json
        {
            "message": "Stack deleted successfully"
        }
        ```
    -   **Cookies:** JWT token required in Authorization header

## Models

### User Model

```go
type User struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Username string             `bson:"username" json:"username"`
    Email    string             `bson:"email" json:"email"`
    Password string             `bson:"password" json:"-"`
}
```

### Content Model

```go
type Content struct {
    ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    UserID string             `json:"user_id" bson:"user_id"`
    Name   string             `json:"name" bson:"name"`
    Url    string             `json:"url" bson:"url"`
    ImgUrl string             `json:"imgUrl" bson:"imgUrl"`
    Stack  []Stack            `json:"stack" bson:"stack"`
}
```

### Stack Model

```go
type Stack struct {
    ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name  string             `bson:"name" json:"name"`
    Color string             `bson:"color" json:"color"`
}
```

## Handlers

### Content Handlers

-   `CreateContentHandler` - Handles the creation of new content
-   `GetContentsHandler` - Retrieves all content from the database
-   `GetContentHandler` - Retrieves content for a specific user
-   `EditContentHandler` - Updates an existing content by ID
-   `DeleteContentHandler` - Deletes a content by ID

### Stack Handlers

-   `CreateStackHandler` - Handles the creation of a new stack
-   `GetStacksHandler` - Retrieves all stacks from the database
-   `EditStackHandler` - Updates an existing stack by ID
-   `DeleteStackHandler` - Deletes a stack by ID

## Middleware

The middleware package includes authentication middleware to protect private routes.

## Authentication Middleware

The authentication middleware ensures that only authenticated users can access private routes. It checks for a valid JWT token in the request headers and verifies it.

```go
package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key")

// AuthMiddleware is a middleware function for authenticating requests
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header missing", http.StatusUnauthorized)
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims := &Claims{}

        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), "userID", claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Claims defines the structure of the JWT claims
type Claims struct {
    UserID string `json:"user_id"`
    jwt.StandardClaims
}
```

## Running the Server

To run the server, use the following command:

```sh
go run main.go
```

Ensure that your MongoDB instance is running and the `.env` file is properly configured.

## Conclusion

This CMS server provides a robust set of API endpoints for managing users, content, and stacks. It uses MongoDB for data storage and Gorilla Mux for routing. The authentication middleware ensures that only authenticated users can access private routes.

Feel free to contribute to this project by submitting issues or pull requests on GitHub.
