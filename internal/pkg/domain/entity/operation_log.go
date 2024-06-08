package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OperationLogModel struct {
	OperationLogID primitive.ObjectID `json:"operation_log_id" bson:"_id"`    // Mongo ObjectId
	UserID         primitive.ObjectID `json:"user_id" bson:"user_id"`         // User ID
	Username       string             `json:"username" bson:"username"`       // Username (for space-time trade-off)
	Email          string             `json:"email" bson:"email"`             // Email (for space-time trade-off)
	IPAddress      string             `json:"ip_address" bson:"ip_address"`   // IP Address
	UserAgent      string             `json:"user_agent" bson:"user_agent"`   // User Agent
	Operation      string             `json:"operation" bson:"operation"`     // Operation, 'CREATE' | 'UPDATE' | 'DELETE'
	EntityID       primitive.ObjectID `json:"entity_id" bson:"entity_id"`     // Entity ID
	EntityType     string             `json:"entity_type" bson:"entity_type"` // Entity noticeType, 'USER' | 'DOCUMENTATION' | 'NOTICE'
	Description    string             `json:"description" bson:"description"` // Description of Operation
	Status         string             `json:"status" bson:"status"`           // Status, 'SUCCESS' | 'FAILURE'
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`   // Created Time in ISO 8601
}
