package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserModel struct {
	UserID       primitive.ObjectID `json:"user_id" bson:"_id"`               // Mongo ObjectId
	Username     string             `json:"username" bson:"username"`         // Username
	Email        string             `json:"email" bson:"email"`               // Email
	Password     string             `json:"password" bson:"password"`         // Password crypt
	Role         string             `json:"role" bson:"role"`                 // Role, 'USER' | 'ADMIN'
	Organization string             `json:"organization" bson:"organization"` // Organization
	LastLogin    time.Time          `json:"last_login" bson:"last_login"`     // Last Login Time in ISO 8601
	Deleted      bool               `json:"deleted" bson:"deleted"`           // Deleted Flag
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`     // Created Time in ISO 8601
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`     // Updated Time in ISO 8601
	DeletedAt    time.Time          `json:"deleted_at" bson:"deleted_at"`     // Deleted Time in ISO 8601
}
