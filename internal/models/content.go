package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Stack struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name  string             `bson:"name" json:"name"`
	Color string             `bson:"color" json:"color"`
}

type Content struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID      string             `json:"user_id" bson:"user_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Url         string             `json:"url" bson:"url"`
	ImgUrl      string             `json:"imgUrl" bson:"imgUrl"`
	Stack       []Stack            `json:"stack" bson:"stack"`
}
