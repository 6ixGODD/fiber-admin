package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocumentationModel struct {
	DocumentID primitive.ObjectID `json:"document_id" bson:"_id"`       // Mongo ObjectId
	Title      string             `json:"title" bson:"title"`           // Title of the document
	Content    string             `json:"content" bson:"content"`       // Content of the document
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"` // Create Time
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"` // Update Time
}
