package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginLogModel struct {
	LoginLogID primitive.ObjectID `json:"login_log_id" bson:"_id"`      // Mongo ObjectId
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`       // User ID
	Username   string             `json:"username" bson:"username"`     // Username (for space-time trade-off)
	Email      string             `json:"email" bson:"email"`           // Email (for space-time trade-off)
	IPAddress  string             `json:"ip_address" bson:"ip_address"` // IP Address
	UserAgent  string             `json:"user_agent" bson:"user_agent"` // User Agent
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"` // Created Time in ISO 8601
}
