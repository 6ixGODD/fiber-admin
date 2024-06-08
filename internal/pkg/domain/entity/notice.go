package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NoticeModel struct {
	NoticeID   primitive.ObjectID `json:"notice_id" bson:"_id"`           // Mongo ObjectId
	Title      string             `json:"title" bson:"title"`             // Title
	Content    string             `json:"content" bson:"content"`         // Content in Markdown format
	NoticeType string             `json:"notice_type" bson:"notice_type"` // NoticeType, 'URGENT' | 'NORMAL'
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`   // Created Time in ISO 8601
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`   // Updated Time in ISO 8601
}
