package entity

import (
	"time"
)

type CacheList struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}

type UserCacheList struct {
	Total int64       `json:"total"`
	List  []UserCache `json:"list"`
}

type UserCache struct {
	UserID       string    `json:"user_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Role         string    `json:"role"`
	Organization string    `json:"organization"`
	LastLogin    time.Time `json:"last_login"`
	Deleted      bool      `json:"deleted"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

type NoticeCacheList struct {
	Total int64         `json:"total"`
	List  []NoticeCache `json:"list"`
}

type NoticeCache struct {
	NoticeID   string    `json:"notice_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	NoticeType string    `json:"notice_type"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DocumentationCacheList struct {
	Total int64                `json:"total"`
	List  []DocumentationCache `json:"list"`
}

type DocumentationCache struct {
	DocumentID string    `json:"document_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
type LoginLogCache struct {
	UserIDHex string    `json:"user_id_hex"` // User ID in Hex
	IPAddress string    `json:"ip_address"`  // IP Address
	UserAgent string    `json:"user_agent"`  // User Agent
	CreatedAt time.Time `json:"created_at"`  // Created Time in ISO 8601
}

type OperationLogCache struct {
	UserIDHex   string    `json:"user_id_hex"`   // User ID in Hex
	IPAddress   string    `json:"ip_address"`    // IP Address
	UserAgent   string    `json:"user_agent"`    // User Agent
	Operation   string    `json:"operation"`     // Operation, 'CREATE' | 'UPDATE' | 'DELETE'
	EntityIDHex string    `json:"entity_id_hex"` // Entity ID in Hex
	EntityType  string    `json:"entity_type"`   // Entity noticeType, 'USER' | 'DOCUMENTATION' | 'NOTICE'
	Description string    `json:"description"`   // Description of Operation
	Status      string    `json:"status"`        // Status, 'SUCCESS' | 'FAILURE'
	CreatedAt   time.Time `json:"created_at"`    // Created Time in ISO 8601
}
